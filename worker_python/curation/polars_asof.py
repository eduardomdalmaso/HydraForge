"""Asynchronous temporal synchronization for video frames and external sensor telemetry (GPS, PTZ, Drones)."""
from __future__ import annotations
import polars as pl


def sync_frames_with_telemetry(
    frames_df: pl.LazyFrame,
    telemetry_df: pl.LazyFrame,
    time_col: str = "timestamp",
    by_col: str | None = "camera_id",
    tolerance: str = "50ms",
    strategy: str = "nearest",
    already_sorted: bool = False,
) -> pl.LazyFrame:
    """Performs asof join to match video frames with nearest telemetry measurements within tolerance.
    
    Optimized based on High Performance Python: if `already_sorted=True`, sets monotonic sorted flag
    directly, avoiding O(N log N) sorting overhead.
    """
    if already_sorted:
        sorted_frames = frames_df.with_columns(pl.col(time_col).set_sorted())
        sorted_telemetry = telemetry_df.with_columns(pl.col(time_col).set_sorted())
    else:
        sorted_frames = frames_df.sort(time_col)
        sorted_telemetry = telemetry_df.sort(time_col)

    return sorted_frames.join_asof(
        sorted_telemetry,
        on=time_col,
        by=by_col,
        tolerance=tolerance,
        strategy=strategy,
    )
