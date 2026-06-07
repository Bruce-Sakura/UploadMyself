#!/usr/bin/env bash
# UploadMyself —— Linux 原生安装（全功能含 GPU ML / CharacterGen）
# 通常由根目录 install.sh --native 调用，也可独立运行。
# 接受环境变量: HOME_DIR LLM_KEY LLM_BASE LLM_MODEL ASSUME_YES REPO_DIR
set -euo pipefail

# ---------- defaults (可被 install.sh 传入的 env 覆盖) ----------
HOME_DIR="${HOME_DIR:-$HOME/.UploadMyself}"
LLM_KEY="${LLM_KEY:-}"
LLM_BASE="${LLM_BASE:-https://token-plan-sgp.xiaomimimo.com/v1}"
LLM_MODEL="${LLM_MODEL:-mimo-v2.5-pro}"
ASSUME_YES="${ASSUME_YES:-0}"
REPO_DIR="${REPO_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"

CONDA_ENV="charactergen"
GO_VERSION="1.22.5"
MINICONDA="$HOME/miniconda3"

log()  { printf '\033[36m[native]\033[0m %s\n' "$*"; }
ok()   { printf '\033[32m[ ok ]\033[0m %s\n' "$*"; }
warn() { printf '\033[33m[warn]\033[0m %s\n' "$*" >&2; }
err()  { printf '\033[31m[error]\033[0m %s\n' "$*" >&2; }

# root 则不用 sudo
if [[ "$(id -u)" -eq 0 ]]; then SUDO=""; else SUDO="sudo"; fi

require_linux() {
  [[ "$(uname -s)" == "Linux" ]] || { err "原生模式仅支持 Linux，Windows/macOS 请用 ./install.sh --docker"; exit 1; }
}

arch_tag() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64";;
    aarch64|arm64) echo "arm64";;
    *) echo "amd64";;
  esac
}

# ---------- 系统依赖 ----------
install_system_deps() {
  if ! command -v apt-get >/dev/null 2>&1; then
    warn "未检测到 apt-get（非 Debian/Ubuntu）。请手动安装: git curl ffmpeg build-essential postgresql python3 python3-pip"
    return
  fi
  log "安装系统依赖 (apt)..."
  $SUDO apt-get update -y
  $SUDO apt-get install -y --no-install-recommends \
    git curl ca-certificates ffmpeg build-essential \
    postgresql postgresql-contrib python3 python3-venv python3-pip
  ok "系统依赖就绪"
}

ensure_go() {
  local need=1
  if command -v go >/dev/null 2>&1; then
    local v; v="$(go version | grep -oE '[0-9]+\.[0-9]+' | head -1)"
    local maj=${v%%.*} min=${v##*.}
    if (( maj > 1 || (maj == 1 && min >= 22) )); then need=0; fi
  fi
  if (( need )); then
    log "安装 Go ${GO_VERSION}..."
    local tgz="go${GO_VERSION}.linux-$(arch_tag).tar.gz"
    curl -fsSL "https://go.dev/dl/${tgz}" -o "/tmp/${tgz}"
    $SUDO rm -rf /usr/local/go
    $SUDO tar -C /usr/local -xzf "/tmp/${tgz}"
    export PATH="/usr/local/go/bin:$PATH"
  fi
  ok "Go: $(go version)"
}

ensure_node() {
  local need=1
  if command -v node >/dev/null 2>&1; then
    local maj; maj="$(node -v | grep -oE '[0-9]+' | head -1)"
    (( maj >= 18 )) && need=0
  fi
  if (( need )); then
    log "安装 Node.js 20 (NodeSource)..."
    curl -fsSL https://deb.nodesource.com/setup_20.x | $SUDO -E bash -
    $SUDO apt-get install -y nodejs
  fi
  ok "Node: $(node -v)"
}

install_miniconda() {
  if command -v conda >/dev/null 2>&1; then
    ok "conda 已安装"
  elif [[ -x "$MINICONDA/bin/conda" ]]; then
    ok "conda 位于 $MINICONDA"
  else
    log "安装 Miniconda 到 $MINICONDA..."
    curl -fsSL "https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-$(uname -m).sh" -o /tmp/miniconda.sh
    bash /tmp/miniconda.sh -b -p "$MINICONDA"
  fi
  # shellcheck disable=SC1091
  source "$MINICONDA/etc/profile.d/conda.sh" 2>/dev/null || source "$(conda info --base)/etc/profile.d/conda.sh"
}

# ---------- PostgreSQL ----------
setup_postgres() {
  log "启动并初始化 PostgreSQL..."
  $SUDO service postgresql start 2>/dev/null || $SUDO systemctl start postgresql 2>/dev/null || warn "请手动启动 postgresql"
  # 幂等创建用户与库
  if ! $SUDO -u postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='uploadmyself'" | grep -q 1; then
    $SUDO -u postgres psql -c "CREATE USER uploadmyself WITH PASSWORD 'uploadmyself';"
  fi
  if ! $SUDO -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname='uploadmyself'" | grep -q 1; then
    $SUDO -u postgres psql -c "CREATE DATABASE uploadmyself OWNER uploadmyself;"
  fi
  ok "PostgreSQL: 库 uploadmyself 就绪"
}

create_home() {
  log "工作目录: $HOME_DIR"
  mkdir -p "$HOME_DIR"/data/{uploads,skills,models} "$HOME_DIR/logs"
}

# ---------- Python / ML ----------
setup_conda_env() {
  log "配置 conda 环境 $CONDA_ENV (python 3.10)..."
  if ! conda env list | grep -qE "^\s*${CONDA_ENV}\s"; then
    conda create -y -n "$CONDA_ENV" python=3.10
  fi
  set +u
  conda activate "$CONDA_ENV"
  set -u

  export HF_HOME="$HOME_DIR/data/models"
  python -m pip install --upgrade pip

  log "安装 PyTorch (CUDA 12.1)..."
  pip install torch torchvision --index-url https://download.pytorch.org/whl/cu121

  log "安装 CharacterGen + ML 依赖..."
  [[ -f "$REPO_DIR/ml/charactergen/requirements_uploadmyself.txt" ]] && \
    pip install -r "$REPO_DIR/ml/charactergen/requirements_uploadmyself.txt"
  [[ -f "$REPO_DIR/ml/scripts/requirements.txt" ]] && \
    pip install -r "$REPO_DIR/ml/scripts/requirements.txt"
  pip install gradio omegaconf einops imageio huggingface_hub
  ok "Python/ML 环境就绪 (模型缓存: $HF_HOME)"
}

# ---------- 构建 ----------
build_backend() {
  log "编译 Go 后端 -> $HOME_DIR/uploadmyself ..."
  ( cd "$REPO_DIR/backend" && go build -o "$HOME_DIR/uploadmyself" . )
  ok "后端二进制已生成"
}

path_hint() { # path_hint <dir>
  case ":${PATH}:" in
    *":$1:"*) : ;;
    *) warn "将 CLI 目录加入 PATH:  echo 'export PATH=\"$1:\$PATH\"' >> ~/.bashrc && source ~/.bashrc" ;;
  esac
}

build_upme() {
  local bindir="$HOME/.local/bin"; mkdir -p "$bindir"
  log "编译 upme CLI -> $bindir/upme ..."
  ( cd "$REPO_DIR/backend" && go build -o "$bindir/upme" ./cmd/upme )
  ok "upme CLI 已安装到 $bindir/upme"
  path_hint "$bindir"
}

build_frontend() {
  log "构建前端 (npm ci + build)..."
  ( cd "$REPO_DIR/frontend" && npm ci && npm run build )
  ok "前端 dist 已生成于 $REPO_DIR/frontend/dist"
}

write_env() {
  cat > "$HOME_DIR/.env" <<EOF
# 由 scripts/install_native.sh 生成（原生模式）
UPLOADMYSELF_HOME=$HOME_DIR
APP_PORT=8000
DB_DSN=host=localhost user=uploadmyself password=uploadmyself dbname=uploadmyself port=5432 sslmode=disable
LLM_API_KEY=$LLM_KEY
LLM_BASE_URL=$LLM_BASE
LLM_MODEL=$LLM_MODEL
PYTHON_BIN=$MINICONDA/envs/$CONDA_ENV/bin/python
ML_SCRIPTS_DIR=$REPO_DIR/ml/scripts
SKILLS_DIR=$HOME_DIR/data/skills
UPLOAD_DIR=$HOME_DIR/data/uploads
ML_SERVICE_URL=http://localhost:8001
HF_HOME=$HOME_DIR/data/models
EOF
  ok ".env 已写入 $HOME_DIR/.env"
}

main() {
  require_linux
  install_system_deps
  ensure_go
  ensure_node
  install_miniconda
  create_home
  setup_postgres
  setup_conda_env
  build_backend
  build_upme
  build_frontend
  write_env
  cat <<EOF

=== 原生安装完成 ✅ ===
  工作目录: $HOME_DIR
  启动: ./run.sh native    （前端 dev + 后端 + ML 服务）
  停止: ./stop.sh
  配置: $HOME_DIR/.env
  CLI:  upme skill list    （命令行控制；详见 README 的 upme 小节）
EOF
}

main "$@"
