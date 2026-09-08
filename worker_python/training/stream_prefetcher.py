"""Asynchronous CUDA Stream Prefetcher for zero-wait PCIe transfer and Tensor Core overlapping."""
from __future__ import annotations
from typing import Any, Iterator, Optional
import torch


class CUDAPrefetcher:
    """Wraps a PyTorch DataLoader to prefetch batches onto the GPU using a dedicated CUDA Stream.
    
    Achieves true compute/transfer overlap based on High Performance Python Chapter 9 & 12.
    """

    def __init__(
        self,
        loader: Any,
        device: torch.device | str = "cuda",
    ) -> None:
        self.loader = loader
        if isinstance(device, str):
            device = torch.device(device)
        self.device = device
        self.is_cuda = (self.device.type == "cuda") and torch.cuda.is_available()
        self.stream = torch.cuda.Stream(device=self.device) if self.is_cuda else None
        self.iterator: Optional[Iterator[Any]] = None
        self.next_batch: Any = None

    def _preload(self) -> None:
        try:
            batch = next(self.iterator)
        except StopIteration:
            self.next_batch = None
            return

        if not self.is_cuda:
            self.next_batch = batch
            return

        with torch.cuda.stream(self.stream):
            self.next_batch = self._to_device_recursive(batch)

    def _to_device_recursive(self, item: Any) -> Any:
        if isinstance(item, torch.Tensor):
            if not item.is_pinned() and not item.is_cuda:
                item = item.pin_memory()
            return item.to(device=self.device, non_blocking=True)
        elif isinstance(item, (list, tuple)):
            return type(item)(self._to_device_recursive(x) for x in item)
        elif isinstance(item, dict):
            return {k: self._to_device_recursive(v) for k, v in item.items()}
        return item

    def __iter__(self) -> CUDAPrefetcher:
        self.iterator = iter(self.loader)
        self._preload()
        return self

    def __next__(self) -> Any:
        if self.is_cuda:
            torch.cuda.current_stream().wait_stream(self.stream)

        batch = self.next_batch
        if batch is None:
            raise StopIteration

        # Record stream on tensors to prevent premature CUDA memory pool reclamation
        if self.is_cuda:
            self._record_stream_recursive(batch, torch.cuda.current_stream())

        self._preload()
        return batch

    def _record_stream_recursive(self, item: Any, stream: torch.cuda.Stream) -> None:
        if isinstance(item, torch.Tensor):
            item.record_stream(stream)
        elif isinstance(item, (list, tuple)):
            for x in item:
                self._record_stream_recursive(x, stream)
        elif isinstance(item, dict):
            for v in item.values():
                self._record_stream_recursive(v, stream)

    def __len__(self) -> int:
        return len(self.loader)
