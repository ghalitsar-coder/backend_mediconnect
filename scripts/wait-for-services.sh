#!/usr/bin/env bash
# scripts/wait-for-services.sh
# Polls TCP ports until all required services are accepting connections.
# Usage: ./scripts/wait-for-services.sh

set -euo pipefail

TIMEOUT=${WAIT_TIMEOUT:-60}

wait_for() {
    local name=$1
    local host=$2
    local port=$3
    local elapsed=0

    echo "⏳  Waiting for ${name} at ${host}:${port}..."
    until nc -z "${host}" "${port}" 2>/dev/null; do
        if [ "${elapsed}" -ge "${TIMEOUT}" ]; then
            echo "❌  Timed out waiting for ${name} after ${TIMEOUT}s"
            exit 1
        fi
        sleep 2
        elapsed=$((elapsed + 2))
    done
    echo "✅  ${name} is ready."
}

wait_for "PostgreSQL" "${DB_HOST:-localhost}" "${DB_PORT:-5432}"
wait_for "Redis"      "${REDIS_HOST:-localhost}" "${REDIS_PORT:-6379}"
wait_for "RabbitMQ"   "${RABBITMQ_HOST:-localhost}" "${RABBITMQ_PORT:-5672}"

echo "🚀  All services ready. Starting MediConnect API..."
exec "$@"
