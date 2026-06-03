"""FastAPI 入口"""

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from backend.api import skill, voice, avatar_2d, avatar_3d, distill, tasks

app = FastAPI(
    title="UploadMyself",
    description="克隆你自己 — 数字分身生成平台",
    version="0.1.0",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:5173", "http://localhost:3000"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 注册路由
app.include_router(skill.router, prefix="/api/v1/skill", tags=["思维框架"])
app.include_router(voice.router, prefix="/api/v1/voice", tags=["语音克隆"])
app.include_router(avatar_2d.router, prefix="/api/v1/avatar/2d", tags=["2D形象"])
app.include_router(avatar_3d.router, prefix="/api/v1/avatar/3d", tags=["3D形象"])
app.include_router(distill.router, prefix="/api/v1/distill", tags=["模型蒸馏"])
app.include_router(tasks.router, prefix="/api/v1/tasks", tags=["任务管理"])


@app.get("/")
async def root():
    return {
        "name": "UploadMyself",
        "version": "0.1.0",
        "docs": "/docs",
    }


@app.get("/health")
async def health():
    return {"status": "ok"}
