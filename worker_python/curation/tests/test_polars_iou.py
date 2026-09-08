"""Unit tests for vectorized IoU calculation in Polars expressions."""
import pytest
import polars as pl
from curation.polars_iou import compute_vectorized_iou


@pytest.mark.unit
def test_compute_vectorized_iou_exact_and_disjoint():
    df = pl.DataFrame({
        # Par 1: Caixas idênticas (IoU = 1.0)
        # Par 2: Caixas disjuntas (IoU = 0.0)
        # Par 3: Sobreposição parcial 50%
        "a_x_center": [0.5, 0.2, 0.5],
        "a_y_center": [0.5, 0.2, 0.5],
        "a_w": [0.2, 0.1, 0.2],
        "a_h": [0.2, 0.1, 0.2],
        "b_x_center": [0.5, 0.8, 0.55],
        "b_y_center": [0.5, 0.8, 0.5],
        "b_w": [0.2, 0.1, 0.2],
        "b_h": [0.2, 0.1, 0.2],
    }).lazy()

    result = compute_vectorized_iou(df).collect()

    assert "iou" in result.columns
    ious = result["iou"].to_list()

    # Par 1 deve ser 1.0
    assert pytest.approx(ious[0], 1e-4) == 1.0
    # Par 2 deve ser 0.0
    assert pytest.approx(ious[1], 1e-4) == 0.0
    # Par 3 deve ser aproximadamente 0.60
    assert 0.50 <= ious[2] <= 0.70
