"""异步任务状态查询"""

from fastapi import APIRouter

router = APIRouter()


@router.get("/{task_id}")
async def get_task_status(task_id: str):
    """查询异步任务状态"""
    # TODO: 从 Redis 查询 Celery 任务状态
    return {
        "task_id": task_id,
        "status": "pending",  # pending | running | completed | failed
        "progress": 0,
        "result": None,
        "error": None,
    }
