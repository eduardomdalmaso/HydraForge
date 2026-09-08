"""Unit tests for TrainingProfiler I/O and GPU metrics."""
import time
import pytest
from training.profiler_tool import TrainingProfiler, ProfilingMetrics


@pytest.mark.unit
def test_training_profiler_metrics_calculation():
    profiler = TrainingProfiler(warm_up_steps=1)

    # Simula 3 passos de treino
    for i in range(3):
        # Simula tempo de espera por dados (10ms)
        time.sleep(0.01)
        profiler.record_data_loaded()

        # Simula tempo de cálculo da GPU (20ms)
        time.sleep(0.02)
        profiler.record_step_completed()

    metrics = profiler.get_summary()

    assert isinstance(metrics, ProfilingMetrics)
    assert metrics.total_steps == 2  # 3 passos - 1 warmup = 2 passos medidos
    assert metrics.avg_compute_ms >= 15.0
    assert metrics.avg_data_wait_ms >= 5.0
    assert 0.0 <= metrics.io_wait_ratio_pct <= 100.0


@pytest.mark.unit
def test_training_profiler_empty():
    profiler = TrainingProfiler()
    metrics = profiler.get_summary()
    assert metrics.total_steps == 0
    assert metrics.avg_step_duration_ms == 0.0
    assert metrics.is_io_starved is False
