from app.settings import settings
from app.models import VehicleSummary


def write_to_file(summary: VehicleSummary):
    with open(settings.filename, "a") as file:
        file.write(
            f"{summary.vehicle_plate}, {summary.entry_date_time}, {summary.exit_date_time}, {summary.duration}\n"
        )
