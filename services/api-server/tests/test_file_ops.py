import pytest
from unittest.mock import MagicMock
from app.models import VehicleSummary
from app.file_ops import write_to_file

@pytest.fixture
def mock_settings(mocker):
    # Mock the settings object
    mock_settings = mocker.patch("app.main.settings")
    mock_settings.filename = "test_log.txt"
    return mock_settings

@pytest.fixture
def mock_open(mocker):
    # Mock the open function
    mock_open = mocker.patch("builtins.open", mocker.mock_open())
    return mock_open


def test_write_to_file(mock_settings, mock_open):

    summary = VehicleSummary(
        vehicle_plate="ABC123",
        entry_date_time="2024-09-11T21:24:56.833597372Z",
        exit_date_time="2024-09-11T22:24:56.833597372Z",
        duration=3600
    )


    write_to_file(summary)


    # Verify if the open function was called correctly
    mock_open.assert_called_once_with("test_log.txt", "a")

    # Verify if the file write was called with the correct content
    mock_open().write.assert_called_once_with(
        "ABC123, 2024-09-11T21:24:56.833597372Z, 2024-09-11T22:24:56.833597372Z, 3600\n"
    )
