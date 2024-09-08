from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    port: int = 8000
    filename: str = "vehicle_summary.txt"

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"


settings = Settings()
