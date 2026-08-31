package cmd

import (
"fmt"
"os"

"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
Use:   "dependabot",
Short: "Dependabot Lite - Orquestador nativo de actualizaciones de código sin contenedores",
Long: `Dependabot Lite es un orquestador nativo de actualizaciones
sin dependencias de contenedores ni scripts externos.

Flujo de trabajo:
  1. Inspecciona el entorno y detecta el gestor de paquetes.
  2. Resuelve dependencias de forma nativa (go_modules, npm, pip, etc.).
  3. Despacha el payload con las modificaciones al API Proxy.`,
Example: `  $ dependabot update go_modules --local-run .
  $ dependabot update npm_and_yarn --local-run
  $ dependabot version`,
Version:       Version(),
SilenceUsage:  true,
SilenceErrors: true,
}

func Execute() {
if err := rootCmd.Execute(); err != nil {
fmt.Fprintf(os.Stderr, "Error: %v\n", err)
os.Exit(1)
}
}

func init() {
rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Activar salida detallada")
}
