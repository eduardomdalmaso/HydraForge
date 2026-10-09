"""
HydraForge Face Quality Assessment Module.
Implements blur detection via Laplacian variance, aspect ratio, and composite scoring.
"""
import math
from typing import Tuple, Optional
import numpy as np

try:
    import cv2
    HAS_CV2 = True
except ImportError:
    HAS_CV2 = False


def compute_laplacian_variance(crop: np.ndarray) -> float:
    """
    Computes sharpness score using Laplacian Variance: Var(nabla^2 I).
    Higher values represent sharper edges; lower values indicate blur.
    """
    if crop is None or crop.size == 0:
        return 0.0

    # Ensure grayscale
    if len(crop.shape) == 3:
        if HAS_CV2:
            gray = cv2.cvtColor(crop, cv2.COLOR_BGR2GRAY)
        else:
            # Fallback luminance conversion: Y = 0.299R + 0.587G + 0.114B
            gray = (0.114 * crop[:, :, 0] + 0.587 * crop[:, :, 1] + 0.299 * crop[:, :, 2]).astype(np.uint8)
    else:
        gray = crop

    if HAS_CV2:
        lap = cv2.Laplacian(gray, cv2.CV_64F)
        return float(lap.var())

    # Pure NumPy fallback 3x3 Laplacian convolution
    h, w = gray.shape
    if h < 3 or w < 3:
        return 0.0
    kernel = np.array([[0, 1, 0], [1, -4, 1], [0, 1, 0]], dtype=np.float64)
    sub = gray.astype(np.float64)
    # Simple valid convolution via slicing
    conv = (
        sub[:-2, 1:-1] * kernel[0, 1] +
        sub[1:-1, :-2] * kernel[1, 0] +
        sub[1:-1, 1:-1] * kernel[1, 1] +
        sub[1:-1, 2:] * kernel[1, 2] +
        sub[2:, 1:-1] * kernel[2, 1]
    )
    return float(conv.var())


def compute_aspect_ratio_score(width: int, height: int) -> float:
    """
    Evaluates whether face width/height matches human biometric norms (~0.75 - 0.90).
    Penalizes extreme profile or corrupted crops.
    """
    if height <= 0 or width <= 0:
        return 0.0
    ratio = width / float(height)
    # Optimum face ratio is around 0.80
    diff = abs(ratio - 0.80)
    score = max(0.0, 1.0 - (diff / 0.50))
    return round(score, 3)


def compute_composite_quality(
    crop: np.ndarray,
    confidence: float,
    blur_threshold: float = 60.0
) -> Tuple[float, float, bool]:
    """
    Computes composite quality score in [0.0, 1.0] and returns:
    (composite_score, blur_score, passes_blur_threshold)
    """
    if crop is None or crop.size == 0:
        return 0.0, 0.0, False

    h, w = crop.shape[:2]
    blur_score = compute_laplacian_variance(crop)
    passes_blur = blur_score >= blur_threshold

    # Normalized blur score via sigmoid centered at threshold
    norm_blur = 1.0 / (1.0 + math.exp(-max(-10.0, min(10.0, (blur_score - blur_threshold) / 25.0))))

    # Aspect ratio score
    aspect_score = compute_aspect_ratio_score(w, h)

    # Resolution score (scales smoothly from 40px to 200px)
    res_score = min(1.0, max(0.2, h / 200.0))

    # Composite weighted score:
    # 40% blur/sharpness, 30% confidence, 15% aspect ratio, 15% resolution
    composite = (
        0.40 * norm_blur +
        0.30 * confidence +
        0.15 * aspect_score +
        0.15 * res_score
    )

    return round(composite, 4), round(blur_score, 2), passes_blur
