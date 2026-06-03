"""思维框架 Skill 生成 API"""

from fastapi import APIRouter, UploadFile, File, Form, BackgroundTasks
from typing import Optional

router = APIRouter()


@router.post("/create")
async def create_skill(
    background_tasks: BackgroundTasks,
    name: str = Form(..., description="Skill 名称"),
    corpus: UploadFile = File(..., description="文本语料文件(.txt/.md)"),
    style: str = Form("auto", description="风格：auto/analytical/creative/pragmatic"),
    provider: Optional[str] = Form(None, description="Provider: local/cloud"),
):
    """上传文本语料，生成思维框架 Skill"""
    # TODO: 实现女娲式 Skill 生成流程
    return {
        "status": "accepted",
        "message": "Skill 生成任务已提交",
        "task_id": "placeholder",
    }


@router.get("/{skill_id}/result")
async def get_skill_result(skill_id: str):
    """获取 Skill 生成结果"""
    # TODO: 从数据库/S3 获取生成的 SKILL.md
    return {"skill_id": skill_id, "status": "pending"}


@router.get("/{skill_id}/download")
async def download_skill(skill_id: str):
    """下载 Skill 文件包"""
    # TODO: 打包为 zip 返回
    return {"skill_id": skill_id}
