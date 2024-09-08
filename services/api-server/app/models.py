from pydantic import BaseModel


class VehicleSummary(BaseModel):
    vehicle_plate: str
    entry_date_time: str
    exit_date_time: str
    duration: str


class ErrorResponse(BaseModel):
    message: str
