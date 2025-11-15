from pydantic_settings import BaseSettings
from typing import Optional

class Settings(BaseSettings):
    DATABASE_URL: str
    JWT_SECRET: str
    JWT_ALGORITHM: str = "HS256"
    JWT_EXPIRATION_HOURS: int = 1
    OPENAI_API_KEY: str
    ENV: str = "development"
    DEBUG: bool = True

    class Config:
        env_file = ".env"

settings = Settings()