"""Defensive AI Annotation Guard against zero-shot model hallucinations (SAM 2 / YOLO-World)."""
from __future__ import annotations
from typing import Any, Dict, List, Optional, Set
from .domain import BoundingBox


class AIAnnotationGuard:
    """Sanitizes and enforces strict deterministic invariants on AI-generated predictions."""

    def __init__(
        self,
        allowed_classes: Set[int] | List[int],
        min_area: float = 0.0001,
        max_area: float = 0.98,
        min_confidence: float = 0.25,
        max_iou_duplicate: float = 0.90,
    ):
        self.allowed_classes = set(allowed_classes)
        self.min_area = min_area
        self.max_area = max_area
        self.min_confidence = min_confidence
        self.max_iou_duplicate = max_iou_duplicate

    def sanitize_predictions(self, raw_predictions: List[Dict[str, Any]]) -> List[BoundingBox]:
        """Filters and validates raw AI predictions, stripping out hallucinations and corruptions."""
        valid_boxes: List[BoundingBox] = []

        for item in raw_predictions:
            try:
                cls_id = int(item["class_id"])
                x_center = float(item["x_center"])
                y_center = float(item["y_center"])
                w = float(item["w"])
                h = float(item["h"])
                conf = float(item.get("confidence", 1.0))
            except (KeyError, ValueError, TypeError):
                # Reject malformed items
                continue

            # Salvaguarda 1: Escopo de classes permitidas
            if cls_id not in self.allowed_classes:
                continue

            # Salvaguarda 2: Limiar de confiança
            if conf < self.min_confidence:
                continue

            box = BoundingBox(
                class_id=cls_id,
                x_center=x_center,
                y_center=y_center,
                w=w,
                h=h,
                confidence=conf,
            )

            # Salvaguarda 3: Validação matemática YOLO [0.0, 1.0]
            if not box.is_valid_yolo():
                continue

            # Salvaguarda 4: Área mínima e máxima (rejeita ruído e background que toma tela toda)
            if box.area < self.min_area or box.area > self.max_area:
                continue

            # Salvaguarda 5: Rejeição de caixas duplicadas com IoU quase idêntico
            is_duplicate = False
            for existing in valid_boxes:
                if existing.class_id == box.class_id and existing.iou(box) > self.max_iou_duplicate:
                    is_duplicate = True
                    break

            if not is_duplicate:
                valid_boxes.append(box)

        return valid_boxes
