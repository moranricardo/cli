<h1 align="center">
    <picture>
        <source media="(prefers-color-scheme: light)" srcset="https://user-images.githubusercontent.com/7659/174594540-5e29e523-396a-465b-9a6e-6cab5b15a568.svg">
        <source media="(prefers-color-scheme: dark)" srcset="https://user-images.githubusercontent.com/7659/174594559-0b3ddaa7-e75b-4f10-9dee-b51431a9fd4c.svg">
        <img src="https://user-images.githubusercontent.com/7659/174594540-5e29e523-396a-465b-9a6e-6cab5b15a568.svg" alt="Dependabot" width="336">
    </picture>
</h1>

The `dependabot` CLI is a tool for running Dependabot update jobs.

## Installation

Use any of the following for a pain-free installation:

* If you have [`go`](https://go.dev/doc/install) installed, you can run:
   ```shell
   go install github.com/dependabot/cli/cmd/dependabot@latest
   ```
   The benefit of this method is that re-running the command will always update to the latest version.
* You can download a pre-built binary from the [releases] page.
* On Mac or Linux, you can run `brew install dependabot`

## Requirements

* [Docker]

## Contributing

Check out our [contributing guidelines][contributing] for instructions on
building the project locally, sharing feedback, and submitting pull requests.

## Usage

```console
$ dependabot
Run Dependabot jobs from the command line.

Usage:
  dependabot [command]

Examples:
  $ dependabot update go_modules dependabot/cli
  $ dependabot test -f input.yml

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  test        Run a smoke test
  update      Perform an update job

