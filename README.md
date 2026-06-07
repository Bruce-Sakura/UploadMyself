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
| **Backend** | Go (Gin + pgx/pgxpool + Viper + Zap) |
| **Frontend** | React 18 + TypeScript + Three.js + Ant Design |
| **Database** | PostgreSQL + Redis |
| **Storage** | MinIO (S3-compatible) |
| **ML** | PyTorch (called from Go via subprocess/gRPC) |

---

## 🚀 Quick Start

### ⚡ One-command install / 一键安装

```bash
git clone https://github.com/Bruce-Sakura/UploadMyself.git
cd UploadMyself
bash install.sh            # 交互选择 Docker / Linux 原生
```

See **[INSTALL.md](INSTALL.md)** — supports **Docker** and **Linux native**, full features incl. GPU ML, runtime data under `~/.UploadMyself`.
详见 **[INSTALL.md](INSTALL.md)**：支持 **Docker** 与 **Linux 原生** 两种方式，全功能含 GPU ML，运行时数据集中在 `~/.UploadMyself`。

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

# Build the upme CLI
cd backend && go build -o upme ./cmd/upme
```

---

## 🖥️ CLI (`upme`)

`upme` is a small CLI client for the backend REST API.

```bash
cd backend && go build -o upme ./cmd/upme
```

The server address defaults to `http://localhost:8000`; override it with the
`-server` flag or the `$UPME_SERVER` environment variable.

```bash
upme health                                   # health check

upme skill list                               # list all skills
upme skill get    -id <id>                     # show one skill (incl. SKILL.md)
upme skill import -url <github/raw url> [-name <n>]   # import from a URL
upme skill import -file <path>          [-name <n>]   # import from a local file
upme skill new    -name <n> -corpus <text|@file>     # generate from corpus (LLM)
upme skill rm     -id <id>                      # delete a skill

upme chat -skill <id> -m "<message>" [-conv <id>]    # chat with your twin
```

---

## 📥 Importing a Skill

Bring a ready-made `SKILL.md` into the system without re-running the LLM:

```
POST /api/v1/skills/import
{ "url": "https://github.com/owner/repo/blob/main/path/SKILL.md", "name": "optional" }
# or
{ "content": "---\nname: my-skill\n---\n...", "name": "optional" }
```

- Provide **either** `url` **or** `content`.
- GitHub `blob` web URLs are automatically rewritten to their `raw.githubusercontent.com` form.
- If `name` is omitted, it is taken from the `name:` field of the YAML frontmatter
  (the Claude Agent Skills format), falling back to `imported-skill`.

---

## 📁 Project Structure

```
UploadMyself/
├── backend/                      # Go backend (Gin, layered modules)
│   ├── main.go                   # Entry point: DI wiring, pgxpool, embedded migrations
│   ├── migrations/
│   │   └── 001_init.sql          # DDL, embedded via go:embed & run at startup
│   ├── internal/
│   │   └── llm/client.go         # OpenAI-compatible LLM client
│   ├── pkg/                      # One module per table, six layers each
│   │   ├── tasks/                # ← module template (entity/dto/mapper/service/service/impl/handler)
│   │   │   ├── entity/           #   table struct
│   │   │   ├── dto/              #   request/response objects
│   │   │   ├── mapper/           #   pgx data access
│   │   │   ├── service/          #   business-logic interface
│   │   │   ├── service/impl/     #   interface implementation
│   │   │   └── handler/          #   HTTP handlers + route Register
│   │   ├── skills/               # Thinking-framework / SKILL.md generation
│   │   ├── voices/               # Voice cloning (train / synthesize)
│   │   ├── avatars/              # 2D/3D avatar processing
│   │   ├── file_uploads/         # File & corpus uploads (OCR/PDF/Word)
│   │   └── messages/             # Agent chat, conversation messages, tool registry
│   ├── cmd/
│   │   └── upme/                 # `upme` CLI client for the REST API
│   └── middleware/               # CORS, etc.
├── ml/                           # ML models & Python scripts
│   ├── models/                   # Pretrained weights
│   └── scripts/                  # Preprocessing & training (ml_service.py)
├── frontend/                     # React + Three.js
│   └── src/
│       ├── pages/                # UI pages
│       └── three/                # 3D rendering
├── skills/                       # Skill packages — <id>/SKILL.md (+ meta.json) per skill
├── docs/                         # Documentation
└── tests/                        # Tests
```

> **Skill storage** — generated and imported thinking frameworks are stored as
> files, not blobs in the DB. Each skill is a directory `<SKILLS_DIR>/<id>/`
> containing `SKILL.md` (the framework) and `meta.json` (name/source/timestamp).
> The Postgres `skills` table is only a metadata index (id, name, status…); the
> `SKILL.md` body is read from disk. `SKILLS_DIR` defaults to `./skills` and is
> mounted as a Docker volume so skills survive container restarts.

> **Module layering** — every package under `backend/pkg/` follows the same six-layer
> shape, with `backend/pkg/tasks` as the reference template:
> `entity` (table struct) → `dto` (I/O objects) → `mapper` (pgx data access) →
> `service` (interface) → `service/impl` (implementation) → `handler` (HTTP + routes).
> Dependencies are wired in `main.go` as `mapper → service → handler`. The legacy
> `backend/agent`, `backend/handler`, and `backend/model` packages have been removed.

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

- **后端**：Golang (Gin + pgx/pgxpool + Viper + Zap，分层模块结构)
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

### CLI (upme)

`upme` 是后端 REST API 的命令行客户端。后端地址默认 `http://localhost:8000`，
可用 `-server` 或环境变量 `$UPME_SERVER` 覆盖。

```bash
# 构建
cd backend && go build -o upme ./cmd/upme

upme health                                          # 健康检查
upme skill list                                      # 列出所有思维框架
upme skill get    -id <id>                           # 查看单个 skill（含 SKILL.md）
upme skill import -url <url> [-name <n>]             # 从 URL/GitHub 导入
upme skill import -file <path> [-name <n>]           # 从本地文件导入
upme skill new    -name <n> -corpus <text|@file>     # 用语料生成（触发 LLM）
upme skill rm     -id <id>                            # 删除
upme chat -skill <id> -m "<消息>" [-conv <id>]        # 与分身对话
```

### 导入 Skill

无需重新跑 LLM，即可导入现成的 `SKILL.md`：

```
POST /api/v1/skills/import
{ "url": "https://github.com/owner/repo/blob/main/path/SKILL.md", "name": "可选" }
# 或
{ "content": "---\nname: my-skill\n---\n...", "name": "可选" }
```

- `url` 与 `content` 二选一。
- GitHub `blob` 网页地址会自动转换为 `raw.githubusercontent.com` 原始地址。
- 不传 `name` 时，从 YAML frontmatter 的 `name:` 字段解析（Claude Agent Skills 格式），
  兜底为 `imported-skill`。

> **Skill 文件存储**：生成/导入的思维框架以文件形式保存在
> `<SKILLS_DIR>/<id>/SKILL.md`（+ `meta.json`），数据库 `skills` 表仅做元数据索引。
> `SKILLS_DIR` 默认 `./skills`，Docker 中挂载为 `skills` 卷。

### 开发进度

详见 [TODO.md](TODO.md) — 包含 7 个阶段、80+ 个子任务的技术路线图。

### 文档

- [架构设计](docs/architecture.md)
- [API 文档](docs/api.md)
- [模型说明](docs/models.md)
