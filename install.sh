#!/usr/bin/env bash
# UploadMyself 一键安装脚本
#   Docker 模式（推荐，跨平台）:  ./install.sh --docker
#   Linux 原生模式（全功能含 GPU）: ./install.sh --native
#   不带参数则交互选择。运行时数据统一存放在工作目录（默认 ~/UploadMyself）。
set -euo pipefail

# ---------- defaults ----------
MODE=""
HOME_DIR="${HOME}/.UploadMyself"
LLM_KEY=""
LLM_BASE="https://token-plan-sgp.xiaomimimo.com/v1"
LLM_MODEL="mimo-v2.5-pro"
ASSUME_YES=0
REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ---------- helpers ----------
log()  { printf '\033[36m[install]\033[0m %s\n' "$*"; }
ok()   { printf '\033[32m[ ok ]\033[0m %s\n' "$*"; }
warn() { printf '\033[33m[warn]\033[0m %s\n' "$*" >&2; }
err()  { printf '\033[31m[error]\033[0m %s\n' "$*" >&2; }

path_hint() { # path_hint <dir>: 若 dir 不在 PATH 则提示加入
  case ":${PATH}:" in
    *":$1:"*) : ;;
    *) warn "将 CLI 目录加入 PATH:  echo 'export PATH=\"$1:\$PATH\"' >> ~/.bashrc && source ~/.bashrc" ;;
  esac
}

usage() {
  cat <<EOF
UploadMyself 一键安装

用法: ./install.sh [选项]

选项:
  --docker            Docker 模式安装（推荐）
  --native            Linux 原生模式安装（全功能含 GPU）
  --home <dir>        工作目录（默认 \$HOME/.UploadMyself）
  --llm-key <key>     LLM API Key（不提供则交互询问）
  --llm-base <url>    LLM Base URL（默认 MiMo）
  --llm-model <name>  LLM 模型名（默认 mimo-v2.5-pro）
  -y, --yes           跳过确认提示
  -h, --help          显示帮助
EOF
}

# ---------- parse args ----------
while [[ $# -gt 0 ]]; do
  case "$1" in
    --docker) MODE=docker; shift;;
    --native) MODE=native; shift;;
    --home) HOME_DIR="$2"; shift 2;;
    --llm-key) LLM_KEY="$2"; shift 2;;
    --llm-base) LLM_BASE="$2"; shift 2;;
    --llm-model) LLM_MODEL="$2"; shift 2;;
    -y|--yes) ASSUME_YES=1; shift;;
    -h|--help) usage; exit 0;;
    *) err "未知参数: $1"; usage; exit 1;;
  esac
done

confirm() { # confirm "提示" -> 0 yes / 1 no
  [[ "$ASSUME_YES" == 1 ]] && return 0
  local ans
  read -rp "$1 [y/N]: " ans || true
  [[ "${ans:-N}" =~ ^[Yy]$ ]]
}

choose_mode() {
  [[ -n "$MODE" ]] && return
  echo "选择安装方式:"
  echo "  1) Docker（推荐，跨平台）"
  echo "  2) Linux 原生（全功能含 GPU）"
  local ans
  read -rp "输入 1/2 [1]: " ans || true
  case "${ans:-1}" in
    2) MODE=native;;
    *) MODE=docker;;
  esac
}

prompt_llm_key() {
  [[ -n "$LLM_KEY" ]] && return
  read -rp "输入 LLM API Key（MiMo / OpenAI 兼容，可留空稍后改 .env）: " LLM_KEY || true
}

create_home() {
  log "工作目录: $HOME_DIR"
  mkdir -p "$HOME_DIR"/data/{postgres,redis,uploads,skills,models} "$HOME_DIR/logs"
}

write_env() {
  local abs_home; abs_home="$(cd "$HOME_DIR" && pwd)"
  cat > "$REPO_DIR/.env" <<EOF
# 由 install.sh 生成
UPLOADMYSELF_HOME=$abs_home
APP_PORT=8000
DB_DSN=host=postgres user=uploadmyself password=uploadmyself dbname=uploadmyself port=5432 sslmode=disable
LLM_API_KEY=$LLM_KEY
LLM_BASE_URL=$LLM_BASE
LLM_MODEL=$LLM_MODEL
PYTHON_BIN=python3
ML_SCRIPTS_DIR=/app/ml/scripts
SKILLS_DIR=/app/skills
ML_SERVICE_URL=http://localhost:8001
EOF
  ok ".env 已写入 $REPO_DIR/.env"
}

print_urls() {
  cat <<EOF

=== UploadMyself 已启动 ===
  前端:  http://localhost:3000
  后端:  http://localhost:8000
  ML:    http://localhost:8001
  工作目录: $HOME_DIR
  管理:  ./run.sh [docker|native] 启动 · ./stop.sh 停止
  CLI:   upme skill list   (命令行控制；详见 README 的 upme 小节)
EOF
}

# ---------- Docker 模式 ----------
install_docker() {
  command -v docker >/dev/null 2>&1 || { err "未找到 docker，请先安装 Docker / Docker Desktop"; exit 1; }
  local COMPOSE
  if docker compose version >/dev/null 2>&1; then COMPOSE="docker compose";
  elif command -v docker-compose >/dev/null 2>&1; then COMPOSE="docker-compose";
  else err "未找到 docker compose"; exit 1; fi
  ok "Docker: $(docker --version)"

  # GPU runtime 检查（全功能需要）
  if docker info 2>/dev/null | grep -qiE 'nvidia'; then
    ok "检测到 NVIDIA Docker runtime"
  else
    warn "未检测到 Docker 的 NVIDIA runtime —— 3D/语音等 GPU 功能将不可用。"
    warn "安装方法见 https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/"
    confirm "仍要继续（仅文字/云端 LLM 功能可用）？" || exit 1
  fi

  create_home
  prompt_llm_key
  write_env

  log "构建并启动容器（首次会拉取镜像 + 编译，耗时较久）..."
  ( cd "$REPO_DIR" && $COMPOSE --env-file .env up -d --build )
  ok "容器已启动"

  # 安装 upme CLI：从镜像拷出 Linux 原生二进制到 PATH
  local bindir="$HOME/.local/bin"; mkdir -p "$bindir"
  if ( cd "$REPO_DIR" && $COMPOSE --env-file .env cp backend:/usr/local/bin/upme "$bindir/upme" ) >/dev/null 2>&1; then
    chmod +x "$bindir/upme" 2>/dev/null || true
    ok "upme CLI 已安装到 $bindir/upme"
    path_hint "$bindir"
  else
    warn "未能拷出 upme（非 Linux 宿主机时改用: $COMPOSE exec backend upme ...）"
  fi

  print_urls
}

# ---------- Linux 原生模式 ----------
install_native() {
  if [[ -x "$REPO_DIR/scripts/install_native.sh" ]]; then
    HOME_DIR="$HOME_DIR" LLM_KEY="$LLM_KEY" LLM_BASE="$LLM_BASE" LLM_MODEL="$LLM_MODEL" \
      ASSUME_YES="$ASSUME_YES" REPO_DIR="$REPO_DIR" bash "$REPO_DIR/scripts/install_native.sh"
  else
    err "原生安装脚本 scripts/install_native.sh 尚未就绪"
    exit 1
  fi
}

# ---------- main ----------
choose_mode
log "安装模式: $MODE"
case "$MODE" in
  docker) install_docker;;
  native) install_native;;
  *) err "未知模式: $MODE"; exit 1;;
esac
