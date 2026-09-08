"""Smoke test for PyTorch / CUDA / AMP training loop boundary."""
import pytest
import torch


@pytest.mark.smoke
@pytest.mark.slow
def test_cuda_training_loop_boundary_minimum_viable():
    """Validates PyTorch + CUDA + AMP FP16 execution on a synthetic 2-image micro-batch."""
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    use_amp = torch.cuda.is_available()

    # Synthetic micro-dataset: 2 fake images (3 channels, 32x32)
    dummy_input = torch.randn(2, 3, 32, 32, device=device)
    dummy_target = torch.tensor([0, 1], device=device, dtype=torch.long)

    model = torch.nn.Sequential(
        torch.nn.Conv2d(3, 8, kernel_size=3, padding=1),
        torch.nn.ReLU(),
        torch.nn.Flatten(),
        torch.nn.Linear(8 * 32 * 32, 2),
    ).to(device)

    optimizer = torch.optim.AdamW(model.parameters(), lr=1e-3)
    device_type = "cuda" if torch.cuda.is_available() else "cpu"
    scaler = torch.amp.GradScaler(device_type, enabled=use_amp)

    # 1 forward + backward step to verify gradient computation and lack of NaN
    with torch.amp.autocast(device_type, enabled=use_amp):
        output = model(dummy_input)
        loss = torch.nn.functional.cross_entropy(output, dummy_target)

    scaler.scale(loss).backward()
    scaler.step(optimizer)
    scaler.update()

    assert not torch.isnan(loss), "Loss must be a valid real number"
    assert model[0].weight.grad is not None, "Convolutional layer gradients must be computed"
    assert loss.item() > 0.0, "Cross entropy loss must be positive"
