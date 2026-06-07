#!/usr/bin/env bash
# 启动 UploadMyself
#   ./run.sh docker   # Docker 模式（默认）
#   ./run.sh native   # Linux 原生模式
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MODE="${1:-docker}"
HOME_DIR="${HOME_DIR:-$HOME/.UploadMyself}"

log()  { printf '\033[36m[run]\033[0m %s\n' "$*"; }
ok()   { printf '\033[32m[ ok ]\033[0m %s\n' "$*"; }
err()  { printf '\033[31m[error]\033[0m %s\n' "$*" >&2; }

compose_cmd() {
  if docker compose version >/dev/null 2>&1; then echo "docker compose";
  elif command -v docker-compose >/dev/null 2>&1; then echo "docker-compose";
  else err "未找到 docker compose"; exit 1; fi
}

run_docker() {
  [[ -f "$REPO_DIR/.env" ]] || { err "缺少 .env，请先运行 ./install.sh --docker"; exit 1; }
  local COMPOSE; COMPOSE="$(compose_cmd)"
  log "启动容器..."
  ( cd "$REPO_DIR" && $COMPOSE --env-file .env up -d )
  ok "已启动"
  cat <<EOF
  前端: http://localhost:3000
  后端: http://localhost:8000
  ML:   http://localhost:8001
EOF
}

run_native() {
  local envf="$HOME_DIR/.env"
  [[ -f "$envf" ]] || { err "缺少 $envf，请先运行 ./install.sh --native"; exit 1; }
  # 导出配置给后端/ML（viper.AutomaticEnv 读环境变量）
  set -a; # shellcheck disable=SC1090
  source "$envf"; set +a
  mkdir -p "$HOME_DIR/logs"

  log "启动 ML 服务 (:8001)..."
  nohup "${PYTHON_BIN:-python3}" "$REPO_DIR/ml/scripts/ml_service.py" --port 8001 \
    > "$HOME_DIR/logs/ml.log" 2>&1 &
  echo $! > "$HOME_DIR/logs/ml.pid"

  log "启动后端 (:8000)..."
  pushd "$HOME_DIR" >/dev/null
  nohup ./uploadmyself > "$HOME_DIR/logs/backend.log" 2>&1 &
  echo $! > "$HOME_DIR/logs/backend.pid"
  popd >/dev/null

  log "启动前端 dev (:5173)..."
  pushd "$REPO_DIR/frontend" >/dev/null
  nohup npm run dev > "$HOME_DIR/logs/frontend.log" 2>&1 &
  echo $! > "$HOME_DIR/logs/frontend.pid"
  popd >/dev/null

  ok "已启动（日志在 $HOME_DIR/logs/）"
  cat <<EOF
  前端: http://localhost:5173
  后端: http://localhost:8000
  ML:   http://localhost:8001
EOF
}

case "$MODE" in
  docker) run_docker;;
  native) run_native;;
  *) err "未知模式: $MODE（用 docker 或 native）"; exit 1;;
esac
