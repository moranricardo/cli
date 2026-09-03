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

func newUpdateCmd() *cobra.Command {
cmd := &cobra.Command{
Use:   "update [ecosystem] [dir]",
Short: "Run dependabot update locally without Docker (Go STABLE)",
Long:  `Diamond STABLE: Ejecuta update nativo sin Docker. Detecta go.mod automáticamente.`,
Args:  cobra.MaximumNArgs(2),
RunE: func(cmd *cobra.Command, args []string) error {
if os.Getenv("DEPENDABOT_LOCAL_RUN") != "1" {
return errors.New("activa DEPENDABOT_LOCAL_RUN=1 (ej: DEPENDABOT_LOCAL_RUN=1 dependabot update go_modules .)")
}

ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Minute)
defer cancel()

ecosystem := "go_modules"
dir := "."

if len(args) >= 1 && args[0] != "." {
ecosystem = args[0]
}
if len(args) == 2 {
dir = args[1]
}

absDir, err := filepath.Abs(dir)
if err != nil {
return fmt.Errorf("dir inválida %q: %w", dir, err)
}

if info, err := os.Stat(absDir); err != nil || !info.IsDir() {
return fmt.Errorf("no existe dir %s", absDir)
}

if ecosystem == "go_modules" {
if _, err := os.Stat(filepath.Join(absDir, "go.mod")); err != nil {
return fmt.Errorf("no se encontró go.mod en %s: %w", absDir, err)
}
}

verbose, _ := cmd.Flags().GetBool("verbose")
if verbose {
fmt.Printf(">> Diamond STABLE\n>> Ecosystem: %s\n>> Dir: %s\n", ecosystem, absDir)
}

job := model.Job{
PackageManager: ecosystem,
Command:        model.UpdateFilesCommand,
}

params := infra.RunParams{
LocalDir: absDir,
Job:      &job,
}

return infra.RunLite(ctx, params)
},
}

return cmd
}
