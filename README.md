# Dependabot CLI

Fork de [dependabot/cli](https://github.com/dependabot/cli) mantenido por @moranricardo.

## Instalación

```sh
go install [github.com/moranricardo/cli/cmd/dependabot@latest](https://github.com/moranricardo/cli/cmd/dependabot@latest)
# o
brew bundle
go build -o dependabot ./cmd/dependabot
```

## Uso

```sh
dependabot --help
dependabot update --help
```

## Desarrollo

```sh
go test ./... -count=1
golangci-lint run
yamllint .
```

## Licencia

MIT - Copyright 2022 GitHub, Inc.
