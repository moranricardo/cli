# Dependabot CLI - Copilot Instructions (forked by moranricardo)

## What is this?
Go CLI that orchestrates Dependabot update jobs via Docker. Does NOT resolve deps itself.
It runs 3 containers: Proxy + Updater + Fake API Server.
Flow: CLI -> proxy+updater on isolated networks -> updater calls fake API -> YAML output.

go.mod: github.com/dependabot/cli, Go 1.22+

## Layout
- cmd/dependabot/ - entrypoint, internal/cmd/ -> Cobra: update, test, graph, version
- internal/infra/ - Docker orchestration: run.go, updater.go, proxy.go, network.go, config.go, cadetails.go
- internal/model/ - types for jobs, creds, smoke tests, API payloads. Kebab-case YAML tags yaml:package-manager
- internal/server/ - fake API api.go + secure input input.go
- testdata/ - fixtures + scripts/*.txt (rsc.io/script)

## Critical Conventions
1. Model tags: kebab-case yaml + json. Add omitempty first for backward compat. See internal/model/job.go for lifecycle.
2. Command pattern: NewXCommand() *cobra.Command -> calls infra.Run(infra.RunParams{...}) -> init() registers. update and graph share extractInput()/processInput().
3. Credentials: $VAR expanded from env. LOCAL_GITHUB_ACCESS_TOKEN auto-injected. Never pass creds to updater directly - go via proxy. checkCredAccess() blocks write tokens.
4. Networking: 2 bridge networks per run: no-internet (internal:true, updater only) and internet (proxy only).
5. Ecosystems: Add in packageManagerLookup map in internal/infra/run.go (go_modules -> gomod).
6. Limited Env / ARM: Binary must work as linux-armv7. Avoid /tmp assumption, use ./. Flag -v requires arg (volume), not verbose. docker.sock may not exist. test -o ./test.yml must work without Docker.

## Build & Test
go build -o dependabot ./cmd/dependabot
./dependabot test -o ./test.yml
./dependabot update --help
go test ./... -race
go test ./cmd/dependabot/ -count=1
script/e2e <pattern>
