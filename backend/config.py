"""配置管理"""

from pydantic_settings import BaseSettings
from functools import lru_cache


class Settings(BaseSettings):
    # App
    app_env: str = "development"
    app_host: str = "0.0.0.0"
    app_port: int = 8000
    secret_key: str = "change-me"

    # Database
    database_url: str = "postgresql+asyncpg://uploadmyself:uploadmyself@localhost:5432/uploadmyself"

    # Redis
    redis_url: str = "redis://localhost:6379/0"

    # MinIO
    minio_endpoint: str = "localhost:9000"
    minio_access_key: str = "minioadmin"
    minio_secret_key: str = "minioadmin"
    minio_bucket: str = "uploadmyself"

    # Provider
    provider_mode: str = "local"  # local | cloud | hybrid

    # Cloud API Keys
    openai_api_key: str = ""
    anthropic_api_key: str = ""
    qwen_api_key: str = ""
    elevenlabs_api_key: str = ""
    fish_audio_api_key: str = ""
    heygen_api_key: str = ""
    tripo3d_api_key: str = ""

    # Model paths
    model_dir: str = "./ml/models"

    # GPU
    cuda_visible_devices: str = "0"
    torch_device: str = "cuda"

    # Logging
    log_level: str = "INFO"

    class Config:
        env_file = ".env"


@lru_cache
def get_settings() -> Settings:
    return Settings()
