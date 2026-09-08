"""
HydraForge Memory Budget Guard (Pre-Flight Resource Estimation & Anti-OOM Guard)
Inspired by 'High Performance Python' and 'Polars' out-of-core memory principles.
"""

from dataclasses import dataclass
from typing import Optional, List, Tuple
import psutil
import torch


@dataclass(frozen=True, slots=True)
class MemoryPlan:
    """Plano de alocação de memória seguro validado antes do treino."""
    batch_size: int
    num_workers: int
    cache: str | bool
    estimated_ram_mb: float
    estimated_vram_mb: float
    available_ram_mb: float
    available_vram_mb: float
    warnings: Tuple[str, ...]
    is_adjusted: bool


class MemoryBudgetGuard:
    """
    Guardião de Orçamento de Memória e Auto-Cura de Parâmetros de Treino.
    Previne que o processo seja morto pelo Linux OOM Killer ao lidar com datasets gigantescos.
    """

    BYTES_PER_PIXEL_RGB = 3
    WORKER_BASE_OVERHEAD_BYTES = 150 * 1024 * 1024  # ~150 MB por worker Python
    MAIN_PROC_BASE_BYTES = 600 * 1024 * 1024        # ~600 MB PyTorch + Ultralytics base
    SAFETY_MARGIN_RAM = 0.80                        # Limite de segurança: 80% da RAM livre
    SAFETY_MARGIN_VRAM = 0.90                       # Limite de segurança: 90% da VRAM livre

    @classmethod
    def estimate_ram_usage(
        cls,
        dataset_size: int,
        imgsz: int,
        batch_size: int,
        num_workers: int,
        prefetch_factor: int = 2,
        cache: str | bool = False
    ) -> float:
        """Calcula o consumo estimado de RAM em Bytes."""
        frame_bytes = imgsz * imgsz * cls.BYTES_PER_PIXEL_RGB
        
        # 1. Custo de Cache em RAM (se o usuário ativar cache=True / cache='ram')
        cache_bytes = (dataset_size * frame_bytes) if (cache is True or cache == 'ram') else 0
        
        # 2. Custo dos Workers do DataLoader (In-flight Batches + Base Process)
        if num_workers > 0:
            in_flight_per_worker = prefetch_factor * batch_size * frame_bytes
            workers_bytes = num_workers * (cls.WORKER_BASE_OVERHEAD_BYTES + in_flight_per_worker)
        else:
            workers_bytes = 0
            
        total_bytes = cls.MAIN_PROC_BASE_BYTES + cache_bytes + workers_bytes
        return float(total_bytes)

    @classmethod
    def estimate_vram_usage(
        cls,
        imgsz: int,
        batch_size: int,
        amp: bool = True
    ) -> float:
        """Calcula o consumo estimado de VRAM em Bytes para o batch loop."""
        bytes_per_elem = 2 if amp else 4  # FP16 vs FP32
        # Estimativa empírica de ativações + gradientes YOLO (~12x tamanho da entrada do batch)
        batch_tensor_bytes = batch_size * (imgsz * imgsz * 3 * bytes_per_elem)
        activations_and_grads = batch_tensor_bytes * 12
        model_weights_bytes = 400 * 1024 * 1024  # ~400 MB pesos + buffers CUDA
        return float(model_weights_bytes + activations_and_grads)

    @classmethod
    def compute_safe_plan(
        cls,
        dataset_size: int,
        imgsz: int,
        batch_size: int,
        num_workers: int,
        cache: str | bool = False,
        amp: bool = True,
        override_available_ram: Optional[int] = None,
        override_available_vram: Optional[int] = None
    ) -> MemoryPlan:
        """
        Analisa os parâmetros de entrada contra os recursos do hardware e aplica
        auto-cura rebaixando flags perigosas antes que o treino inicie.
        """
        # Obter recursos do sistema
        if override_available_ram is not None:
            available_ram = override_available_ram
        else:
            available_ram = psutil.virtual_memory().available

        if override_available_vram is not None:
            available_vram = override_available_vram
        elif torch.cuda.is_available():
            try:
                free_b, total_b = torch.cuda.mem_get_info()
                available_vram = free_b
            except Exception:
                available_vram = 32 * 1024 * 1024 * 1024  # Fallback 32GB
        else:
            available_vram = 0

        max_allowed_ram = available_ram * cls.SAFETY_MARGIN_RAM
        warnings: List[str] = []
        is_adjusted = False

        adj_batch = batch_size
        adj_workers = num_workers
        adj_cache = cache

        # 1. Checagem Anti-OOM de Cache em RAM
        if adj_cache is True or adj_cache == 'ram':
            cache_cost = dataset_size * (imgsz * imgsz * cls.BYTES_PER_PIXEL_RGB)
            if cache_cost > max_allowed_ram * 0.5:  # Cache não pode ocupar mais de 50% do limite
                adj_cache = False
                is_adjusted = True
                warnings.append(
                    f"[RAM_GUARD] 'cache=ram' desativado automaticamente! "
                    f"Dataset ({dataset_size:,} fotos) exigiria {cache_cost / (1024**3):.2f} GB de RAM."
                )

        # 2. Checagem e Redução de Workers se estourar o orçamento de RAM
        est_ram = cls.estimate_ram_usage(dataset_size, imgsz, adj_batch, adj_workers, cache=adj_cache)
        while est_ram > max_allowed_ram and adj_workers > 2:
            adj_workers = max(2, adj_workers // 2)
            is_adjusted = True
            est_ram = cls.estimate_ram_usage(dataset_size, imgsz, adj_batch, adj_workers, cache=adj_cache)
            warnings.append(f"[RAM_GUARD] Workers rebaixados para {adj_workers} para respeitar teto de RAM.")

        # 3. Checagem e Redução de Batch Size se ainda exceder RAM ou VRAM
        est_vram = cls.estimate_vram_usage(imgsz, adj_batch, amp=amp)
        max_allowed_vram = available_vram * cls.SAFETY_MARGIN_VRAM if available_vram > 0 else float('inf')

        while (est_ram > max_allowed_ram or (available_vram > 0 and est_vram > max_allowed_vram)) and adj_batch > 4:
            adj_batch = max(4, adj_batch // 2)
            is_adjusted = True
            est_ram = cls.estimate_ram_usage(dataset_size, imgsz, adj_batch, adj_workers, cache=adj_cache)
            est_vram = cls.estimate_vram_usage(imgsz, adj_batch, amp=amp)
            warnings.append(f"[RAM_GUARD] Batch size rebaixado para {adj_batch} devido a limites de VRAM/RAM.")

        final_ram_mb = est_ram / (1024 * 1024)
        final_vram_mb = est_vram / (1024 * 1024)
        avail_ram_mb = available_ram / (1024 * 1024)
        avail_vram_mb = available_vram / (1024 * 1024)

        return MemoryPlan(
            batch_size=adj_batch,
            num_workers=adj_workers,
            cache=adj_cache,
            estimated_ram_mb=round(final_ram_mb, 2),
            estimated_vram_mb=round(final_vram_mb, 2),
            available_ram_mb=round(avail_ram_mb, 2),
            available_vram_mb=round(avail_vram_mb, 2),
            warnings=tuple(warnings),
            is_adjusted=is_adjusted
        )
