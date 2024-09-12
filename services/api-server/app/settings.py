from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file='.env', env_file_encoding='utf-8')

    port: int = 8000
    log_level: str = "DEBUG"
    filename: str

settings = Settings()



