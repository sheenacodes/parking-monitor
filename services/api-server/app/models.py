from pydantic import BaseModel, field_validator
from datetime import datetime
import re


def validate_date_format(date: str):
    # Regular expression to match ISO 8601 format YYYY-MM-DDTHH:MM:SSZ
    # Check only UTC that ending with Z ( no time offsets)
    iso8601_regex = r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z?$"
    if not re.match(iso8601_regex, date):
        raise ValueError(
            "Invalid datetime format. Please use ISO 8601 format YYYY-MM-DDTHH:MM:SSZ "
        )

    # try parse it as a datetime to check it's valid datetime
    try:
        datetime.fromisoformat(date.replace("Z", "+00:00"))
    except ValueError:
        raise ValueError(
            "Invalid datetime. Please use ISO 8601 format YYYY-MM-DDTHH:MM:SSZ  "
        )


# TODO validate duration str


class VehicleSummary(BaseModel):
    vehicle_plate: str
    entry_date_time: str
    exit_date_time: str
    duration: str

    @field_validator("entry_date_time", "entry_date_time")
    @classmethod
    def validate_iso8601(cls, v):
        validate_date_format(v)
        return v


class SuccessResponse(BaseModel):
    detail: str


class ErrorResponse(BaseModel):
    detail: str
