# UploadMyself — 安装指南 / Installation Guide

运行时数据统一存放在工作目录 **`~/.UploadMyself`**（隐藏目录），包含数据库、上传文件、skill 包与 ML 模型缓存。

Runtime data lives in the working directory **`~/.UploadMyself`** (hidden): database, uploads, skill packages, and the ML model cache.

---

## 一键安装 / One-command install

```bash
git clone https://github.com/Bruce-Sakura/UploadMyself.git
cd UploadMyself

# 交互选择 Docker / 原生 (interactive)
bash install.sh

# 或显式指定 (explicit)
bash install.sh --docker      # Docker 模式（推荐 / recommended）
bash install.sh --native      # Linux 原生模式 / native

# 常用选项 / options
bash install.sh --docker --llm-key sk-xxx --home ~/.UploadMyself -y
```

安装脚本会：创建工作目录 → 生成 `.env` → 构建并启动服务。
The installer creates the working dir, generates `.env`, then builds & starts everything.

---

## 方式一：Docker（推荐）/ Option A: Docker

**前置条件 / Prerequisites**
- Docker + Docker Compose v2
- 全功能（3D/语音）需要 **NVIDIA GPU + [nvidia-container-toolkit](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/)**
  Full features (3D/voice) require an NVIDIA GPU and the NVIDIA Container Toolkit.

```bash
bash install.sh --docker
```

- 前端 Frontend: http://localhost:3000
- 后端 Backend:  http://localhost:8000
- ML 服务 ML:    http://localhost:8001

数据卷以 bind mount 落到 `~/.UploadMyself/data/{postgres,redis,uploads,skills,models}`。
HuggingFace 模型缓存持久化到 `~/.UploadMyself/data/models`，重启不重复下载。

---

## 方式二：Linux 原生 / Option B: Linux native

**前置条件 / Prerequisites**
- Linux (Debian/Ubuntu，使用 apt) / Debian/Ubuntu-based Linux
- NVIDIA 驱动 + CUDA 运行时（全功能）/ NVIDIA driver + CUDA runtime
- 具备 sudo 权限（安装系统依赖）/ sudo privileges

```bash
bash install.sh --native
```

脚本自动安装：git/curl/ffmpeg/PostgreSQL、Go 1.22、Node 20、Miniconda，
创建 conda 环境 `charactergen`（PyTorch cu121 + CharacterGen 依赖），
编译后端到 `~/.UploadMyself/uploadmyself`，构建前端。

The script installs system deps, Go 1.22, Node 20, Miniconda, creates the
`charactergen` conda env, builds the backend binary and the frontend.

- 前端 Frontend (vite dev): http://localhost:5173
- 后端 Backend:  http://localhost:8000
- ML 服务 ML:    http://localhost:8001

---

## 启动 / 停止 / 卸载 — Run / Stop / Uninstall

```bash
bash run.sh docker     # 启动（Docker）/ start (docker)
bash run.sh native     # 启动（原生）/ start (native)

bash stop.sh           # 停止（自动识别两种模式）/ stop both
bash stop.sh docker
bash stop.sh native

bash uninstall.sh          # 停止 + 清理容器/卷（保留数据目录）
bash uninstall.sh --purge  # 同时删除 ~/.UploadMyself 全部数据
```

---

## 配置 / Configuration

Docker 模式配置在仓库根 `.env`；原生模式在 `~/.UploadMyself/.env`。
Docker config: repo-root `.env`. Native config: `~/.UploadMyself/.env`.

| 变量 / Variable | 说明 / Description |
|-----------------|-------------------|
| `UPLOADMYSELF_HOME` | 工作目录绝对路径 / working dir |
| `LLM_API_KEY` | LLM 密钥（MiMo / OpenAI 兼容）/ LLM key |
| `LLM_BASE_URL` | LLM 接口地址，默认 MiMo / LLM base url |
| `LLM_MODEL` | 模型名，默认 `mimo-v2.5-pro` / model |
| `DB_PATH` | SQLite 数据库文件路径 / SQLite db file path |
| `UPLOAD_DIR` `SKILLS_DIR` | 上传 / skill 包目录 / dirs |
| `ML_SERVICE_URL` | ML 服务地址 / ML service url |
| `HF_HOME` | HuggingFace 模型缓存 / model cache |

改完配置后重启：`bash stop.sh && bash run.sh <mode>`。
After editing, restart with `bash stop.sh && bash run.sh <mode>`.

---

## 常见问题 / FAQ

- **没有 GPU？/ No GPU?** 文字对话与 skill 导入仍可用；3D 形象与语音克隆需要 GPU。
  Text chat and skill import still work; 3D avatar and voice cloning need a GPU.
- **首次启动很慢？/ Slow first start?** Docker 首次需拉镜像 + 编译；原生首次需下载模型（缓存到 `~/.UploadMyself/data/models`）。
- **改了端口？/ Custom ports?** 编辑对应 `.env` 的 `APP_PORT` 与 `docker-compose.yml` 的 ports。
- **CLI 工具 / CLI:** 安装脚本会自动把 `upme` 装到 `~/.local/bin/upme`，安装后可直接 `upme skill list`、`upme chat -m "你好"`、`upme skill import -url <github>` 等（确保 `~/.local/bin` 在 `PATH`）。Docker 在 Windows/macOS 宿主机上改用 `docker compose exec backend upme ...`。
