## Maintainer docs (Fork Lite v0.5.3)

Target: maintainers del fork Lite low-resource.

### Creating a lite release

1. Limpia el árbol:
   git status
   env -u GOOS -u GOARCH CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -tags netgo,osusergo -ldflags="-s -w" -trimpath -o dependabot ./cmd/dependabot && ./dependabot --help

2. Test rápido sin modificar:
   DRY_RUN=1 DEPENDABOT_LOCAL_RUN=1 ./dependabot update go_modules . --local-run -o /tmp/out.yml -v

3. Tag local lite:
   git tag -a v0.5.3-lite -m "release: v0.5.3-lite ARM 32-bit 2.8M strict error"

4. Push:
   git push origin main --tags

5. Release en GitHub:
   Draft new release -> tag v0.5.3-lite -> Generate notes -> Publish.

### Checklist antes de push
- [ ] `find cmd -type f -name "*.go" | sort` solo 4 archivos
- [ ] `internal/infra/run-lite.go` con timeout 10s y error propagation
- [ ] build linux/arm 32-bit ~2.8M sin cgo
- [ ] `docs/debugging.md` menciona --local-run -v (no --debug)
- [ ] `./dependabot update --help` funciona, sin flag --version
