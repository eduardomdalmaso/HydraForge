"""Rolling windows and dynamic time aggregation for Active Learning incident burst detection."""
from __future__ import annotations
import polars as pl


def detect_incident_bursts(
    df: pl.LazyFrame,
    time_col: str = "timestamp",
    by_col: str = "camera_id",
    every: str = "10s",
    period: str = "1m",
    min_detections: int = 20,
    target_class_id: int | None = None,
) -> pl.LazyFrame:
    """Uses dynamic rolling windows to identify time intervals with abnormal detection densities."""
    sorted_df = df.sort(time_col)

    agg_exprs = [
        pl.len().alias("detection_count"),
        pl.col("frame_id").n_unique().alias("unique_frames_in_burst"),
        pl.col("confidence").mean().alias("avg_confidence"),
    ]

    if target_class_id is not None:
        agg_exprs.append(
            pl.col("class_id").filter(pl.col("class_id") == target_class_id).count().alias("target_class_count")
        )

    query = (
        sorted_df.group_by_dynamic(
            time_col,
            every=every,
            period=period,
            group_by=by_col,
        )
        .agg(agg_exprs)
        .filter(pl.col("detection_count") >= min_detections)
    )

    return query
