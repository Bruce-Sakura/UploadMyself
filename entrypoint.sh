#!/bin/bash
set -e

echo "[entrypoint] Starting ML service on port 8001..."
cd /app
python3 ml/scripts/ml_service.py --port 8001 &
ML_PID=$!

# Wait for ML service to be ready
for i in $(seq 1 30); do
    if curl -s http://localhost:8001/health > /dev/null 2>&1; then
        echo "[entrypoint] ML service ready ✅"
        break
    fi
    sleep 2
done

echo "[entrypoint] Starting Go backend on port 8000..."
exec ./uploadmyself
