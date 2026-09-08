"""High-performance curation, flattening, and deduplication engine using Polars expressions."""
from __future__ import annotations
import polars as pl


def flatten_nested_detections(df: pl.LazyFrame) -> pl.LazyFrame:
    """Flattens nested detection lists and structs into a flat, columnar representation.
    
    Transforms `detections: List[Struct]` into independent columns using native Rust `.explode()` and `.unnest()`.
    """
    return (
        df.explode("detections", empty_as_null=True)
        .unnest("detections")
        .with_columns([
            pl.col("bbox").arr.get(0).alias("x_center"),
            pl.col("bbox").arr.get(1).alias("y_center"),
            pl.col("bbox").arr.get(2).alias("w"),
            pl.col("bbox").arr.get(3).alias("h"),
        ])
        .drop("bbox")
    )


def deduplicate_consecutive_frames(df: pl.LazyFrame) -> pl.LazyFrame:
    """Removes identical consecutive frames from static cameras using window expressions.
    
    Applies `(phash == phash.shift(1)).over("camera_id")` to maintain embarrassingly parallel multi-core execution.
    """
    return (
        df.with_columns([
            (pl.col("phash") == pl.col("phash").shift(1).over("camera_id"))
            .fill_null(False)
            .alias("is_duplicate_frame")
        ])
        .filter(pl.col("is_duplicate_frame") == False)
        .drop("is_duplicate_frame")
    )


def filter_and_audit_annotations(
    df: pl.LazyFrame,
    min_confidence: float = 0.25,
    min_area: float = 0.0001,
    max_area: float = 0.98,
    allowed_classes: list[int] | set[int] | None = None,
) -> pl.LazyFrame:
    """Applies strict geometric and semantic sanity filters against corrupted or hallucinated annotations."""
    eps = 1e-6
    coords = ["x_center", "y_center", "w", "h"]

    # 1. Base validity: finite numbers, bounds in [0.0, 1.0], positive width and height
    is_finite = pl.all_horizontal([pl.col(c).is_finite() for c in coords])
    is_in_bounds = pl.all_horizontal([
        (pl.col(c) >= 0.0) & (pl.col(c) <= 1.0) for c in coords
    ])
    is_positive_dims = (pl.col("w") > 0.0) & (pl.col("h") > 0.0)

    # 2. Area boundaries
    area_expr = pl.col("w") * pl.col("h")
    is_valid_area = (area_expr >= min_area) & (area_expr <= max_area)

    # 3. Image boundary confinement
    x_min = pl.col("x_center") - (pl.col("w") / 2.0)
    x_max = pl.col("x_center") + (pl.col("w") / 2.0)
    y_min = pl.col("y_center") - (pl.col("h") / 2.0)
    y_max = pl.col("y_center") + (pl.col("h") / 2.0)

    is_inside_frame = (
        (x_min >= -eps)
        & (y_min >= -eps)
        & (x_max <= (1.0 + eps))
        & (y_max <= (1.0 + eps))
    )

    # 4. Confidence threshold
    is_confident = pl.col("confidence") >= min_confidence

    filtered = df.filter(
        is_finite & is_in_bounds & is_positive_dims & is_valid_area & is_inside_frame & is_confident
    )

    if allowed_classes is not None:
        filtered = filtered.filter(pl.col("class_id").is_in(list(allowed_classes)))

    return filtered


def compute_dataset_telemetry(df: pl.LazyFrame) -> pl.LazyFrame:
    """Computes instant aggregation statistics per camera (LazyFrame)."""
    return (
        df.group_by("camera_id")
        .agg([
            pl.len().alias("total_annotations"),
            pl.col("frame_id").n_unique().alias("unique_frames"),
            pl.col("class_id").n_unique().alias("distinct_classes"),
            pl.col("confidence").mean().alias("avg_confidence"),
            (pl.col("w") * pl.col("h")).mean().alias("avg_bbox_area"),
        ])
    )
