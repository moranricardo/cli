package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dependabot",
	Short: "Dependabot Lite v0.1-lite - without containers",
	Long:  "Orquestador nativo de Dependabot sin dependencias externas.",
	Example: `  $ dependabot update go_modules --local .
  $ dependabot version`,
	Version: Version(),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
