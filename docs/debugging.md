# CLI Debugging Guide (Lite v0.5+)

Este documento explica cómo depurar ejecuciones del CLI usando el entorno nativo Go y el flujo `v0.5-lite` en dispositivos de recursos limitados.

## 1. Verificación rápida (Nativa, sin Docker)

Prueba que el CLI pueda evaluar las dependencias localmente y generar el archivo de salida:

./dependabot update go_modules dependabot/cli -o out.yml
head -n 40 out.yml

Deberías ver la estructura YAML con las llamadas `create_pull_request`.

## 2. Depuración del Runner Nativo Go (Lite / Low-resource)

Para rastrear la ejecución interna sin contenedor ni Docker, lanza el runner en modo local con logs extendidos:

go run ./cmd/dependabot update go_modules ./ --local --debug

Puntos clave de la arquitectura a inspeccionar:
- internal/infra/run-lite.go: Controla el aislamiento en /tmp y evita modificar el repositorio base.
- internal/server/: Proxy local API HTTP que simula el backend de GitHub/Dependabot.

## 3. Depuración del Updater (dependabot-core)

Si necesitas depurar los parsers de Ruby desde la fuente original:

git clone https://github.com/dependabot/dependabot-core
cd dependabot-core
script/dependabot update go_modules dependabot/cli --debug

Dentro de la sesión interactiva:
bin/run fetch_files
bin/run update_files

### Uso del Depurador de Ruby (rdbg)
Inserta la sentencia `debugger` en el código de dependabot-core (por ejemplo, en `go_modules/lib/dependabot/go_modules/update_checker.rb`):

def latest_resolvable_version
  debugger
  latest_version_finder.latest_version
end

Al ejecutar `bin/run update_files`, la consola se detendrá en esa línea para inspeccionar variables como `dependency`.

## 4. Diagnóstico de bloqueos (Hangs)

Si el proceso se congela durante una actualización:
- En ejecuciones locales de Go, genera un dump del stack trace enviando la señal ABRT:
  pkill -ABRT dependabot
- Revisa las peticiones enviadas al proxy local observando la salida HTTP del servidor local (puerto por defecto :40505).
