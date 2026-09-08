"""Stratified train/val/test splitting and Hive-partitioned Parquet/YAML export."""
from __future__ import annotations
from pathlib import Path
from typing import Dict, List, Sequence
import polars as pl
import yaml


def stratified_train_val_split(
    df: pl.LazyFrame,
    train_ratio: float = 0.70,
    val_ratio: float = 0.15,
    seed: int = 42,
    group_col: str | None = None,
    temporal_window: str | None = None,  # Ex: "30s" para agrupar cenas contíguas
    class_col: str = "class_id",
) -> pl.LazyFrame:
    """Performs stratified split across classes while grouping sequential frames or temporal scenes to prevent data leakage."""
    val_boundary = train_ratio + val_ratio

    # 1. Se temporal_window foi fornecido, cria um identificador de cena única (camera_id + timestamp truncado)
    if temporal_window is not None:
        group_expr = pl.concat_str([
            pl.col("camera_id").cast(pl.String),
            pl.lit("_scene_"),
            pl.col("timestamp").dt.truncate(temporal_window).cast(pl.String)
        ])
    elif group_col is not None:
        group_expr = pl.col(group_col)
    else:
        group_expr = pl.col("frame_id")

    return (
        df.with_columns([
            # Hash determinístico do grupo/cena para garantir que a cena inteira vá para o mesmo split
            group_expr.hash(seed=seed).alias("_group_hash")
        ])
        .with_columns([
            # Ranking percentual do bloco dentro de cada classe para balanceamento exato
            (pl.col("_group_hash").rank(method="dense") / pl.col("_group_hash").count())
            .over(class_col)
            .alias("_stratified_ratio")
        ])
        .with_columns([
            pl.when(pl.col("_stratified_ratio") <= train_ratio)
            .then(pl.lit("train"))
            .when(pl.col("_stratified_ratio") <= val_boundary)
            .then(pl.lit("val"))
            .otherwise(pl.lit("test"))
            .alias("split_destination")
        ])
        .drop(["_group_hash", "_stratified_ratio"])
    )


def export_partitioned_splits(
    df: pl.DataFrame | pl.LazyFrame,
    output_dir: str | Path,
    partition_by: Sequence[str] = ("split_destination",),
    compression: str = "zstd",
    compression_level: int = 3,
) -> None:
    """Exports structured DataFrame/LazyFrame into native Hive-partitioned directory tree using multithreaded Rust writer."""
    out_path = Path(output_dir)
    out_path.mkdir(parents=True, exist_ok=True)

    if isinstance(df, pl.LazyFrame):
        materialized_df = df.collect()
    else:
        materialized_df = df

    materialized_df.write_parquet(
        str(out_path),
        use_pyarrow=False,  # Native Polars Rust writer
        compression=compression,
        compression_level=compression_level,
        partition_by=list(partition_by),
    )


def generate_data_yaml(
    output_dir: str | Path,
    class_map: Dict[int, str] | List[str],
    dataset_name: str = "hydra_dataset",
) -> Path:
    """Generates standardized Ultralytics YOLO data.yaml configuration from splits."""
    out_path = Path(output_dir)
    out_path.mkdir(parents=True, exist_ok=True)

    if isinstance(class_map, list):
        names_dict = {i: name for i, name in enumerate(class_map)}
    else:
        names_dict = class_map

    yaml_content = {
        "path": str(out_path.resolve()),
        "train": "split_destination=train",
        "val": "split_destination=val",
        "test": "split_destination=test",
        "names": names_dict,
        "nc": len(names_dict),
    }

    yaml_path = out_path / "data.yaml"
    with open(yaml_path, "w", encoding="utf-8") as f:
        yaml.dump(yaml_content, f, sort_keys=False)

    return yaml_path
