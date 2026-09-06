# CLI - Dependabot Lite
Native Go, CGO_ENABLED=0. No Docker.

## Relación ra-pulse
Consumido como binario directo por ra-pulse-orchestrator. No cambiar firma de update sin bump major.

## Firma estable
dependabot update [ecosystem] [dir] --verbose
Requiere DEPENDABOT_LOCAL_RUN=1

## Estructura
- cmd/dependabot/ - entrypoint
- internal/infra/ - run.go sin Docker
- internal/model/ - models

## Security & Compliance
No incluir rutas absolutas de entorno local en logs.
