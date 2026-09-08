"""
Unit and Behavior Tests for HydraForge MemoryBudgetGuard.
"""

import pytest
from training.memory_guard import MemoryBudgetGuard, MemoryPlan


@pytest.mark.unit
def test_ram_estimation_scaling():
    # 1.000 imagens 640x640 sem cache
    ram_1k = MemoryBudgetGuard.estimate_ram_usage(
        dataset_size=1000,
        imgsz=640,
        batch_size=32,
        num_workers=4,
        cache=False
    )
    # Deve ser em torno de ~1.5 GB
    assert 1e9 < ram_1k < 3e9

    # 100.000 imagens COM cache=True em RAM
    ram_100k_cached = MemoryBudgetGuard.estimate_ram_usage(
        dataset_size=100_000,
        imgsz=640,
        batch_size=32,
        num_workers=4,
        cache=True
    )
    # Deve ultrapassar 100 GB
    assert ram_100k_cached > 100 * (1024**3)


@pytest.mark.unit
def test_anti_oom_guard_disables_cache_on_large_dataset():
    # Simula máquina com 32GB de RAM livre
    simulated_ram = 32 * (1024**3)
    simulated_vram = 32 * (1024**3)

    plan = MemoryBudgetGuard.compute_safe_plan(
        dataset_size=80_000,  # 80k imagens exigiria ~98GB
        imgsz=640,
        batch_size=32,
        num_workers=8,
        cache='ram',
        override_available_ram=simulated_ram,
        override_available_vram=simulated_vram
    )

    # O guardião deve ter forçado cache para False
    assert plan.cache is False
    assert plan.is_adjusted is True
    assert any("[RAM_GUARD] 'cache=ram' desativado" in w for w in plan.warnings)


@pytest.mark.unit
def test_anti_oom_guard_scales_down_workers_on_low_memory():
    # Simula máquina com apenas 4GB de RAM livre
    low_ram = 4 * (1024**3)
    plan = MemoryBudgetGuard.compute_safe_plan(
        dataset_size=10_000,
        imgsz=1280,  # Imagens grandes
        batch_size=32,
        num_workers=16,  # 16 workers consumiria toda a RAM
        cache=False,
        override_available_ram=low_ram
    )

    assert plan.num_workers < 16
    assert plan.is_adjusted is True
    assert plan.estimated_ram_mb <= (low_ram * 0.80) / (1024**2)


@pytest.mark.unit
def test_anti_oom_guard_accepts_healthy_parameters():
    # Máquina com 64GB de RAM e 32GB VRAM rodando dataset normal
    plan = MemoryBudgetGuard.compute_safe_plan(
        dataset_size=5_000,
        imgsz=640,
        batch_size=32,
        num_workers=8,
        cache=False,
        override_available_ram=64 * (1024**3),
        override_available_vram=32 * (1024**3)
    )

    assert plan.is_adjusted is False
    assert plan.batch_size == 32
    assert plan.num_workers == 8
    assert len(plan.warnings) == 0
