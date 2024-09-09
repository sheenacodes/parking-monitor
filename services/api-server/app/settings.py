from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    port: int = 8000
    filename: str = "/project/log/log.txt"
    log_level: str = "DEBUG"

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"


settings = Settings()
