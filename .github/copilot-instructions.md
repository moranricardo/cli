# deja solo esto:# Dependabot Lite - 
# Copilot Instructions (moranricardo/cli)
# ## Relación ra-pulse
# Consumido como binario directo por 
# ra-pulse-orchestrator...## What is this?
# ## RunnerFork reescrito: orquestador 
# ## NATIVO sin Docker. No levanta 
# ## Proxy/Updater.
# Native Go, CGO_ENABLED=0, linux-armv7 
# compatibleGo 1.23+ STABLE, optimizado 
# para entornos con recursos limitados 
# (linux-armv7 compatible).

## Commands reales
# 4. quita el binario compilado del repo 
# (no se debe subir)- dependabot update 
# [ecosystem] [dir] - único comando
rm -f dependabot- Flags: -v verbose, -h 
help - NO existe: test, graph, version, 
--version echo "dependabot" >> .gitignore 
echo "/data/" >> .gitignore## Layout 
actual - cmd/dependabot/ - entrypoint 
Cobra lite - internal/infra/ - run.go 
simplificado sin Docker

## Relación ra-pulse
# 5. commit limpioConsumido como binario 
# directo por ra-pulse-orchestrator. No 
# cambiar firma de update sin bump major.
git add .gitignore 
.github/copilot-instructions.md git 
status## Security & Compliance - 
SECURITY.md: private advisory first, 90d 
disclosure - Integridad: siempre commitear 
go.sum git commit -m "chore: excluye 
compatible"
git push## Build
go build -o dependabot ./cmd/dependabot
./dependabot update --help
go vet ./...
