"""GPU and I/O Profiling Tool based on High Performance Python Chapters 1, 2 & 9."""
from __future__ import annotations
import time
from dataclasses import dataclass
from typing import Any, Callable, Dict, List, Optional
import torch


@dataclass
class ProfilingMetrics:
    """Metrics capturing GPU compute vs Dataloader I/O efficiency."""
    total_steps: int
    avg_step_duration_ms: float
    avg_data_wait_ms: float
    avg_compute_ms: float
    io_wait_ratio_pct: float
    peak_vram_mb: float
    is_io_starved: bool  # True if data wait > 10% of total step time


class TrainingProfiler:
    """Lightweight step profiler measuring I/O wait vs Tensor Core compute duration."""

    def __init__(self, warm_up_steps: int = 2) -> None:
        self.warm_up_steps = warm_up_steps
        self.step_durations: List[float] = []
        self.data_wait_times: List[float] = []
        self.compute_times: List[float] = []
        self._last_step_end: Optional[float] = None
        self._batch_start: Optional[float] = None
        self.step_count = 0

    def record_data_loaded(self) -> None:
        """Call immediately after `batch = next(loader)`."""
        now = time.perf_counter()
        if self._last_step_end is not None and self.step_count >= self.warm_up_steps:
            wait_time = (now - self._last_step_end) * 1000.0  # ms
            self.data_wait_times.append(wait_time)
        self._batch_start = now

    def record_step_completed(self) -> None:
        """Call immediately after `optimizer.step()` and `scaler.update()`."""
        now = time.perf_counter()
        if self._batch_start is not None and self.step_count >= self.warm_up_steps:
            compute_time = (now - self._batch_start) * 1000.0  # ms
            self.compute_times.append(compute_time)
            total_step = compute_time + (self.data_wait_times[-1] if self.data_wait_times else 0.0)
            self.step_durations.append(total_step)

        self._last_step_end = now
        self.step_count += 1

    def get_summary(self) -> ProfilingMetrics:
        """Computes summary diagnostic metrics."""
        steps = len(self.compute_times)
        if steps == 0:
            return ProfilingMetrics(
                total_steps=0,
                avg_step_duration_ms=0.0,
                avg_data_wait_ms=0.0,
                avg_compute_ms=0.0,
                io_wait_ratio_pct=0.0,
                peak_vram_mb=0.0,
                is_io_starved=False,
            )

        avg_compute = sum(self.compute_times) / steps
        avg_wait = sum(self.data_wait_times) / len(self.data_wait_times) if self.data_wait_times else 0.0
        avg_total = avg_compute + avg_wait
        wait_ratio = (avg_wait / avg_total * 100.0) if avg_total > 0 else 0.0

        peak_vram = 0.0
        if torch.cuda.is_available():
            peak_vram = torch.cuda.max_memory_allocated() / (1024.0 * 1024.0)

        return ProfilingMetrics(
            total_steps=steps,
            avg_step_duration_ms=avg_total,
            avg_data_wait_ms=avg_wait,
            avg_compute_ms=avg_compute,
            io_wait_ratio_pct=wait_ratio,
            peak_vram_mb=peak_vram,
            is_io_starved=wait_ratio > 10.0,  # I/O starved if waiting > 10%
        )
