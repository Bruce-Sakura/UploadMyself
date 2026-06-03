"""3D 虚拟形象 API"""

from fastapi import APIRouter, UploadFile, File, Form
from typing import Optional

router = APIRouter()


@router.post("/upload")
async def upload_photo_3d(
    photos: list[UploadFile] = File(..., description="照片(1-4张，多角度更佳)"),
    name: str = Form(...),
):
    """上传照片用于 3D 重建"""
    return {"status": "uploaded", "message": f"{len(photos)} 张照片已上传"}


@router.post("/reconstruct")
async def reconstruct_3d(
    avatar_id: str = Form(...),
    quality: str = Form("medium", description="质量: low/medium/high"),
    format: str = Form("glb", description="输出格式: glb/vrm/obj"),
):
    """从照片重建 3D 模型"""
    # TODO: InstantMesh 3D 重建
    return {"status": "reconstructing", "avatar_id": avatar_id}


@router.get("/{avatar_id}/model")
async def download_3d_model(avatar_id: str):
    """下载 3D 模型文件"""
    return {"avatar_id": avatar_id}


@router.get("/{avatar_id}/preview")
async def preview_3d_model(avatar_id: str):
    """获取 3D 模型预览（Three.js 渲染数据）"""
    return {"avatar_id": avatar_id, "model_url": ""}
