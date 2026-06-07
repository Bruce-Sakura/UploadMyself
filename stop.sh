#!/usr/bin/env bash
# 停止 UploadMyself（不带参数则两种模式都尝试停止）
#   ./stop.sh [docker|native]
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MODE="${1:-all}"
HOME_DIR="${HOME_DIR:-$HOME/.UploadMyself}"

log()  { printf '\033[36m[stop]\033[0m %s\n' "$*"; }
ok()   { printf '\033[32m[ ok ]\033[0m %s\n' "$*"; }

stop_docker() {
  if [[ -f "$REPO_DIR/.env" ]] && command -v docker >/dev/null 2>&1; then
    local COMPOSE=""
    if docker compose version >/dev/null 2>&1; then COMPOSE="docker compose";
    elif command -v docker-compose >/dev/null 2>&1; then COMPOSE="docker-compose"; fi
    if [[ -n "$COMPOSE" ]]; then
      log "停止容器..."
      ( cd "$REPO_DIR" && $COMPOSE --env-file .env down ) || true
    fi
  fi
}

stop_native() {
  local p f pid
  for p in frontend backend ml; do
    f="$HOME_DIR/logs/$p.pid"
    if [[ -f "$f" ]]; then
      pid="$(cat "$f")"
      if kill "$pid" 2>/dev/null; then log "已停止 $p (pid $pid)"; fi
      rm -f "$f"
    fi
  done
}

case "$MODE" in
  docker) stop_docker;;
  native) stop_native;;
  all) stop_native; stop_docker;;
  *) echo "用法: ./stop.sh [docker|native]"; exit 1;;
esac
ok "已停止"
