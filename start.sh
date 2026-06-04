#!/bin/bash
# UploadMyself 启动脚本
# 启动 Docker 服务 + 主机 ML 服务

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

echo "=== UploadMyself 启动 ==="

# 1. 启动 Docker 服务 (PostgreSQL + Redis + Go Backend + Frontend)
echo "[1/3] 启动 Docker 服务..."
docker compose up -d

# 等待 PostgreSQL 就绪
echo "[2/3] 等待 PostgreSQL..."
sleep 3

# 2. 启动 ML 服务 (主机 GPU)
echo "[3/3] 启动 ML 服务 (GPU)..."
source ~/miniconda3/etc/profile.d/conda.sh
conda activate charactergen

# 后台启动 ML 服务
nohup python ml/scripts/ml_service.py --port 8001 > /tmp/uploadmyself-ml.log 2>&1 &
ML_PID=$!
echo "ML 服务 PID: $ML_PID"

# 等待 ML 服务就绪
sleep 5
if curl -s http://localhost:8001/health > /dev/null 2>&1; then
    echo "✅ ML 服务就绪"
else
    echo "⚠️ ML 服务启动中，查看日志: tail -f /tmp/uploadmyself-ml.log"
fi

echo ""
echo "=== 服务就绪 ==="
echo "  前端: http://localhost:3000"
echo "  后端: http://localhost:8000"
echo "  ML:   http://localhost:8001"
echo ""
echo "停止: ./stop.sh"
