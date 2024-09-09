from fastapi import APIRouter, HTTPException, status
from app.models import VehicleSummary
from app.file_ops import write_to_file
from typing import Dict
import logging

logger = logging.getLogger(__name__)

router = APIRouter()


@router.post(
    "/parkinglog",
    status_code=status.HTTP_201_CREATED,
    response_model=Dict[str, str],
    summary="Record vehicle parking log",
    description="Records a parking log containing entry, exit, and parking duration to a local file.",
)
async def log_vehicle_exit(summary: VehicleSummary):
    logger.debug(f"post /parkinglog called")
    try:
        write_to_file(summary)
        return {"message": "Vehicle summary recorded successfully"}
    except Exception as e:
        logger.error(str(e))
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to record event to file",
        )
