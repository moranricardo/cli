package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dependabot/cli/internal/infra"
	"github.com/dependabot/cli/internal/model"
	"github.com/spf13/cobra"
)

var localRun bool

var updateCmd = &cobra.Command{
	Use:   "update [ecosystem]",
	Short: "Run dependabot update locally without Docker",
	Long: `Ejecuta el update de forma nativa:
- Detecta go.mod, package.json, etc.
- No necesita Docker ni test.yml`,
	Example: `  $ dependabot update go_modules --local-run .
  $ dependabot update go_modules --local-run -v .`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !localRun {
			return fmt.Errorf("usa --local-run (ej: dependabot update go_modules --local-run .)")
		}
		verbose, _ := cmd.Root().PersistentFlags().GetBool("verbose")
		return RunLite(args, verbose)
	},
}

func init() {
	updateCmd.Flags().BoolVar(&localRun, "local-run", false, "run without docker")
	rootCmd.AddCommand(updateCmd)
}

func RunLite(args []string, verbose bool) error {
	ecosystem := "go_modules"
	if len(args) > 0 && args[0] != "." {
		ecosystem = args[0]
	}

	cwd, _ := os.Getwd()
	if len(args) > 0 && args[len(args)-1] != ecosystem {
		if filepath.IsAbs(args[len(args)-1]) || args[len(args)-1] == "." {
			cwd, _ = filepath.Abs(args[len(args)-1])
		}
	}

	if verbose {
		fmt.Printf(">> Ecosystem: %s\n", ecosystem)
		fmt.Printf(">> Dir: %s\n", cwd)
	}

	// Validación nativa previa
	if ecosystem == "go_modules" {
		if _, err := os.Stat(filepath.Join(cwd, "go.mod")); os.IsNotExist(err) {
			return fmt.Errorf("no se encontró go.mod en %s", cwd)
		}
		if verbose {
			fmt.Println(">> go.mod encontrado, resolviendo nativo...")
		}
	}

	job := model.Job{
		PackageManager: ecosystem,
		Command:        model.UpdateFilesCommand,
	}

	params := infra.RunParams{
		LocalDir: cwd,
		Job:      &job,
	}

	if verbose {
		fmt.Println(">> Dispatch a infra.RunLite (modo lite nativo)")
	}

	return infra.RunLite(params)
}
