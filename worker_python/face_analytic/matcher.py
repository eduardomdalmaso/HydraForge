"""
HydraForge Face Matcher and Cross-Camera ReID Registry.
Handles known gallery matching, unknown classification, and multi-camera association.
"""
import time
from dataclasses import dataclass, field
from typing import Dict, List, Optional, Tuple
import numpy as np


def cosine_similarity(v1: np.ndarray, v2: np.ndarray) -> float:
    """Computes cosine similarity between two normalized feature vectors."""
    n1 = np.linalg.norm(v1)
    n2 = np.linalg.norm(v2)
    if n1 == 0 or n2 == 0:
        return 0.0
    return float(np.dot(v1, v2) / (n1 * n2))


@dataclass
class IdentityProfile:
    subject_id: str
    name: str
    is_known: bool
    embedding: np.ndarray
    first_seen: float
    last_seen: float
    camera_trail: List[str] = field(default_factory=list)


class FaceEmbeddingMatcher:
    """
    Manages known biometric gallery, dynamic unknown visitors,
    and cross-camera re-identification (Cross-Camera ReID).
    """

    def __init__(
        self,
        unknown_score: float = 0.60,
        recognition_threshold: float = 0.75,
        max_trail_history: int = 50
    ):
        self.unknown_score = unknown_score
        self.recognition_threshold = recognition_threshold
        self.max_trail_history = max_trail_history
        self.known_gallery: Dict[str, IdentityProfile] = {}
        self.active_unknowns: Dict[str, IdentityProfile] = {}
        self.next_unknown_idx = 1

    def register_known_subject(self, subject_id: str, name: str, embedding: np.ndarray) -> None:
        """Registers an authorized/known identity in the gallery."""
        emb = np.array(embedding, dtype=np.float32)
        norm = np.linalg.norm(emb)
        if norm > 0:
            emb = emb / norm
        self.known_gallery[subject_id] = IdentityProfile(
            subject_id=subject_id,
            name=name,
            is_known=True,
            embedding=emb,
            first_seen=time.time(),
            last_seen=time.time(),
            camera_trail=[]
        )

    def identify_or_register(
        self, embedding: np.ndarray, camera_id: str
    ) -> Tuple[str, str, bool, float]:
        """
        Matches embedding against known gallery first, then cross-camera unknown registry.
        Returns: (subject_id, subject_name, is_known, match_confidence)
        """
        emb = np.array(embedding, dtype=np.float32)
        norm = np.linalg.norm(emb)
        if norm > 0:
            emb = emb / norm

        now = time.time()

        # 1. Match against Known Gallery
        best_known_id = None
        best_known_sim = -1.0
        for s_id, profile in self.known_gallery.items():
            sim = cosine_similarity(emb, profile.embedding)
            if sim > best_known_sim:
                best_known_sim = sim
                best_known_id = s_id

        if best_known_sim >= self.recognition_threshold and best_known_id:
            profile = self.known_gallery[best_known_id]
            profile.last_seen = now
            if not profile.camera_trail or profile.camera_trail[-1] != camera_id:
                profile.camera_trail.append(camera_id)
            return profile.subject_id, profile.name, True, round(best_known_sim, 3)

        # 2. Match against Cross-Camera Unknown Visitors Registry
        best_unk_id = None
        best_unk_sim = -1.0
        for u_id, profile in self.active_unknowns.items():
            sim = cosine_similarity(emb, profile.embedding)
            if sim > best_unk_sim:
                best_unk_sim = sim
                best_unk_id = u_id

        if best_unk_sim >= self.recognition_threshold and best_unk_id:
            # Re-identified unknown person moving across cameras!
            profile = self.active_unknowns[best_unk_id]
            profile.last_seen = now
            if not profile.camera_trail or profile.camera_trail[-1] != camera_id:
                profile.camera_trail.append(camera_id)
            return profile.subject_id, profile.name, False, round(best_unk_sim, 3)

        # 3. New Unknown Subject
        new_id = f"UNKNOWN_{self.next_unknown_idx:04d}"
        self.next_unknown_idx += 1
        name = f"Desconhecido #{self.next_unknown_idx - 1}"
        new_profile = IdentityProfile(
            subject_id=new_id,
            name=name,
            is_known=False,
            embedding=emb,
            first_seen=now,
            last_seen=now,
            camera_trail=[camera_id]
        )
        self.active_unknowns[new_id] = new_profile

        # Keep active unknowns bounded
        if len(self.active_unknowns) > 2000:
            oldest = min(self.active_unknowns.values(), key=lambda p: p.last_seen)
            del self.active_unknowns[oldest.subject_id]

        conf = max(0.0, best_known_sim if best_known_sim > 0 else 0.50)
        return new_id, name, False, round(conf, 3)
