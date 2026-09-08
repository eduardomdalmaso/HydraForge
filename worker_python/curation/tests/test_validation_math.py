"""Property-based and triangulation tests for bounding box math and invariants."""
import pytest
from hypothesis import given, strategies as st
from curation.domain import BoundingBox


# ---------------------------------------------------------------------------
# 1. Triangulação: Casos de borda e limites discretos
# ---------------------------------------------------------------------------

@pytest.mark.unit
def test_triangulation_boundary_cases():
    """Triangulação com valores discretos nos limites 0.0 e 1.0."""
    # Caixa central válida
    assert BoundingBox(0, 0.5, 0.5, 0.5, 0.5).is_valid_yolo() is True

    # Coordenadas negativas
    assert BoundingBox(0, -0.01, 0.5, 0.2, 0.2).is_valid_yolo() is False
    assert BoundingBox(0, 0.5, -0.01, 0.2, 0.2).is_valid_yolo() is False

    # Coordenadas acima de 1.0
    assert BoundingBox(0, 1.01, 0.5, 0.2, 0.2).is_valid_yolo() is False
    assert BoundingBox(0, 0.5, 1.01, 0.2, 0.2).is_valid_yolo() is False

    # Dimensões degeneradas (zero ou negativas)
    assert BoundingBox(0, 0.5, 0.5, 0.0, 0.2).is_valid_yolo() is False
    assert BoundingBox(0, 0.5, 0.5, 0.2, -0.1).is_valid_yolo() is False

    # Caixa cujo centro está em 0.9 e largura 0.4 (estende além de 1.0)
    assert BoundingBox(0, 0.9, 0.5, 0.4, 0.2).is_valid_yolo() is False


# ---------------------------------------------------------------------------
# 2. Property-Based Testing com Hypothesis
# ---------------------------------------------------------------------------

@pytest.mark.property
@given(
    class_id=st.integers(min_value=0, max_value=80),
    w=st.floats(min_value=0.01, max_value=0.90),
    h=st.floats(min_value=0.01, max_value=0.90),
)
def test_property_valid_contained_boxes(class_id: int, w: float, h: float):
    """Gera caixas geometricamente seguras e garante que são sempre válidas."""
    # Escolhe centro garantido para não vazar a borda
    x_center = w / 2.0 + (1.0 - w) * 0.5
    y_center = h / 2.0 + (1.0 - h) * 0.5
    
    bbox = BoundingBox(class_id=class_id, x_center=x_center, y_center=y_center, w=w, h=h)
    assert bbox.is_valid_yolo() is True
    assert 0.0 < bbox.area <= 1.0


@pytest.mark.property
@given(
    invalid_coord=st.one_of(
        st.floats(max_value=-0.001),
        st.floats(min_value=1.001),
        st.sampled_from([float("nan"), float("inf"), float("-inf")]),
    )
)
def test_property_rejects_corrupted_coordinates(invalid_coord: float):
    """Garante que coordenadas corrompidas (NaN, Inf, fora do intervalo) são rejeitadas."""
    bbox = BoundingBox(class_id=0, x_center=invalid_coord, y_center=0.5, w=0.2, h=0.2)
    assert bbox.is_valid_yolo() is False


@pytest.mark.property
@given(
    w=st.floats(min_value=0.05, max_value=0.5),
    h=st.floats(min_value=0.05, max_value=0.5),
)
def test_property_iou_symmetry_and_range(w: float, h: float):
    """Propriedade fundamental: IoU(A, B) == IoU(B, A) e 0.0 <= IoU <= 1.0."""
    b1 = BoundingBox(class_id=0, x_center=0.4, y_center=0.4, w=w, h=h)
    b2 = BoundingBox(class_id=0, x_center=0.5, y_center=0.5, w=w, h=h)

    iou1 = b1.iou(b2)
    iou2 = b2.iou(b1)

    assert pytest.approx(iou1, 1e-6) == iou2
    assert 0.0 <= iou1 <= 1.0
