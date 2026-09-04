# CLI Debugging Guide (Lite v0.5.3)

Guía para depurar el CLI en modo nativo Go sin Docker, optimizado para dispositivos de recursos limitados.

## 1. Build nativo ARM 32-bit ultra-ligero

Build reproducible sin cgo para entornos low-resource:

env -u GOOS -u GOARCH -u GOARM CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -tags netgo,osusergo -ldflags="-s -w" -trimpath -o dependabot ./cmd/dependabot

./dependabot --help
# binario ~2.8M ARM 32-bit

## 2. Verificación rápida (sin Docker)

DEPENDABOT_LOCAL_RUN=1 E6_LOW_MODE=1 ./dependabot update go_modules . --local-run -o /tmp/out.yml -v
head -n 40 /tmp/out.yml

Debe generar YAML con `create_pull_request`.

## 3. Runner Nativo Lite (Diamond)

export DEPENDABOT_LOCAL_RUN=1
export E6_LOW_MODE=1
go run ./cmd/dependabot update go_modules ./ --local-run -v

Arquitectura:
- internal/infra/run-lite.go: aislamiento en /tmp, timeout 10s, error propagation estricta
- internal/server/: proxy local API :40505

## 4. Diagnóstico de bloqueos

pkill -ABRT dependabot
ls /tmp/dependabot-*
