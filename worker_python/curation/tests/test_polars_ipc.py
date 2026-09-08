"""Unit tests for POSIX /dev/shm Arrow IPC streaming."""
import tempfile
from pathlib import Path
import pytest
import polars as pl
from curation.polars_ipc import sink_to_shm_ipc, read_from_shm_ipc


@pytest.mark.unit
def test_sink_and_read_shm_ipc():
    with tempfile.TemporaryDirectory() as tmp_dir:
        ipc_file = Path(tmp_dir) / "test_stream.arrow"

        orig_df = pl.DataFrame({
            "frame_id": [101, 102, 103],
            "camera_id": ["cam_alpha", "cam_beta", "cam_alpha"],
            "confidence": [0.85, 0.92, 0.78],
        }, schema={
            "frame_id": pl.UInt64,
            "camera_id": pl.Categorical,
            "confidence": pl.Float32,
        })

        sink_path = sink_to_shm_ipc(orig_df, ipc_file)
        assert sink_path.exists()

        recovered_df = read_from_shm_ipc(sink_path)
        assert len(recovered_df) == 3
        assert recovered_df["frame_id"].to_list() == [101, 102, 103]
        assert recovered_df["confidence"].to_list() == [pytest.approx(0.85, 1e-4), pytest.approx(0.92, 1e-4), pytest.approx(0.78, 1e-4)]
