# 🧬 UploadMyself

> **Clone Yourself** — Upload your photo, text corpus, and voice sample to generate your digital twin.

[中文版](#中文说明)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go 1.22+](https://img.shields.io/badge/go-1.22+-00ADD8.svg)](https://go.dev/)
[![React 18](https://img.shields.io/badge/react-18-61DAFB.svg)](https://react.dev/)

---

## ✨ Features

| Feature | Description | Model |
|---------|-------------|-------|
| 🧠 **Skill Cloning** | Generate your thinking framework from text corpus | Nuwa-style LLM analysis |
| 🎤 **Voice Cloning** | Clone your voice from 1-10 min audio | GPT-SoVITS / CosyVoice2 / Fish Speech |
| 🖼️ **2D Avatar** | One photo → animated talking video | LivePortrait / HunyuanVideo-Avatar |
| 🧊 **3D Avatar** | One photo → interactive 3D model | InstantMesh / TripoSR + Three.js |
| 🔬 **Model Distillation** | Compress large models for edge deployment | Knowledge Distillation pipeline |

### Dual Mode

- **Local Mode** — GPU inference, data never leaves your machine
- **Cloud Mode** — Third-party APIs, zero deployment
- **Hybrid Mode** — Core local + enhancement via cloud

---

## 🏗️ Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go (Gin + GORM + Asynq + Viper + Zap) |
| **Frontend** | React 18 + TypeScript + Three.js + Ant Design |
| **Database** | PostgreSQL + Redis |
| **Storage** | MinIO (S3-compatible) |
| **ML** | PyTorch (called from Go via subprocess/gRPC) |

---

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- Node.js 18+
- Docker & Docker Compose
- NVIDIA GPU (for local ML inference)

### Install

```bash
git clone https://github.com/Bruce-Sakura/UploadMyself.git
cd UploadMyself

# Backend
cd backend && go mod tidy && cd ..

# Frontend
cd frontend && npm install && cd ..

# Download ML models
make models-download
```

### Run

```bash
# Start dependencies
docker-compose up -d redis postgres minio

# Backend (port 8000)
cd backend && go run .

# Frontend (port 5173)
cd frontend && npm run dev
```

Open http://localhost:5173

### Build

```bash
make build          # Build Go binary
make frontend-build # Build React app
make test           # Run all tests
make lint           # Run linters
```

---

## 📁 Project Structure

```
UploadMyself/
├── backend/                 # Go backend (Gin)
│   ├── api/                 # Routes & Handlers
│   ├── config/              # Viper config
│   ├── core/                # Core engines
│   │   ├── skill_engine/    # Thinking framework cloning
│   │   ├── voice_engine/    # Voice cloning
│   │   ├── avatar_engine/   # 2D/3D avatar
│   │   └── distill_engine/  # Model distillation
│   ├── models/              # Data models
│   ├── services/            # Service layer
│   │   └── provider/        # Local / Cloud provider
│   └── workers/             # Asynq async tasks
├── ml/                      # ML models & Python scripts
│   ├── models/              # Pretrained weights
│   └── scripts/             # Preprocessing & training
├── frontend/                # React + Three.js
│   └── src/
│       ├── pages/           # UI pages
│       └── three/           # 3D rendering
├── skills/                  # Generated user skills
├── docs/                    # Documentation
└── tests/                   # Tests
```

---

## 📖 Docs

- [Architecture](docs/architecture.md)
- [API Reference](docs/api.md) — http://localhost:8000/docs after startup
- [Deployment Guide](docs/deployment.md)
- [Model Details](docs/models.md)
- [Development TODO](TODO.md)

---

## 🤝 Contributing

Contributions welcome! See [contributing.md](docs/contributing.md).

---

## 📄 License

MIT License — see [LICENSE](LICENSE).

---

> *"The best way to predict the future is to create it."* — Upload yourself into the digital world.

---

# 中文说明

## 🧬 UploadMyself — 克隆你自己

> 上传你的照片、文本语料、语音样本，生成你的数字分身。

### 功能一览

| 功能 | 说明 | 模型方案 |
|------|------|---------|
| 🧠 **思维框架克隆** | 输入文本语料，生成你的思维模式 Skill | 仿女娲 + LLM |
| 🎤 **语音克隆** | 1-10 分钟音频克隆你的声音 | GPT-SoVITS v2 / CosyVoice2 / Fish Speech v1.5 |
| 🖼️ **2D 虚拟形象** | 一张照片生成动态说话视频 | LivePortrait / HunyuanVideo-Avatar (2025最新) |
| 🧊 **3D 虚拟形象** | 一张照片生成可交互 3D 模型 | InstantMesh / TripoSR + Three.js |
| 🔬 **模型蒸馏** | 大模型压缩，降低部署成本 | NVIDIA TensorRT + 知识蒸馏 |

### 技术栈

- **后端**：Golang (Gin + GORM + Asynq + Viper + Zap)
- **前端**：React 18 + TypeScript + Three.js + Ant Design
- **数据库**：PostgreSQL + Redis
- **存储**：MinIO
- **ML 推理**：PyTorch (通过 Go subprocess/gRPC 调用)

### 快速开始

```bash
git clone https://github.com/Bruce-Sakura/UploadMyself.git
cd UploadMyself

# 后端
cd backend && go mod tidy && cd ..

# 前端
cd frontend && npm install && cd ..

# 启动依赖
docker-compose up -d redis postgres minio

# 启动后端 (端口 8000)
cd backend && go run .

# 启动前端 (端口 5173)
cd frontend && npm run dev
```

### 开发进度

详见 [TODO.md](TODO.md) — 包含 7 个阶段、80+ 个子任务的技术路线图。

### 文档

- [架构设计](docs/architecture.md)
- [API 文档](docs/api.md)
- [模型说明](docs/models.md)
