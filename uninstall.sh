#!/usr/bin/env bash
# 卸载 UploadMyself：停止服务，可选删除工作目录与 Docker 卷/镜像
#   ./uninstall.sh [--purge]   --purge 同时删除工作目录数据
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOME_DIR="${HOME_DIR:-$HOME/.UploadMyself}"
PURGE=0
[[ "${1:-}" == "--purge" ]] && PURGE=1

log()  { printf '\033[36m[uninstall]\033[0m %s\n' "$*"; }
ok()   { printf '\033[32m[ ok ]\033[0m %s\n' "$*"; }
warn() { printf '\033[33m[warn]\033[0m %s\n' "$*" >&2; }

log "停止所有服务..."
bash "$REPO_DIR/stop.sh" all || true

# Docker 卷/容器清理
if [[ -f "$REPO_DIR/.env" ]] && command -v docker >/dev/null 2>&1; then
  if docker compose version >/dev/null 2>&1; then
    ( cd "$REPO_DIR" && docker compose --env-file .env down -v --remove-orphans ) || true
  fi
fi

if [[ "$PURGE" == 1 ]]; then
  read -rp "确认删除工作目录 $HOME_DIR 及其全部数据？此操作不可恢复 [y/N]: " ans || true
  if [[ "${ans:-N}" =~ ^[Yy]$ ]]; then
    rm -rf "$HOME_DIR"
    ok "已删除 $HOME_DIR"
  else
    warn "已跳过删除工作目录"
  fi
else
  warn "工作目录 $HOME_DIR 已保留（加 --purge 可一并删除）"
fi

ok "卸载完成"
