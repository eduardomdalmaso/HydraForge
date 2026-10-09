"""
HydraForge Face Analytics Package.
"""
from .config import FaceAnalyticConfig
from .quality import compute_laplacian_variance, compute_composite_quality
from .best_frame import BestFrameBuffer, CandidateFace
from .tracker import PersonGatedTracker, TrackedPerson
from .matcher import FaceEmbeddingMatcher, IdentityProfile
from .pipeline import FaceAnalyticPipeline

__all__ = [
    "FaceAnalyticConfig",
    "compute_laplacian_variance",
    "compute_composite_quality",
    "BestFrameBuffer",
    "CandidateFace",
    "PersonGatedTracker",
    "TrackedPerson",
    "FaceEmbeddingMatcher",
    "IdentityProfile",
    "FaceAnalyticPipeline",
]
