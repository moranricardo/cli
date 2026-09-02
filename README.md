# cli - sistema interno Moto E6 Plus Android 9

Fork de dependabot/cli adaptado por @moranricardo para dispositivos de baja gama.

> Termux X es solo desplazador de texto, no parte del proyecto.

## Función
Herramienta que orquesta jobs de Dependabot sin Docker en modo lite para Android 9.

## Instalación sistema interno
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dependabot-internal./cmd/dependabot

## Uso
./dependabot-internal --help
./dependabot-internal update go_modules moranricardo/cli --provider github
./dependabot-internal test -f testdata/smoke.yml

## Modo Lite (Moto E6)
El binario interno usa `internal/infra/run-lite.go` sin Docker:
- Menos de 10MB
- Sin daemon Docker
- Compatible con shell Android estándar (am, logcat)

## Desarrollo
go test./... -count=1 -short
go vet./...

## Licencia
MIT - Copyright 2022 GitHub, Inc.
