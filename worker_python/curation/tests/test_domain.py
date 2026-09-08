"""Pure domain and fake repository tests."""
import pytest
from curation.domain import BoundingBox
from curation.fakes import FakeJobRepository


@pytest.mark.unit
def test_repository_saves_and_retrieves_bboxes_without_database():
    # Arrange (Given)
    fake_repo = FakeJobRepository()
    bbox = BoundingBox(class_id=0, x_center=0.5, y_center=0.5, w=0.2, h=0.2)

    # Act (When)
    fake_repo.save("job-123", [bbox])

    # Assert (Then)
    saved_boxes = fake_repo.get_bboxes("job-123")
    assert len(saved_boxes) == 1
    assert saved_boxes[0].is_valid_yolo() is True
    assert saved_boxes[0].class_id == 0


@pytest.mark.unit
def test_repository_delete_and_empty_retrieval():
    fake_repo = FakeJobRepository()
    assert fake_repo.get_bboxes("non-existent") == []
    
    bbox = BoundingBox(class_id=1, x_center=0.4, y_center=0.4, w=0.1, h=0.1)
    fake_repo.save("job-abc", [bbox])
    assert fake_repo.delete("job-abc") is True
    assert fake_repo.get_bboxes("job-abc") == []


@pytest.mark.unit
def test_iou_calculation_exact_match_and_disjoint():
    b1 = BoundingBox(class_id=0, x_center=0.5, y_center=0.5, w=0.2, h=0.2)
    b2 = BoundingBox(class_id=0, x_center=0.5, y_center=0.5, w=0.2, h=0.2)
    b3 = BoundingBox(class_id=0, x_center=0.9, y_center=0.9, w=0.1, h=0.1)

    # Exact match must yield IoU = 1.0
    assert pytest.approx(b1.iou(b2), 1e-5) == 1.0
    # Disjoint boxes must yield IoU = 0.0
    assert b1.iou(b3) == 0.0
