#!/usr/bin/env bash
#
# 服务器侧部署脚本：链路里的第 5、6 步。
# 由 GitHub Actions 通过 SSH 调用，也可以手动执行。
#
# 用法：
#   bash scripts/deploy.sh <image_tag>   部署指定版本（tag 是 commit sha 前 12 位）
#   bash scripts/deploy.sh --status      查看当前版本与容器状态
#
# 做的事：拉新镜像 -> 重建 app 容器 -> 轮询 /health -> 失败则切回上一个版本。

set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_DIR"

COMPOSE_FILE="docker/compose.prod.yaml"
ENV_FILE="docker/.env"
STATE_FILE=".deploy_state"

log() { printf '[%s] %s\n' "$(date '+%F %T')" "$*"; }

compose() {
  docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" "$@"
}

# 从 .env 读一个变量；读不到就用兜底值。用它可以避免把端口写死在两处。
read_env() {
  local key="$1" fallback="$2" value=""
  if [ -f "$ENV_FILE" ]; then
    value="$(grep -E "^${key}=" "$ENV_FILE" | tail -n1 | cut -d= -f2- | tr -d '\r' || true)"
  fi
  printf '%s' "${value:-$fallback}"
}

HEALTH_URL="http://127.0.0.1:$(read_env PUBLIC_PORT 8080)/health"

# 最多探 30 次、每次间隔 2 秒，即最长等 60 秒。
# 这个时间的意义是：把「发布成功」定义为应用真的能响应，而不是容器启动了。
wait_health() {
  local attempt
  for attempt in $(seq 1 30); do
    if curl -fsS --max-time 3 "$HEALTH_URL" >/dev/null 2>&1; then
      log "健康检查通过（第 ${attempt} 次探测）"
      return 0
    fi
    sleep 2
  done
  return 1
}

if [ "${1:-}" = "--status" ]; then
  log "当前版本：$(cat "$STATE_FILE" 2>/dev/null || echo '（无记录）')"
  compose ps
  exit 0
fi

NEW_TAG="${1:-}"
if [ -z "$NEW_TAG" ]; then
  log "用法：bash scripts/deploy.sh <image_tag>"
  exit 2
fi

if [ ! -f "$ENV_FILE" ]; then
  log "缺少 ${ENV_FILE}，请先从 docker/.env.example 复制并填写"
  exit 2
fi

PREV_TAG="$(cat "$STATE_FILE" 2>/dev/null || true)"
log "当前版本：${PREV_TAG:-（无记录）}"
log "开始部署：${NEW_TAG}"

export IMAGE_TAG="$NEW_TAG"

log "拉取新镜像（数据库容器不动）"
compose pull app

log "重建 app 容器"
compose up -d app

if wait_health; then
  # 只有健康检查通过才把新版本记为「当前版本」，
  # 这样下一次部署的回滚目标一定是最后一个真正跑起来的版本。
  printf '%s\n' "$NEW_TAG" > "$STATE_FILE"
  log "部署成功，当前版本：${NEW_TAG}"
  exit 0
fi

log "健康检查连续 30 次未通过，判定本次发布失败"
compose logs --tail=80 app || true

if [ -z "$PREV_TAG" ]; then
  log "没有可回滚的历史版本，需要人工介入"
  exit 1
fi

log "回滚到上一个版本：${PREV_TAG}"
export IMAGE_TAG="$PREV_TAG"
compose up -d app

if wait_health; then
  log "已回滚到 ${PREV_TAG}，服务已恢复"
else
  log "回滚后健康检查仍未通过，需要人工介入"
fi

# 无论如何都以失败退出，让 GitHub Actions 上这次发布显示为红色
exit 1
