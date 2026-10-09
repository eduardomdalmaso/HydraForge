"""
HydraForge Person-Gated Face Tracker.
Associates face detections with persistent person tracks.
"""
from typing import List, Dict, Tuple, Optional
import numpy as np


def compute_iou(boxA: Tuple[int, int, int, int], boxB: Tuple[int, int, int, int]) -> float:
    """Computes Intersection over Union (IoU) between two bounding boxes (xyxy)."""
    xA = max(boxA[0], boxB[0])
    yA = max(boxA[1], boxB[1])
    xB = min(boxA[2], boxB[2])
    yB = min(boxA[3], boxB[3])

    inter = max(0, xB - xA) * max(0, yB - yA)
    if inter == 0:
        return 0.0

    areaA = (boxA[2] - boxA[0]) * (boxA[3] - boxA[1])
    areaB = (boxB[2] - boxB[0]) * (boxB[3] - boxB[1])
    return float(inter / (areaA + areaB - inter))


class TrackedPerson:
    def __init__(self, track_id: int, bbox_xyxy: Tuple[int, int, int, int]):
        self.track_id = track_id
        self.bbox_xyxy = bbox_xyxy
        self.face_hits = 0
        self.missed_frames = 0
        self.last_face_bbox: Optional[Tuple[int, int, int, int]] = None
        self.subject_id: Optional[str] = None
        self.subject_name: Optional[str] = None
        self.is_recognized = False

    def update_box(self, new_bbox: Tuple[int, int, int, int]):
        self.bbox_xyxy = new_bbox
        self.missed_frames = 0

    def contains_face(self, face_bbox: Tuple[int, int, int, int]) -> bool:
        """
        Validates if face resides within the upper anatomical region (top 45%) of person box.
        """
        fx1, fy1, fx2, fy2 = face_bbox
        px1, py1, px2, py2 = self.bbox_xyxy
        face_center_x = (fx1 + fx2) / 2.0
        face_center_y = (fy1 + fy2) / 2.0

        p_height = py2 - py1
        p_head_bottom = py1 + (p_height * 0.45)

        # Relaxed horizontal boundary (+15% margin)
        p_width = px2 - px1
        margin_x = p_width * 0.15

        x_inside = (px1 - margin_x) <= face_center_x <= (px2 + margin_x)
        y_inside = (py1 - margin_x) <= face_center_y <= p_head_bottom
        return x_inside and y_inside


class PersonGatedTracker:
    """
    Orchestrates person tracks and assigns face detections to corresponding individuals.
    """

    def __init__(self, max_missed: int = 30, iou_match_threshold: float = 0.35):
        self.tracks: Dict[int, TrackedPerson] = {}
        self.next_track_id = 1
        self.max_missed = max_missed
        self.iou_match_threshold = iou_match_threshold

    def update_persons(self, person_boxes: List[Tuple[int, int, int, int]]) -> List[TrackedPerson]:
        """Updates person trajectories matching current detections."""
        unmatched_dets = list(person_boxes)
        matched_tracks = []

        for track in list(self.tracks.values()):
            best_iou = 0.0
            best_match_idx = -1
            for idx, det in enumerate(unmatched_dets):
                iou = compute_iou(track.bbox_xyxy, det)
                if iou > best_iou:
                    best_iou = iou
                    best_match_idx = idx

            if best_iou >= self.iou_match_threshold and best_match_idx >= 0:
                track.update_box(unmatched_dets.pop(best_match_idx))
                matched_tracks.append(track)
            else:
                track.missed_frames += 1
                if track.missed_frames > self.max_missed:
                    del self.tracks[track.track_id]

        # Register new person tracks
        for new_det in unmatched_dets:
            new_t = TrackedPerson(self.next_track_id, new_det)
            self.tracks[self.next_track_id] = new_t
            self.next_track_id += 1
            matched_tracks.append(new_t)

        return matched_tracks

    def match_face_to_person(
        self, face_bbox: Tuple[int, int, int, int]
    ) -> Optional[TrackedPerson]:
        """Finds the best matching person track for a detected face."""
        candidates = []
        for track in self.tracks.values():
            if track.contains_face(face_bbox):
                candidates.append(track)

        if not candidates:
            return None
        # If multiple overlap, pick the one where face is closest to head center
        fx_c = (face_bbox[0] + face_bbox[2]) / 2.0
        fy_c = (face_bbox[1] + face_bbox[3]) / 2.0
        return min(
            candidates,
            key=lambda t: abs(fx_c - (t.bbox_xyxy[0] + t.bbox_xyxy[2]) / 2.0)
        )
