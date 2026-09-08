"""Unit tests for Zero-Copy Polars to PyTorch bridge with Pinned Memory."""
import pytest
import polars as pl
import torch
from curation.polars_bridge import (
    polars_column_to_tensor,
    polars_matrix_to_tensor,
    transfer_to_device_async,
)


@pytest.mark.unit
def test_polars_column_to_tensor_zero_copy():
    series = pl.Series("confidence", [0.95, 0.88, 0.72, 0.99], dtype=pl.Float32)
    tensor = polars_column_to_tensor(series)

    assert isinstance(tensor, torch.Tensor)
    assert tensor.dtype == torch.float32
    assert tensor.shape == (4,)
    assert torch.allclose(tensor, torch.tensor([0.95, 0.88, 0.72, 0.99], dtype=torch.float32))


@pytest.mark.unit
def test_polars_column_to_tensor_pinned_memory():
    series = pl.Series("confidence", [0.95, 0.88], dtype=pl.Float32)
    tensor = polars_column_to_tensor(series, pin_memory=True)

    assert isinstance(tensor, torch.Tensor)
    if torch.cuda.is_available():
        assert tensor.is_pinned()


@pytest.mark.unit
def test_polars_matrix_to_tensor():
    df = pl.DataFrame({
        "x_center": [0.5, 0.2],
        "y_center": [0.5, 0.3],
        "w": [0.2, 0.1],
        "h": [0.2, 0.1],
    }, schema={
        "x_center": pl.Float32,
        "y_center": pl.Float32,
        "w": pl.Float32,
        "h": pl.Float32,
    })

    tensor = polars_matrix_to_tensor(df, ["x_center", "y_center", "w", "h"], pin_memory=True)
    assert isinstance(tensor, torch.Tensor)
    assert tensor.shape == (2, 4)
    assert tensor.dtype == torch.float32
    assert tensor[0, 0].item() == pytest.approx(0.5, 1e-4)


@pytest.mark.unit
def test_transfer_to_device_async():
    tensor = torch.randn(4, 4, dtype=torch.float32)
    device = "cuda" if torch.cuda.is_available() else "cpu"
    transferred = transfer_to_device_async(tensor, device=device, non_blocking=True)

    assert isinstance(transferred, torch.Tensor)
    assert transferred.device.type == device
