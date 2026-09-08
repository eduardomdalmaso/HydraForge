"""Unit tests for CUDAPrefetcher overlapping."""
import pytest
import torch
from torch.utils.data import DataLoader, TensorDataset
from training.stream_prefetcher import CUDAPrefetcher


@pytest.mark.unit
def test_cuda_prefetcher_cpu_fallback():
    # Cria dataset simples com 10 tensores em lotes de 2
    x = torch.randn(10, 4)
    y = torch.randint(0, 2, (10,))
    dataset = TensorDataset(x, y)
    loader = DataLoader(dataset, batch_size=2)

    prefetcher = CUDAPrefetcher(loader, device="cpu")
    assert len(prefetcher) == 5

    batches_seen = 0
    for batch_x, batch_y in prefetcher:
        assert batch_x.shape == (2, 4)
        assert batch_y.shape == (2,)
        batches_seen += 1

    assert batches_seen == 5


@pytest.mark.smoke
def test_cuda_prefetcher_gpu_execution():
    if not torch.cuda.is_available():
        pytest.skip("CUDA não disponível no ambiente atual")

    x = torch.randn(8, 3, 32, 32)
    y = torch.randint(0, 2, (8,))
    dataset = TensorDataset(x, y)
    loader = DataLoader(dataset, batch_size=4, pin_memory=True)

    prefetcher = CUDAPrefetcher(loader, device="cuda")
    assert len(prefetcher) == 2

    for batch_x, batch_y in prefetcher:
        assert batch_x.is_cuda
        assert batch_y.is_cuda
