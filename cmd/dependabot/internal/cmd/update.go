package cmd

import (
 "context"
 "errors"
 "fmt"
 "os"
 "path/filepath"
 "time"
 "github.com/moranricardo/cli/internal/infra"
 "github.com/moranricardo/cli/internal/model"
 "github.com/spf13/cobra"
)

func detectEcosystem(dir string) string {
 if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil { return "go_modules" }
 if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil { return "npm_and_yarn" }
 if _, err := os.Stat(filepath.Join(dir, "requirements.txt")); err == nil { return "pip" }
 if _, err := os.Stat(filepath.Join(dir, "pyproject.toml")); err == nil { return "pip" }
 if _, err := os.Stat(filepath.Join(dir, "Cargo.toml")); err == nil { return "cargo" }
 return "go_modules"
}

func newUpdateCmd() *cobra.Command {
 cmd := &cobra.Command{
  Use: "update [ecosystem] [dir]",
  Short: "Run dependabot update locally (Go 1.27 STABLE)",
  Args: cobra.MaximumNArgs(2),
  RunE: func(cmd *cobra.Command, args []string) error {
   if os.Getenv("DEPENDABOT_LOCAL_RUN")!= "1" {
    return errors.New("activa DEPENDABOT_LOCAL_RUN=1")
   }
   ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Minute)
   defer cancel()
   ecosystem, dir := "", "."
   if len(args) == 1 {
    if info, err := os.Stat(args[0]); err == nil && info.IsDir() { dir = args[0] } else if args[0]!= "." { ecosystem = args[0] }
   }
   if len(args) == 2 { ecosystem = args[0]; dir = args[1] }
   absDir, _ := filepath.Abs(dir)
   if ecosystem == "" { ecosystem = detectEcosystem(absDir) }
   verbose, _ := cmd.Flags().GetBool("verbose")
   if verbose { fmt.Printf(">> Diamond STABLE [%s] dir=%s\n", ecosystem, dir) } // <-- relativo, no leak termux
   job := model.Job{PackageManager: ecosystem, Command: model.UpdateFilesCommand}
   return infra.RunLite(ctx, infra.RunParams{LocalDir: absDir, Job: &job})
  },
 }
 return cmd
}
