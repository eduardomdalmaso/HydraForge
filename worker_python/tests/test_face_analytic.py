"""
Comprehensive Unit Tests for HydraForge Face Analytics Pipeline.
Tests height filters (min/max=0), person gating, best-frame blur, and cross-camera ReID.
"""
import pytest
import numpy as np

from face_analytic.config import FaceAnalyticConfig
from face_analytic.quality import (
    compute_laplacian_variance,
    compute_aspect_ratio_score,
    compute_composite_quality
)
from face_analytic.best_frame import BestFrameBuffer
from face_analytic.tracker import PersonGatedTracker, compute_iou
from face_analytic.matcher import FaceEmbeddingMatcher, cosine_similarity
from face_analytic.pipeline import FaceAnalyticPipeline


def test_config_height_rules():
    # max_height = 0 means unbounded / infinity
    cfg = FaceAnalyticConfig(min_height=40, max_height=0, min_faces=3)
    cfg.validate()
    assert cfg.is_height_valid(39) is False
    assert cfg.is_height_valid(40) is True
    assert cfg.is_height_valid(1000) is True

    # If max_height > 0, enforces range
    cfg2 = FaceAnalyticConfig(min_height=40, max_height=120)
    cfg2.validate()
    assert cfg2.is_height_valid(35) is False
    assert cfg2.is_height_valid(80) is True
    assert cfg2.is_height_valid(150) is False

    # Invalid configurations raise ValueError
    with pytest.raises(ValueError):
        FaceAnalyticConfig(min_height=5).validate()
    with pytest.raises(ValueError):
        FaceAnalyticConfig(min_height=100, max_height=50).validate()
    with pytest.raises(ValueError):
        FaceAnalyticConfig(min_faces=0).validate()


def test_quality_and_blur_scoring():
    # Synthetic flat image has zero variance (blurred)
    flat_crop = np.full((60, 50, 3), 128, dtype=np.uint8)
    var_flat = compute_laplacian_variance(flat_crop)
    assert var_flat == 0.0

    # Textured image has high edge variance (sharp)
    sharp_crop = np.random.randint(0, 255, (60, 50, 3), dtype=np.uint8)
    var_sharp = compute_laplacian_variance(sharp_crop)
    assert var_sharp > 50.0

    score_flat, _, passes_flat = compute_composite_quality(flat_crop, 0.90, blur_threshold=60.0)
    score_sharp, _, passes_sharp = compute_composite_quality(sharp_crop, 0.90, blur_threshold=60.0)

    assert passes_flat is False
    assert passes_sharp is True
    assert score_sharp > score_flat


def test_best_frame_buffer_and_adaptive_blur():
    buf = BestFrameBuffer(track_id=1, max_size=5, blur_threshold=50.0)
    dummy_frame = np.zeros((100, 100, 3), dtype=np.uint8)
    dummy_crop = np.zeros((40, 40, 3), dtype=np.uint8)

    # Add 3 moving/blurry frames (blur_score=25, below threshold 50)
    for i in range(3):
        buf.add_candidate(
            frame=dummy_frame,
            crop=dummy_crop,
            bbox_xyxy=(10, 10, 50, 50),
            confidence=0.85 + (i * 0.02),
            blur_score=25.0 + (i * 5.0),
            composite_score=0.40 + (i * 0.10),
            passes_blur=False
        )

    # Without adaptive blur, no candidate returned because all were below threshold
    assert buf.get_best_candidate(adaptive_blur=False) is None

    # With adaptive blur enabled, picks the relative peak candidate
    best_cand = buf.get_best_candidate(adaptive_blur=True)
    assert best_cand is not None
    assert best_cand.blur_score == 35.0  # highest in the blurry sequence

    # Now a sharp frame arrives!
    buf.add_candidate(
        frame=dummy_frame,
        crop=dummy_crop,
        bbox_xyxy=(10, 10, 50, 50),
        confidence=0.95,
        blur_score=85.0,
        composite_score=0.92,
        passes_blur=True
    )
    best_sharp = buf.get_best_candidate(adaptive_blur=True)
    assert best_sharp.blur_score == 85.0


def test_person_gated_tracking():
    tracker = PersonGatedTracker()
    # Frame 1: Person at (100, 100, 200, 400)
    persons = tracker.update_persons([(100, 100, 200, 400)])
    assert len(persons) == 1
    track = persons[0]

    # Face inside head region (120, 110, 170, 160) -> matches
    matched = tracker.match_face_to_person((120, 110, 170, 160))
    assert matched is not None
    assert matched.track_id == track.track_id

    # Face at the feet (120, 350, 170, 390) -> rejected
    assert tracker.match_face_to_person((120, 350, 170, 390)) is None


def test_cross_camera_reid_and_known_unknown():
    matcher = FaceEmbeddingMatcher(unknown_score=0.60, recognition_threshold=0.75)

    # Register known employee
    emp_emb = np.zeros(512, dtype=np.float32)
    emp_emb[0] = 1.0  # Unit vector
    matcher.register_known_subject("EMP_001", "Alice Silva", emp_emb)

    # Query with Alice's embedding on Camera 1 -> Known
    alice_query = emp_emb + np.random.randn(512).astype(np.float32) * 0.01
    s_id, name, is_known, conf = matcher.identify_or_register(alice_query, camera_id="cam_01")
    assert s_id == "EMP_001"
    assert name == "Alice Silva"
    assert is_known is True
    assert conf >= 0.75

    # Query with unknown visitor on Camera 1 -> Unknown created
    visitor_emb = np.zeros(512, dtype=np.float32)
    visitor_emb[10] = 1.0
    u_id, u_name, u_known, _ = matcher.identify_or_register(visitor_emb, camera_id="cam_01")
    assert u_known is False
    assert "UNKNOWN" in u_id

    # Same visitor arrives at Camera 2! -> Re-identified across cameras
    visitor_cam2 = visitor_emb + np.random.randn(512).astype(np.float32) * 0.01
    reid_id, _, _, reid_conf = matcher.identify_or_register(visitor_cam2, camera_id="cam_02")
    assert reid_id == u_id  # Re-identified as same visitor!
    assert reid_conf >= 0.75
    assert matcher.active_unknowns[u_id].camera_trail == ["cam_01", "cam_02"]


def test_end_to_end_pipeline():
    cfg = FaceAnalyticConfig(min_height=30, max_height=0, min_faces=2, camera_id="gate_01")
    pipeline = FaceAnalyticPipeline(cfg)
    frame = np.random.randint(0, 255, (480, 640, 3), dtype=np.uint8)

    # Person and Face boxes
    person_box = (100, 50, 250, 450)
    face_box = (130, 60, 190, 130)  # height = 70px

    # Frame 1: 1st hit -> no event yet (min_faces = 2)
    ev1 = pipeline.process_frame(frame, [person_box], [face_box])
    assert len(ev1) == 0

    # Frame 2: 2nd hit -> event emitted!
    ev2 = pipeline.process_frame(frame, [person_box], [face_box])
    assert len(ev2) == 1
    event = ev2[0]
    assert event["event_type"] == "AI.DETECTION.FACE_RECOGNITION"
    assert event["camera_id"] == "gate_01"
    assert event["face_hits"] == 2
    assert "best_frame_quality" in event
