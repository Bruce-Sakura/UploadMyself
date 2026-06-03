"""模型提供者抽象"""

from abc import ABC, abstractmethod
from typing import Any


class BaseProvider(ABC):
    """统一模型提供者接口"""

    @abstractmethod
    async def inference(self, *args, **kwargs) -> Any:
        ...


class LocalProvider(BaseProvider):
    """本地 GPU 推理"""

    async def inference(self, *args, **kwargs):
        # TODO: 本地模型推理
        raise NotImplementedError


class CloudProvider(BaseProvider):
    """云端 API 调用"""

    async def inference(self, *args, **kwargs):
        # TODO: 云端 API 调用
        raise NotImplementedError


def get_provider(mode: str = "local") -> BaseProvider:
    """根据配置返回对应 Provider"""
    if mode == "cloud":
        return CloudProvider()
    return LocalProvider()
