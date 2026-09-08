"""Unit tests for rolling windows and anomaly burst detection."""
import pytest
import polars as pl
from curation.polars_dynamic import detect_incident_bursts


@pytest.mark.unit
def test_detect_incident_bursts():
    # Cria uma sequência de 50 detecções concentradas em 20 segundos
    timestamps = pl.datetime_range(
        pl.datetime(2026, 9, 8, 14, 0, 0),
        pl.datetime(2026, 9, 8, 14, 0, 20),
        interval="400ms",
        eager=True,
    )
    df = pl.DataFrame({
        "timestamp": timestamps,
        "camera_id": ["cam_01"] * len(timestamps),
        "frame_id": range(len(timestamps)),
        "confidence": [0.9] * len(timestamps),
        "class_id": [0] * len(timestamps),
    }).lazy()

    # Janela deslizante de 1 minuto a cada 10 segundos buscando no mínimo 25 detecções
    bursts = detect_incident_bursts(
        df, every="10s", period="1m", min_detections=25, target_class_id=0
    ).collect()

    # Deve detectar a explosão de atividade
    assert len(bursts) >= 1
    assert bursts["detection_count"].max() >= 25
    assert "target_class_count" in bursts.columns
