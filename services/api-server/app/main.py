from fastapi import FastAPI
from fastapi.exceptions import RequestValidationError
from app.routes import router as parking_log_router
from app.settings import settings
from app.exception_handlers import custom_validation_exception_handler


app = FastAPI(
    title="Vehicle Event Recorder",
    description="A REST API to record vehicle plate, time of entry and exit from parking lot, and duration of parking.",
    version="1.0.0",
)
app.include_router(parking_log_router)


@app.get("/")
def read_root():
    return {"message": "Welcome to the Vehicle Parking Summary API"}


app.add_exception_handler(RequestValidationError, custom_validation_exception_handler)


def main():
    import uvicorn
    print(settings.port)
    uvicorn.run(app, host="0.0.0.0", port=settings.port)


if __name__ == "__main__":
    main()