"""Unit tests for defensive AI Annotation Guard against zero-shot hallucinations."""
import pytest
from curation.ai_guard import AIAnnotationGuard


@pytest.mark.unit
def test_ai_guard_filters_hallucinated_classes():
    # Classes autorizadas: 0 (Person), 1 (Car)
    guard = AIAnnotationGuard(allowed_classes={0, 1})

    raw_predictions = [
        {"class_id": 0, "x_center": 0.5, "y_center": 0.5, "w": 0.2, "h": 0.2, "confidence": 0.85},
        # Classe alucinada 99
        {"class_id": 99, "x_center": 0.3, "y_center": 0.3, "w": 0.1, "h": 0.1, "confidence": 0.92},
    ]

    sanitized = guard.sanitize_predictions(raw_predictions)
    assert len(sanitized) == 1
    assert sanitized[0].class_id == 0


@pytest.mark.unit
def test_ai_guard_filters_low_confidence_and_degenerate_area():
    guard = AIAnnotationGuard(allowed_classes={0, 1}, min_confidence=0.5, min_area=0.001, max_area=0.95)

    raw_predictions = [
        # Baixa confiança (0.3 < 0.5)
        {"class_id": 0, "x_center": 0.5, "y_center": 0.5, "w": 0.2, "h": 0.2, "confidence": 0.3},
        # Área gigantesca que toma 100% da imagem (fundo erroneamente detectado)
        {"class_id": 1, "x_center": 0.5, "y_center": 0.5, "w": 0.99, "h": 0.99, "confidence": 0.9},
        # Válido
        {"class_id": 1, "x_center": 0.5, "y_center": 0.5, "w": 0.2, "h": 0.3, "confidence": 0.88},
    ]

    sanitized = guard.sanitize_predictions(raw_predictions)
    assert len(sanitized) == 1
    assert sanitized[0].class_id == 1
    assert pytest.approx(sanitized[0].area, 1e-4) == 0.06


@pytest.mark.unit
def test_ai_guard_deduplicates_near_identical_boxes():
    guard = AIAnnotationGuard(allowed_classes={0}, max_iou_duplicate=0.85)

    raw_predictions = [
        {"class_id": 0, "x_center": 0.5, "y_center": 0.5, "w": 0.2, "h": 0.2, "confidence": 0.90},
        # Quase a mesma caixa com variação minúscula
        {"class_id": 0, "x_center": 0.501, "y_center": 0.501, "w": 0.2, "h": 0.2, "confidence": 0.87},
    ]

    sanitized = guard.sanitize_predictions(raw_predictions)
    assert len(sanitized) == 1
