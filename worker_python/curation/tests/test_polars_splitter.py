"""Unit tests for stratified splitting and Hive Parquet / YAML export."""
import tempfile
from pathlib import Path
import pytest
import polars as pl
import yaml
from curation.polars_splitter import (
    stratified_train_val_split,
    export_partitioned_splits,
    generate_data_yaml,
)


@pytest.mark.unit
def test_stratified_train_val_split():
    # Dataset balanceado de 1000 linhas com 2 classes
    df = pl.DataFrame({
        "frame_id": range(1000),
        "class_id": [0, 1] * 500,
    }).lazy()

    split_df = stratified_train_val_split(df, train_ratio=0.70, val_ratio=0.15, seed=42).collect()

    assert "split_destination" in split_df.columns
    counts = split_df["split_destination"].value_counts().to_dicts()
    counts_map = {item["split_destination"]: item["count"] for item in counts}

    # Aproximações de proporção (70% train, 15% val, 15% test)
    assert counts_map["train"] == pytest.approx(700, abs=40)
    assert counts_map["val"] == pytest.approx(150, abs=30)
    assert counts_map["test"] == pytest.approx(150, abs=30)


@pytest.mark.unit
def test_export_partitioned_splits_and_data_yaml():
    with tempfile.TemporaryDirectory() as tmp_dir:
        tmp_path = Path(tmp_dir)

        df = pl.DataFrame({
            "frame_id": [1, 2, 3],
            "class_id": [0, 1, 0],
            "split_destination": ["train", "val", "train"],
        })

        # Exporta particionado
        export_partitioned_splits(df, tmp_path / "splits", partition_by=["split_destination"])

        # Verifica criação dos diretórios Hive
        train_dir = tmp_path / "splits" / "split_destination=train"
        val_dir = tmp_path / "splits" / "split_destination=val"
        assert train_dir.exists()
        assert val_dir.exists()

        # Gera data.yaml
        class_map = {0: "person", 1: "car"}
        yaml_path = generate_data_yaml(tmp_path, class_map)

        assert yaml_path.exists()
        with open(yaml_path, "r", encoding="utf-8") as f:
            data = yaml.safe_load(f)

        assert data["nc"] == 2
        assert data["names"] == {0: "person", 1: "car"}
        assert data["train"] == "split_destination=train"


@pytest.mark.unit
def test_stratified_split_with_temporal_window():
    # 10 frames de cam_01 dentro do mesmo minuto (mesma cena)
    timestamps = pl.datetime_range(
        pl.datetime(2026, 9, 8, 12, 0, 0),
        pl.datetime(2026, 9, 8, 12, 0, 9),
        interval="1s",
        eager=True,
    )
    df = pl.DataFrame({
        "frame_id": range(10),
        "camera_id": ["cam_01"] * 10,
        "timestamp": timestamps,
        "class_id": [0] * 10,
    }).lazy()

    # Janela temporal de 30s deve agrupar todos os 10 frames na mesma cena
    split_df = stratified_train_val_split(df, temporal_window="30s", seed=42).collect()

    # Todos os 10 frames da cena devem ir para o mesmo split_destination
    distinct_destinations = split_df["split_destination"].n_unique()
    assert distinct_destinations == 1, "Todos os frames da mesma cena temporal devem ficar no mesmo split (sem vazamento)"

