"""
HydraForge Best Frame Selector Module.
Maintains best-frame buffers per track with adaptive blur handling.
"""
import time
from dataclasses import dataclass
from typing import List, Optional, Tuple
import numpy as np


@dataclass
class CandidateFace:
    frame: np.ndarray
    crop: np.ndarray
    bbox_xyxy: Tuple[int, int, int, int]
    confidence: float
    blur_score: float
    composite_score: float
    passes_blur: bool
    timestamp: float
    track_id: int


class BestFrameBuffer:
    """
    Manages face candidates for a single track across video frames.
    Selects the optimal frame using multi-criteria quality scoring while
    employing adaptive fallback to avoid losing detections during motion.
    """

    def __init__(self, track_id: int, max_size: int = 15, blur_threshold: float = 60.0):
        self.track_id = track_id
        self.max_size = max_size
        self.blur_threshold = blur_threshold
        self.candidates: List[CandidateFace] = []
        self.total_seen = 0
        self.sharp_count = 0

    def add_candidate(
        self,
        frame: np.ndarray,
        crop: np.ndarray,
        bbox_xyxy: Tuple[int, int, int, int],
        confidence: float,
        blur_score: float,
        composite_score: float,
        passes_blur: bool
    ) -> None:
        self.total_seen += 1
        if passes_blur:
            self.sharp_count += 1

        cand = CandidateFace(
            frame=frame.copy() if frame is not None else None,
            crop=crop.copy() if crop is not None else None,
            bbox_xyxy=bbox_xyxy,
            confidence=confidence,
            blur_score=blur_score,
            composite_score=composite_score,
            passes_blur=passes_blur,
            timestamp=time.time(),
            track_id=self.track_id
        )

        self.candidates.append(cand)
        # Keep buffer sorted by composite score descending
        self.candidates.sort(key=lambda c: c.composite_score, reverse=True)
        if len(self.candidates) > self.max_size:
            self.candidates = self.candidates[:self.max_size]

    def get_best_candidate(self, adaptive_blur: bool = True) -> Optional[CandidateFace]:
        """
        Retrieves the optimal face candidate for recognition.
        Prioritizes sharp frames; if none exist and adaptive_blur is True,
        selects the highest relative quality frame seen during the track.
        """
        if not self.candidates:
            return None

        sharp_candidates = [c for c in self.candidates if c.passes_blur]
        if sharp_candidates:
            # Pick highest composite score among sharp frames
            return sharp_candidates[0]

        if not adaptive_blur:
            return None

        # Adaptive fallback: If person was in motion throughout the track,
        # choose the relative peak frame (highest blur score and confidence)
        # with at least 30% of normal blur threshold to prevent pure noise.
        relaxed_min = self.blur_threshold * 0.30
        acceptable = [c for c in self.candidates if c.blur_score >= relaxed_min]
        if acceptable:
            return max(acceptable, key=lambda c: c.composite_score)

        return self.candidates[0]

    def has_sufficient_quality(self, min_faces: int, adaptive_blur: bool = True) -> bool:
        """Determines if the track has collected enough qualifying face observations."""
        if self.total_seen < min_faces:
            return False
        best = self.get_best_candidate(adaptive_blur=adaptive_blur)
        return best is not None
