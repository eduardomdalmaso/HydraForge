"""Pure functional core for dataset curation, bounding boxes and labels."""
from __future__ import annotations
import math
from dataclasses import dataclass
from typing import List, Optional, Protocol


@dataclass(frozen=True, slots=True)
class BoundingBox:
    """Normalized YOLO bounding box [class_id, x_center, y_center, w, h].
    
    Optimized based on High Performance Python Chapter 11: `slots=True` eliminates the
    `__dict__` overhead, reducing memory usage by ~70% per instance.
    """
    class_id: int
    x_center: float
    y_center: float
    w: float
    h: float
    confidence: float = 1.0

    def is_valid_yolo(self) -> bool:
        """Validates that all coordinates and dimensions are strictly within [0.0, 1.0] and non-degenerate."""
        values = [self.x_center, self.y_center, self.w, self.h]
        for v in values:
            if math.isnan(v) or math.isinf(v) or v < 0.0 or v > 1.0:
                return False
        # Width and height must be strictly positive to avoid degenerate zero-area boxes
        if self.w <= 0.0 or self.h <= 0.0:
            return False
        # Box bounds must stay inside image boundary
        x_min = self.x_center - (self.w / 2.0)
        x_max = self.x_center + (self.w / 2.0)
        y_min = self.y_center - (self.h / 2.0)
        y_max = self.y_center + (self.h / 2.0)
        
        # Tolerance for float rounding (e.g. 1.0000001)
        eps = 1e-6
        if x_min < -eps or y_min < -eps or x_max > (1.0 + eps) or y_max > (1.0 + eps):
            return False

        return self.class_id >= 0

    @property
    def area(self) -> float:
        """Returns the normalized area (0.0 to 1.0]."""
        return self.w * self.h

    def iou(self, other: BoundingBox) -> float:
        """Computes Intersection over Union (IoU) with another bounding box."""
        x1_min = self.x_center - (self.w / 2.0)
        x1_max = self.x_center + (self.w / 2.0)
        y1_min = self.y_center - (self.h / 2.0)
        y1_max = self.y_center + (self.h / 2.0)

        x2_min = other.x_center - (other.w / 2.0)
        x2_max = other.x_center + (other.w / 2.0)
        y2_min = other.y_center - (other.h / 2.0)
        y2_max = other.y_center + (other.h / 2.0)

        inter_xmin = max(x1_min, x2_min)
        inter_ymin = max(y1_min, y2_min)
        inter_xmax = min(x1_max, x2_max)
        inter_ymax = min(y1_max, y2_max)

        inter_w = max(0.0, inter_xmax - inter_xmin)
        inter_h = max(0.0, inter_ymax - inter_ymin)
        intersection = inter_w * inter_h

        union = self.area + other.area - intersection
        if union <= 0.0:
            return 0.0
        return intersection / union


class JobRepository(Protocol):
    """Port for curation job persistence."""
    def save(self, job_id: str, bboxes: List[BoundingBox]) -> None: ...
    def get_bboxes(self, job_id: str) -> List[BoundingBox]: ...
    def delete(self, job_id: str) -> bool: ...
