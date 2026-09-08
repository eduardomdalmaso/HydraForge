"""Optimized Apache Arrow and Polars schema definitions for HydraForge metadata."""
from __future__ import annotations
import polars as pl

# 1. Schema para struct de uma única detecção dentro do frame
DETECTION_STRUCT_SCHEMA = pl.Struct({
    "class_id": pl.UInt16,
    "confidence": pl.Float32,
    "bbox": pl.Array(pl.Float32, shape=4),       # [x_center, y_center, w, h] - Comprimento fixo -> SIMD
    "segmentation": pl.List(pl.Float32),        # Polígonos de vértices variáveis
})

# 2. Schema principal colunar para anotações brutas por frame
RAW_ANNOTATIONS_SCHEMA = {
    "frame_id": pl.UInt64,
    "camera_id": pl.Categorical,
    "timestamp": pl.Datetime(time_unit="us"),
    "image_path": pl.String,
    "frame_width": pl.UInt16,
    "frame_height": pl.UInt16,
    "phash": pl.UInt64,
    "detections": pl.List(DETECTION_STRUCT_SCHEMA),
}

# 3. Schema plano pós-desaninhamento (Ready for Curation & Splitting)
FLAT_DETECTION_SCHEMA = {
    "frame_id": pl.UInt64,
    "camera_id": pl.Categorical,
    "timestamp": pl.Datetime(time_unit="us"),
    "image_path": pl.String,
    "frame_width": pl.UInt16,
    "frame_height": pl.UInt16,
    "phash": pl.UInt64,
    "class_id": pl.UInt16,
    "confidence": pl.Float32,
    "x_center": pl.Float32,
    "y_center": pl.Float32,
    "w": pl.Float32,
    "h": pl.Float32,
    "segmentation": pl.List(pl.Float32),
}


def ingest_hydra_metadata_lazy(
    parquet_glob_path: str,
    low_memory: bool = True,
    parallel: str = "prefiltered",
) -> pl.LazyFrame:
    """Initiates lazy scan on Parquet files, applying predicate and projection pushdowns directly at disk scan."""
    return pl.scan_parquet(
        source=parquet_glob_path,
        parallel=parallel,
        low_memory=low_memory,
    )
