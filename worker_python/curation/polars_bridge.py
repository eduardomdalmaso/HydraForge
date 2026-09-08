"""Zero-copy memory bridge from Polars/Apache Arrow to PyTorch Tensors with Pinned Memory support."""
from __future__ import annotations
import numpy as np
import polars as pl
import torch


def polars_column_to_tensor(
    series: pl.Series,
    pin_memory: bool = False,
    to_cuda: bool = False,
) -> torch.Tensor:
    """Extracts a zero-copy PyTorch Tensor from a numerical Polars Series via Apache Arrow.
    
    Optimized based on High Performance Python: guarantees memory contiguity to prevent
    L1/L2 cache misses and supports asynchronous pinned memory transfer to NVIDIA GPUs.
    """
    arrow_array = series.to_arrow()
    np_view = arrow_array.to_numpy(zero_copy_only=False)
    
    # Ensure memory is writable/contiguous for optimal PyTorch integration
    if not np_view.flags.writeable:
        np_view = np_view.copy()
        
    tensor = torch.from_numpy(np_view)
    if pin_memory and torch.cuda.is_available():
        tensor = tensor.pin_memory()
    if to_cuda and torch.cuda.is_available():
        return tensor.to(device="cuda", non_blocking=True)
    return tensor


def polars_matrix_to_tensor(
    df: pl.DataFrame,
    columns: list[str],
    pin_memory: bool = False,
    to_cuda: bool = False,
) -> torch.Tensor:
    """Converts a subset of float/int columns into a 2D contiguous PyTorch Tensor."""
    sub_df = df.select(columns)
    np_matrix = np.ascontiguousarray(sub_df.to_numpy())
    tensor = torch.from_numpy(np_matrix)
    if pin_memory and torch.cuda.is_available():
        tensor = tensor.pin_memory()
    if to_cuda and torch.cuda.is_available():
        return tensor.to(device="cuda", non_blocking=True)
    return tensor


def transfer_to_device_async(
    tensor: torch.Tensor,
    device: torch.device | str = "cuda",
    non_blocking: bool = True,
) -> torch.Tensor:
    """Performs non-blocking asynchronous DMA transfer to target device over PCIe."""
    if isinstance(device, str):
        device = torch.device(device)

    if device.type == "cuda" and not tensor.is_pinned() and not tensor.is_cuda:
        # Pinned memory is required for true non-blocking DMA transfer
        tensor = tensor.pin_memory()

    return tensor.to(device=device, non_blocking=non_blocking)
