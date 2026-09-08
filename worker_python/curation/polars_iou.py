"""Vectorized Intersection over Union (IoU) calculation using native Polars expressions."""
from __future__ import annotations
import polars as pl


def compute_vectorized_iou(
    df: pl.LazyFrame,
    box_a_prefix: str = "a_",
    box_b_prefix: str = "b_",
    output_col: str = "iou",
) -> pl.LazyFrame:
    """Computes IoU between two sets of bounding boxes directly in Polars expressions.
    
    Expects columns:
      - {box_a_prefix}x_center, {box_a_prefix}y_center, {box_a_prefix}w, {box_a_prefix}h
      - {box_b_prefix}x_center, {box_b_prefix}y_center, {box_b_prefix}w, {box_b_prefix}h
    """
    ax = pl.col(f"{box_a_prefix}x_center")
    ay = pl.col(f"{box_a_prefix}y_center")
    aw = pl.col(f"{box_a_prefix}w")
    ah = pl.col(f"{box_a_prefix}h")

    bx = pl.col(f"{box_b_prefix}x_center")
    by = pl.col(f"{box_b_prefix}y_center")
    bw = pl.col(f"{box_b_prefix}w")
    bh = pl.col(f"{box_b_prefix}h")

    # Coordenadas mínimas e máximas de A
    a_xmin = ax - (aw / 2.0)
    a_xmax = ax + (aw / 2.0)
    a_ymin = ay - (ah / 2.0)
    a_ymax = ay + (ah / 2.0)

    # Coordenadas mínimas e máximas de B
    b_xmin = bx - (bw / 2.0)
    b_xmax = bx + (bw / 2.0)
    b_ymin = by - (bh / 2.0)
    b_ymax = by + (bh / 2.0)

    # Interseção
    inter_xmin = pl.max_horizontal(a_xmin, b_xmin)
    inter_ymin = pl.max_horizontal(a_ymin, b_ymin)
    inter_xmax = pl.min_horizontal(a_xmax, b_xmax)
    inter_ymax = pl.min_horizontal(a_ymax, b_ymax)

    inter_w = (inter_xmax - inter_xmin).clip(lower_bound=0.0)
    inter_h = (inter_ymax - inter_ymin).clip(lower_bound=0.0)
    intersection = inter_w * inter_h

    # Áreas e União
    area_a = aw * ah
    area_b = bw * bh
    union = area_a + area_b - intersection

    iou_expr = pl.when(union <= 0.0).then(0.0).otherwise(intersection / union).alias(output_col)

    return df.with_columns(iou_expr)
