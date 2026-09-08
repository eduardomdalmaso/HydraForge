"""In-memory Fake Repositories for ultra-fast, mock-free unit testing."""
from __future__ import annotations
from typing import Dict, List, Optional
from .domain import BoundingBox, JobRepository


class FakeJobRepository:
    """In-memory fake implementation of JobRepository."""

    def __init__(self) -> None:
        self._store: Dict[str, List[BoundingBox]] = {}

    def save(self, job_id: str, bboxes: List[BoundingBox]) -> None:
        """Saves bboxes associated with job_id."""
        self._store[job_id] = list(bboxes)

    def get_bboxes(self, job_id: str) -> List[BoundingBox]:
        """Retrieves bboxes for job_id, or returns empty list if not found."""
        return list(self._store.get(job_id, []))

    def delete(self, job_id: str) -> bool:
        """Deletes job_id from store."""
        if job_id in self._store:
            del self._store[job_id]
            return True
        return False

    def clear(self) -> None:
        """Clears all in-memory entries."""
        self._store.clear()
