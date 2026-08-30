package cmd

import (
	"fmt"
	"os"

	"github.com/dependabot/cli/internal/infra"
	"github.com/dependabot/cli/internal/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var localRun bool

type JobInput struct {
	Input struct {
		Job model.Job `yaml:"job"`
	} `yaml:"input"`
}

var updateCmd = &cobra.Command{
	Use:   "update [ecosystem]",
	Short: "Run dependabot update locally without Docker",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !localRun {
			return fmt.Errorf("use --local-run flag")
		}
		return RunLite(args)
	},
}

func init() {
	updateCmd.Flags().BoolVar(&localRun, "local-run", false, "run without docker")
	rootCmd.AddCommand(updateCmd)
}

func RunLite(args []string) error {
	ecosystem := "generic"
	if len(args) > 0 {
		ecosystem = args[0]
	}

	cwd, _ := os.Getwd()
	var job model.Job
	data, err := os.ReadFile("test.yml")
	if err == nil {
		var inputJob JobInput
		if err := yaml.Unmarshal(data, &inputJob); err == nil {
			job = inputJob.Input.Job
		}
	}

	if job.PackageManager == "" {
		job.PackageManager = ecosystem
	}
	if job.Command == "" {
		job.Command = model.UpdateFilesCommand
	}

	params := infra.RunParams{
		LocalDir: cwd,
		Input:    "test.yml",
		Job:      &job,
	}

	return infra.RunLite(params)
}

func runLocal(args []string) error {
	return RunLite(args)
}
