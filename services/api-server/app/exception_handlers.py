from fastapi import Request
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse
from app.models import ErrorResponse

import logging

logger = logging.getLogger(__name__)


async def custom_validation_exception_handler(
    request: Request, exc: RequestValidationError
):
    error_response = ErrorResponse(detail="Invalid Request data")
    return JSONResponse(status_code=422, content=error_response.dict())
