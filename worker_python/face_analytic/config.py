"""
HydraForge Face Analytics Configuration Module.
Defines parameters, constraints, and validation rules.
"""
from dataclasses import dataclass, field
from typing import Dict, Any


@dataclass
class FaceAnalyticConfig:
    min_height: int = 40
    max_height: int = 0  # 0 indicates unbounded / infinity
    min_faces: int = 25  # Minimum hits associated with a person track (Vezha standard)
    min_confidence: float = 0.60
    blur_threshold: float = 60.0
    adaptive_blur: bool = True
    unknown_score: float = 0.60
    recognition_threshold: float = 0.75
    best_frame_buffer_size: int = 15
    camera_id: str = "cam_01"
    iou_person_face_threshold: float = 0.30
    save_best_frame_disk: bool = True
    gallery_path: str = ""
    detector_sensitivity: str = "normal"  # "baixo", "normal", "alto"
    detect_gender_age: bool = False

    def validate(self) -> None:
        """Validates configuration boundaries and logical consistency."""
        if self.min_height < 10:
            raise ValueError(f"min_height must be >= 10, got {self.min_height}")
        if self.max_height < 0:
            raise ValueError(f"max_height must be >= 0 (0 means unbounded), got {self.max_height}")
        if self.max_height > 0 and self.max_height < self.min_height:
            raise ValueError(
                f"max_height ({self.max_height}) cannot be less than min_height ({self.min_height})"
            )
        if self.min_faces < 1:
            raise ValueError(f"min_faces must be >= 1, got {self.min_faces}")
        if not (0.1 <= self.min_confidence <= 1.0):
            raise ValueError(f"min_confidence must be in [0.1, 1.0], got {self.min_confidence}")
        if self.blur_threshold <= 0:
            raise ValueError(f"blur_threshold must be positive, got {self.blur_threshold}")
        if not (0.1 <= self.unknown_score < self.recognition_threshold <= 1.0):
            raise ValueError(
                f"Requires 0.1 <= unknown_score ({self.unknown_score}) < "
                f"recognition_threshold ({self.recognition_threshold}) <= 1.0"
            )

    def is_height_valid(self, height: int) -> bool:
        """Evaluates whether detected face height satisfies constraints."""
        if height < self.min_height:
            return False
        if self.max_height > 0 and height > self.max_height:
            return False
        return True

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "FaceAnalyticConfig":
        valid_keys = {f.name for f in cls.__dataclass_fields__.values()}
        filtered = {k: v for k, v in data.items() if k in valid_keys}
        cfg = cls(**filtered)
        cfg.validate()
        return cfg
