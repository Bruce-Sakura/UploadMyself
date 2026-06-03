"""语音克隆 API"""

from fastapi import APIRouter, UploadFile, File, Form, BackgroundTasks
from typing import Optional

router = APIRouter()


@router.post("/upload")
async def upload_voice_sample(
    audio: UploadFile = File(..., description="语音样本文件(.wav/.mp3/.flac)"),
    name: str = Form(..., description="声音名称"),
):
    """上传语音样本"""
    # TODO: 存储音频，启动预处理
    return {"status": "uploaded", "message": "语音样本已上传"}


@router.post("/train")
async def train_voice(
    background_tasks: BackgroundTasks,
    voice_id: str = Form(...),
    text: Optional[str] = Form(None, description="对应文本（可选，用于强制对齐）"),
    epochs: int = Form(5, description="训练轮数"),
    provider: Optional[str] = Form(None),
):
    """启动声音模型训练"""
    # TODO: Celery 异步训练任务
    return {"status": "training_started", "voice_id": voice_id}


@router.post("/synthesize")
async def synthesize_speech(
    voice_id: str = Form(...),
    text: str = Form(..., description="要合成的文本"),
    speed: float = Form(1.0, description="语速"),
    emotion: str = Form("neutral", description="情感"),
):
    """用克隆的声音合成语音"""
    # TODO: 推理合成，返回音频文件
    return {"status": "synthesized", "voice_id": voice_id}


@router.get("/{voice_id}/samples")
async def get_voice_samples(voice_id: str):
    """获取声音样本和试听"""
    return {"voice_id": voice_id, "samples": []}
