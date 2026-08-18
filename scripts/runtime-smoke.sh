#!/usr/bin/env bash
set -Eeuo pipefail

# 真实启动 Compose backend；探活通过后由 runtime_smoke 终止本脚本，trap 负责清理。
export COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-runtime_a_contract_ticket_api}"
# contractapi does not publish the database port; DB_PORT is used for the internal backend connection and is left unchanged.
compose=(docker compose --env-file .env.example)
cleanup() {
  if [[ -n "${logs_pid:-}" ]]; then kill "$logs_pid" >/dev/null 2>&1 || true; fi
  "${compose[@]}" down --volumes --remove-orphans >/dev/null 2>&1 || true
}
trap 'cleanup; exit 143' TERM INT
trap cleanup EXIT
cleanup
"${compose[@]}" up --build --no-color --detach backend
"${compose[@]}" logs --follow --no-color backend &
logs_pid=$!
wait "$logs_pid"
