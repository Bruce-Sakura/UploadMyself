# 🧬 UploadMyself

> **克隆你自己** — 输入照片 + 文本语料 + 语音样本，生成你的数字分身

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go 1.22+](https://img.shields.io/badge/go-1.22+-00ADD8.svg)](https://go.dev/)
[![React 18](https://img.shields.io/badge/react-18-61DAFB.svg)](https://react.dev/)

---

## ✨ 功能特性

| 功能 | 说明 | 模型方案 |
|------|------|---------|
| 🧠 **思维框架克隆** | 输入文本语料，生成你的思维模式 Skill | 仿女娲 Skill + LLM |
| 🎤 **语音克隆** | 3-10 分钟音频，克隆你的声音 | GPT-SoVITS / CosyVoice2 |
| 🖼️ **2D 虚拟形象** | 一张照片生成动态说话视频 | LivePortrait + SadTalker |
| 🧊 **3D 虚拟形象** | 一张照片生成可交互 3D 模型 | InstantMesh + Three.js |
| 🔬 **模型蒸馏** | 大模型压缩为轻量版 | 知识蒸馏 pipeline |

### 🔀 双模式运行

- **本地模式**：GPU 推理，数据不出本机
- **云端模式**：调用第三方 API，免部署
- **混合模式**：核心本地 + 增强云端

---

## 🏗️ 技术栈

**后端**：Golang (Gin + GORM + Asynq) + Redis + PostgreSQL
**前端**：React 18 + TypeScript + Three.js + Ant Design
**ML**：GPT-SoVITS / CosyVoice2 / LivePortrait / InstantMesh

详见 [PROJECT_PLAN.md](PROJECT_PLAN.md)

---

## 🚀 快速开始

### 环境要求

- Go 1.22+
- Node.js 18+
- CUDA 11.8+ (本地推理需要 NVIDIA GPU)
- Docker & Docker Compose (推荐)

### 安装

```bash
# 克隆仓库
git clone https://github.com/Bruce-Sakura/UploadMyself.git
cd UploadMyself

# 后端依赖
cd backend && go mod tidy && cd ..

# 前端依赖
cd frontend && npm install && cd ..

# 下载模型
make models-download
```

### 启动

```bash
# 启动依赖服务
docker-compose up -d redis postgres minio

# 后端
cd backend && go run .

# 前端
cd frontend && npm run dev
```

访问 http://localhost:5173

---

## 📁 项目结构

```
UploadMyself/
├── backend/          # Golang 后端 (Gin)
│   ├── api/          # 路由 & Handler
│   ├── core/         # 核心引擎
│   │   ├── skill_engine/    # 思维框架克隆
│   │   ├── voice_engine/    # 语音克隆
│   │   ├── avatar_engine/   # 虚拟形象
│   │   └── distill_engine/  # 模型蒸馏
│   ├── models/       # 数据模型
│   ├── services/     # 服务层
│   │   └── provider/ # Provider (本地/云端)
│   ├── workers/      # Asynq 异步任务
│   └── config/       # 配置管理
├── ml/               # ML 模型与脚本 (Python)
├── frontend/         # React + Three.js 前端
├── skills/           # 生成的 Skill 存放
├── docs/             # 文档
└── tests/            # 测试
```

---

## 📖 文档

- [架构设计](docs/architecture.md)
- [API 文档](docs/api.md) — 启动后访问 http://localhost:8000/docs
- [部署指南](docs/deployment.md)
- [模型说明](docs/models.md)

---

## 🤝 贡献

欢迎贡献！请阅读 [contributing.md](docs/contributing.md)

---

## 📄 许可证

MIT License — 详见 [LICENSE](LICENSE)

---

> *"The best way to predict the future is to create it."* — Upload yourself into the digital world.
