"""Unit tests for Garbage Collection tuning utility."""
import gc
import pytest
from training.gc_tuning import disable_gc_context, tune_gc_thresholds
from curation.domain import BoundingBox


@pytest.mark.unit
def test_disable_gc_context():
    gc.enable()
    assert gc.isenabled() is True

    with disable_gc_context():
        assert gc.isenabled() is False

    assert gc.isenabled() is True


@pytest.mark.unit
def test_bounding_box_has_slots():
    bbox = BoundingBox(class_id=0, x_center=0.5, y_center=0.5, w=0.2, h=0.2)
    # Instances with __slots__ do not have a __dict__ attribute
    assert not hasattr(bbox, "__dict__")
    assert hasattr(bbox, "__slots__")


@pytest.mark.unit
def test_tune_gc_thresholds():
    orig_threshold = gc.get_threshold()
    tune_gc_thresholds(alloc_threshold=20000)
    new_threshold = gc.get_threshold()
    assert new_threshold[0] == 20000
    # Restore
    gc.set_threshold(*orig_threshold)
