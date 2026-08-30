# Maintainer Docs (Lite Edition)

Guía de publicación y mantenimiento de versiones para el CLI nativo ligero.

## Crear un Release Local / Tag

1. Verificar el árbol de trabajo:
   git status

2. Compilar y probar el binario nativo:
   go build -o dependabot ./cmd/dependabot
   ./dependabot --version

3. Crear el tag anotado de la versión:
   git tag -a v0.X.Y-lite -m "release: versión v0.X.Y-lite nativa"

4. Publicar cambios y tags en GitHub:
   git push origin main --tags

## Manejo de errores en versiones

- Eliminar tag local: git tag -d v0.X.Y-lite
- Eliminar tag remoto: git push origin :refs/tags/v0.X.Y-lite
