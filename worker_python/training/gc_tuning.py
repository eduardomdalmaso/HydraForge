"""Garbage Collection tuning utilities based on High Performance Python Chapter 11."""
from __future__ import annotations
import gc
from contextlib import contextmanager
from typing import Iterator


@contextmanager
def disable_gc_context() -> Iterator[None]:
    """Disables automatic garbage collection during high-frequency compute loops.
    
    Prevents stop-the-world micro-stutters during tensor training iterations.
    Re-enables and triggers a single manual collection on exit.
    """
    is_initially_enabled = gc.isenabled()
    if is_initially_enabled:
        gc.disable()
    try:
        yield
    finally:
        if is_initially_enabled:
            gc.enable()
            gc.collect(generation=1)


def tune_gc_thresholds(alloc_threshold: int = 50000) -> None:
    """Tunes GC allocation thresholds to trigger collections less frequently during training."""
    # Default is typically (700, 10, 10)
    current = gc.get_threshold()
    gc.set_threshold(alloc_threshold, current[1] * 2, current[2] * 2)
