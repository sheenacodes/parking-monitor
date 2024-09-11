import pytest
from fastapi.testclient import TestClient
from fastapi import HTTPException
from app.routes import router
from app.models import VehicleSummary, SuccessResponse, ErrorResponse
from app.file_ops import write_to_file
import logging

# Create a FastAPI app for testing
from fastapi import FastAPI
app = FastAPI()
app.include_router(router)

@pytest.fixture
def client():
    return TestClient(app)

@pytest.fixture
def mock_write_to_file(mocker):
    # Mock the write_to_file function
    return mocker.patch("app.file_ops.write_to_file")

def test_log_vehicle_exit_success(client, mock_write_to_file):
    # Arrange
    summary = VehicleSummary(
        vehicle_plate="ABC123",
        entry_date_time="2024-09-11T21:24:56.833597372Z",
        exit_date_time="2024-09-11T22:24:56.833597372Z",
        duration="3600"
    )
    mock_write_to_file.return_value = None  # Simulate successful file write

    # Act
    response = client.post("/parkinglog", json=summary.dict())

    # Assert
    assert response.status_code == 201
    assert response.json() == {"detail": "Vehicle summary recorded successfully"}
    mock_write_to_file.assert_called_once_with(summary)

def test_log_vehicle_exit_failure(client, mock_write_to_file):
    # Arrange
    summary = VehicleSummary(
        vehicle_plate="ABC123",
        entry_date_time="2024-09-11T21:24:56.833597372Z",
        exit_date_time="2024-09-11T22:24:56.833597372Z",
        duration="3600"
    )
    mock_write_to_file.side_effect = Exception("File write error")  # Simulate a failure

    # Act
    response = client.post("/parkinglog", json=summary.dict())

    # Assert
    assert response.status_code == 500
    assert response.json() == {"detail": "Failed to record event to file"}
    mock_write_to_file.assert_called_once_with(summary)
