#!/bin/sh
# wait-for-backend.sh
# 等待 web-server 的 /api/v1/healthz 返回 200，最多等待 60 次（约 60s）
set -e
HOST="http://web-server:8080/api/v1/healthz"
MAX_RETRIES=60
SLEEP=1
count=0
while [ $count -lt $MAX_RETRIES ]; do
  if curl -fsS "$HOST" >/dev/null 2>&1; then
    echo "backend is up"
    exit 0
  fi
  echo "waiting for backend... ($count)"
  count=$((count+1))
  sleep $SLEEP
done
echo "backend did not become ready in time"
exit 1
