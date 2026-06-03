"""模型注册与管理"""

from dataclasses import dataclass
from enum import Enum


class ModelType(Enum):
    LLM = "llm"
    VOICE = "voice"
    AVATAR_2D = "avatar_2d"
    AVATAR_3D = "avatar_3d"


@dataclass
class ModelInfo:
    name: str
    model_type: ModelType
    path: str
    provider: str  # local | cloud
    loaded: bool = False


class ModelRegistry:
    """管理所有可用模型"""

    def __init__(self):
        self._models: dict[str, ModelInfo] = {}

    def register(self, model: ModelInfo):
        self._models[model.name] = model

    def get(self, name: str) -> ModelInfo | None:
        return self._models.get(name)

    def list_by_type(self, model_type: ModelType) -> list[ModelInfo]:
        return [m for m in self._models.values() if m.model_type == model_type]


# 全局注册表
registry = ModelRegistry()
