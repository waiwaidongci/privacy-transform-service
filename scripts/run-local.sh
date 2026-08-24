#!/usr/bin/env bash
set -euo pipefail
PRIVACY_TRANSFORM_HTTP_ADDR="${PRIVACY_TRANSFORM_HTTP_ADDR:-:8083}" go run ./cmd/privacy-transform-service
