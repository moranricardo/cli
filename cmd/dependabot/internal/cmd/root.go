package cmd

import (
  "fmt"
  "os"
  "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "dependabot",
  Short: "Dependabot Lite - Orquestador nativo ",
  Long: `Dependabot Lite es un orquestador nativo sin Docker`,
  SilenceUsage: true,
  SilenceErrors: true,
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
  }
}

func init() {
  rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose")
  rootCmd.AddCommand(newUpdateCmd())
}
