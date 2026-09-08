"""Unit and property tests for Polars-based curation, flattening, and deduplication."""
import pytest
import polars as pl
from curation.polars_curator import (
    flatten_nested_detections,
    deduplicate_consecutive_frames,
    filter_and_audit_annotations,
    compute_dataset_telemetry,
)


@pytest.mark.unit
def test_flatten_nested_detections():
    # Cria DataFrame com lista de structs e bboxes de tamanho fixo 4
    raw_df = pl.DataFrame({
        "frame_id": [101, 102],
        "camera_id": ["cam_01", "cam_02"],
        "phash": [12345, 67890],
        "detections": [
            [
                {"class_id": 0, "confidence": 0.9, "bbox": [0.5, 0.5, 0.2, 0.2], "segmentation": [0.5, 0.5]},
                {"class_id": 1, "confidence": 0.8, "bbox": [0.2, 0.3, 0.1, 0.1], "segmentation": []},
            ],
            [
                {"class_id": 2, "confidence": 0.95, "bbox": [0.8, 0.8, 0.1, 0.1], "segmentation": [0.8, 0.8]},
            ],
        ],
    }, schema={
        "frame_id": pl.UInt64,
        "camera_id": pl.String,
        "phash": pl.UInt64,
        "detections": pl.List(pl.Struct({
            "class_id": pl.UInt16,
            "confidence": pl.Float32,
            "bbox": pl.Array(pl.Float32, shape=4),
            "segmentation": pl.List(pl.Float32),
        })),
    })

    flat_df = flatten_nested_detections(raw_df.lazy()).collect()

    # Deve ter expandido para 3 linhas de detecções
    assert len(flat_df) == 3
    assert set(flat_df.columns) >= {"frame_id", "camera_id", "class_id", "confidence", "x_center", "y_center", "w", "h"}
    assert flat_df.filter(pl.col("frame_id") == 101).shape[0] == 2
    assert flat_df["w"].to_list() == [pytest.approx(0.2, 1e-4), pytest.approx(0.1, 1e-4), pytest.approx(0.1, 1e-4)]


@pytest.mark.unit
def test_deduplicate_consecutive_frames():
    df = pl.DataFrame({
        "camera_id": ["cam_01", "cam_01", "cam_01", "cam_02", "cam_02"],
        "frame_id": [1, 2, 3, 4, 5],
        "phash": [0xAAAA, 0xAAAA, 0xBBBB, 0xCCCC, 0xCCCC],
    }).lazy()

    dedup = deduplicate_consecutive_frames(df).collect()

    # cam_01 tem 1 duplicado (frame 2) e cam_02 tem 1 duplicado (frame 5) -> sobram 3 frames
    assert len(dedup) == 3
    assert dedup["frame_id"].to_list() == [1, 3, 4]


@pytest.mark.unit
def test_filter_and_audit_annotations():
    df = pl.DataFrame({
        "frame_id": [1, 2, 3, 4, 5],
        "class_id": [0, 0, 1, 99, 0],
        "confidence": [0.90, 0.10, 0.85, 0.90, 0.95],
        "x_center": [0.5, 0.5, 0.5, 0.5, 1.5], # 1.5 está fora do limite
        "y_center": [0.5, 0.5, 0.5, 0.5, 0.5],
        "w": [0.2, 0.2, 0.2, 0.2, 0.2],
        "h": [0.2, 0.2, 0.2, 0.2, 0.2],
    }).lazy()

    # Apenas classes 0 e 1, confiança >= 0.5
    filtered = filter_and_audit_annotations(df, min_confidence=0.5, allowed_classes={0, 1}).collect()

    # Apenas o frame 1 e 3 devem sobreviver
    assert len(filtered) == 2
    assert filtered["frame_id"].to_list() == [1, 3]


@pytest.mark.unit
def test_compute_dataset_telemetry():
    df = pl.DataFrame({
        "camera_id": ["cam_01", "cam_01", "cam_02"],
        "frame_id": [1, 1, 2],
        "class_id": [0, 1, 0],
        "confidence": [0.8, 0.9, 0.95],
        "w": [0.2, 0.3, 0.4],
        "h": [0.2, 0.3, 0.4],
    }).lazy()

    stats = compute_dataset_telemetry(df).collect()
    cam1 = stats.filter(pl.col("camera_id") == "cam_01").to_dicts()[0]
    assert cam1["total_annotations"] == 2
    assert cam1["unique_frames"] == 1
    assert cam1["distinct_classes"] == 2
