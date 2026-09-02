# Política de seguridad

## Versiones compatibles

Usamos Go semver. Solo damos soporte de seguridad a la rama principal y último release.

| Versión | Compatible |
| --- | --- |
| main (rama principal) | ✅ |
| >= 0.5.9 | ✅ |
| < 0.5.9 | ❌ |

## Cómo informar sobre una vulnerabilidad

**No abras un issue público.**

1. Usa **Private vulnerability reporting** de GitHub:
   https://github.com/moranricardo/cli/security/advisories/new

2. Incluye:
   - Descripción, impacto, PoC si tienes
   - Versión afectada: `moran-internal --version`
   - ¿Es explotable en `arm` / `arm64`?

Respondemos en **72h** y aplicamos divulgación coordinada de **90 días**.

## Alcance

**Dentro:** RCE, inyección, path traversal, leaks de tokens en `internal/*`, `cmd/dependabot`

**Fuera:** DoS de dependencias, reportes de `go mod` sin PoC, vulnerabilidades de Docker (este CLI es modo lite sin Docker)

Gracias por mantener el proyecto seguro.
