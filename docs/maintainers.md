## Maintainer docs (Fork Lite v0.5)

This documentation is targeting maintainers of the Lite fork.

### Creating a lite release

1. Limpia el árbol y verifica versión:
   git status
   go build -o dependabot ./cmd/dependabot && ./dependabot --version

2. Tag local lite:
   git tag -a v0.5.2-lite -m "fix: native go runner no storage pollution"

3. Push:
   git push origin main --tags

4. Release en GitHub:
   - Go to https://github.com/moranricardo/cli/releases
   - Draft a new release -> Choose tag `v0.5.x-lite`
   - Generate release notes, borra cambios menores de workflows/README
   - Publish. No dependes del workflow `release.yml` oficial, tu binario es nativo.

### Nota de versionado
El binario genera su versión basándose en `cmd/dependabot/internal/cmd/version.go` + `git describe`.
Si la salida incluye el sufijo `.dirty`, limpia el working tree antes de crear el tag.

### Rollback de tags
- Eliminar local: `git tag -d v0.5.x-lite`
- Eliminar remoto: `git push origin :refs/tags/v0.5.x-lite`

### Checklist antes de push
- [ ] `find cmd -type f -name "*.go" | sort` solo 4 archivos (sin `cmd/dependabot/cmd/`)
- [ ] `internal/infra/run-lite.go` no escribe en storage no deseado
- [ ] `docs/debugging.md` menciona modo `--local --debug`
- [ ] `./dependabot --version` no muestra el sufijo `.dirty`
