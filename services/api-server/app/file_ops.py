from app.settings import settings
from app.models import VehicleSummary
import logging

logger = logging.getLogger(__name__)


def write_to_file(summary: VehicleSummary):
    logger.debug(f"writing log to file {settings.filename}")
    with open(settings.filename, "a") as file:
        file.write(
            f"{summary.vehicle_plate}, {summary.entry_date_time}, {summary.exit_date_time}, {summary.duration}\n"
        )
        file.flush()
