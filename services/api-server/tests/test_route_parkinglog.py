import pytest
from unittest.mock import patch
from app.models import VehicleSummary
from app.file_ops import write_to_file
from fastapi.testclient import TestClient
from app.main import app 

client = TestClient(app)

@pytest.fixture
def mock_write_to_file():
    with patch("app.routes.write_to_file") as mock:
        yield mock


def test_log_vehicle_exit_success(mock_write_to_file):
 
    summary = VehicleSummary(
        vehicle_plate="ABC123",
        entry_date_time="2024-09-11T21:24:56.833597372Z",
        exit_date_time="2024-09-11T22:24:56.833597372Z",
        duration="3600",
    )

    # Mock successful file write
    mock_write_to_file.return_value = None  


    response = client.post("/parkinglog", json=summary.model_dump())


    assert response.status_code == 201
    assert response.json() == {"detail": "Vehicle summary recorded successfully"}
    mock_write_to_file.assert_called_once_with(summary)


def test_log_vehicle_exit_failure(mock_write_to_file):
   
    summary = VehicleSummary(
        vehicle_plate="ABC123",
        entry_date_time="2024-09-11T21:24:56.833597372Z",
        exit_date_time="2024-09-11T22:24:56.833597372Z",
        duration="3600",
    )
    # Simulate a failure on file write
    mock_write_to_file.side_effect = Exception("File write error")  

  
    response = client.post("/parkinglog", json=summary.model_dump())

 
    assert response.status_code == 500
    assert response.json() == {"detail": "Failed to record vehicle summary"}
