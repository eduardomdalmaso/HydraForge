from datetime import datetime
import pytest
import polars as pl
from curation.polars_asof import sync_frames_with_telemetry


@pytest.mark.unit
def test_sync_frames_with_telemetry():
    # 3 frames de vídeo com timestamps em microssegundos
    frames_df = pl.DataFrame({
        "timestamp": pl.datetime_range(
            datetime(2026, 9, 8, 10, 0, 0),
            datetime(2026, 9, 8, 10, 0, 0, 100000),
            interval="50ms",
            eager=True,
        ),
        "camera_id": ["cam_01", "cam_01", "cam_01"],
        "frame_id": [1, 2, 3],
    }).lazy()

    # Sensor de GPS emitindo com pequeno descasamento temporal
    telemetry_df = pl.DataFrame({
        "timestamp": [
            datetime(2026, 9, 8, 10, 0, 0, 2000),   # +2ms do frame 1
            datetime(2026, 9, 8, 10, 0, 0, 51000),  # +1ms do frame 2
            datetime(2026, 9, 8, 10, 0, 0, 103000), # +3ms do frame 3
        ],
        "camera_id": ["cam_01", "cam_01", "cam_01"],
        "altitude_m": [150.2, 150.8, 151.1],
        "speed_kmh": [45.0, 46.2, 47.0],
    }).lazy()

    synced = sync_frames_with_telemetry(
        frames_df, telemetry_df, tolerance="10ms", strategy="nearest"
    ).collect()

    assert len(synced) == 3
    assert "altitude_m" in synced.columns
    assert synced["altitude_m"].to_list() == [150.2, 150.8, 151.1]
