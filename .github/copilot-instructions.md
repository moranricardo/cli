# Dependabot Lite - Copilot Instructions (moranricardo/cli)

## What is this?
Fork reescrito: orquestador NATIVO sin Docker. No levanta Proxy/Updater.
Go 1.23+ STABLE, optimizado para entornos con recursos limitados (linux-armv7 compatible).

## Commands reales
- dependabot update [ecosystem] [dir] - único comando
- Flags: -v verbose, -h help
- NO existe: test, graph, version, --version

## Layout actual
- cmd/dependabot/ - entrypoint Cobra lite
- internal/infra/ - run.go simplificado sin Docker

## Relación ra-pulse
Consumido como binario directo por ra-pulse-orchestrator. No cambiar firma de update sin bump major.

## Security & Compliance
- SECURITY.md: private advisory first, 90d disclosure
- Integridad: siempre commitear go.sum

## Build
go build -o dependabot ./cmd/dependabot
./dependabot update --help
go vet ./...
