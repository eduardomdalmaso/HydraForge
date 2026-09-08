"""Zero-copy Arrow IPC streaming directly over POSIX /dev/shm for inter-process communication."""
from __future__ import annotations
from pathlib import Path
import polars as pl


def sink_to_shm_ipc(
    df: pl.DataFrame,
    shm_path: str | Path = "/dev/shm/hydra_curation.arrow",
) -> Path:
    """Sinks a Polars DataFrame into a fast POSIX shared memory Arrow IPC stream."""
    path = Path(shm_path)
    # Garante que o diretório existe se for arquivo temporário local
    path.parent.mkdir(parents=True, exist_ok=True)

    with open(path, "wb") as f:
        df.write_ipc_stream(f)

    return path


def read_from_shm_ipc(
    shm_path: str | Path = "/dev/shm/hydra_curation.arrow",
) -> pl.DataFrame:
    """Reads a Polars DataFrame instantly from a POSIX shared memory Arrow IPC stream."""
    path = Path(shm_path)
    if not path.exists():
        raise FileNotFoundError(f"Shared memory IPC stream not found at: {path}")

    with open(path, "rb") as f:
        return pl.read_ipc_stream(f)
