## Maintainer docs (Fork Lite v0.5.2)

This documentation is targeting maintainers of the Lite fork.

### Creating a lite release

1. Limpia el árbol:
   git status
   go build -o dependabot ./cmd/dependabot && ./dependabot --version

2. Test rápido sin modificar nada:
   DRY_RUN=1 ./dependabot update go_modules --local-run .

3. Tag local lite:
   git tag -a v0.5.2-lite -m "release: v0.5.2-lite with strict error propagation and http timeout"

4. Push:
   git push origin main --tags

5. Release en GitHub:
   - Go to https://github.com/moranricardo/cli/releases
   - Draft a new release -> Choose tag v0.5.2-lite
   - Generate release notes, borra cambios menores
   - Publish. No dependes del workflow release.yml oficial.

### Checklist antes de push
- [ ] `find cmd -type f -name "*.go" | sort` solo 4 archivos
- [ ] `internal/infra/run-lite.go` con timeout 10s y error propagation
- [ ] `docs/debugging.md` menciona --local --debug
- [ ] `./dependabot --version` sale v0.5.2-lite limpio (sin +dirty)
