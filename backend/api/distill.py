"""模型蒸馏 API"""

from fastapi import APIRouter, Form
from typing import Optional

router = APIRouter()


@router.post("/start")
async def start_distillation(
    teacher_model: str = Form(..., description="教师模型"),
    student_model: str = Form(..., description="学生模型架构"),
    task_type: str = Form(..., description="任务类型: llm/voice/avatar_2d"),
    dataset: Optional[str] = Form(None, description="蒸馏数据集路径"),
    epochs: int = Form(10),
    temperature: float = Form(2.0, description="蒸馏温度"),
    alpha: float = Form(0.5, description="Loss 权重 (KD vs CE)"),
):
    """启动模型蒸馏任务"""
    # TODO: Celery 异步蒸馏
    return {"status": "distillation_started", "task_id": "placeholder"}


@router.get("/{task_id}/status")
async def distillation_status(task_id: str):
    """查询蒸馏进度"""
    return {
        "task_id": task_id,
        "status": "training",
        "epoch": 0,
        "total_epochs": 10,
        "loss": 0.0,
    }


@router.get("/{task_id}/metrics")
async def distillation_metrics(task_id: str):
    """获取蒸馏指标对比"""
    return {
        "task_id": task_id,
        "teacher": {"accuracy": 0.0, "latency_ms": 0, "size_mb": 0},
        "student": {"accuracy": 0.0, "latency_ms": 0, "size_mb": 0},
    }
