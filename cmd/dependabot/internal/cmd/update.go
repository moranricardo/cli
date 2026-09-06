# deja: if verbose { fmt.Printf(">> 
# Diamond STABLE [%s] dir=%s\n", 
# ecosystem, dir) }package cmd

import ( "context" "errors" "fmt" "os" 
 "path/filepath" "time" 
 "github.com/moranricardo/cli/internal/infra" 
 "github.com/moranricardo/cli/internal/model" 
 "github.com/spf13/cobra"
# 2. revisa copilot-instructions)
nano .github/copilot-instructions.md func 
detectEcosystem(dir string) string {
# borra cualquier línea que diga "termux" 
# if _, err := os.Stat(filepath.Join(dir, 
# "go.mod")); err == nil { return 
# "go_modules" }
 if _, err := os.Stat(filepath.Join(dir, 
 "package.json")); err == nil { return 
 "npm_and_yarn" } if _, err := 
 os.Stat(filepath.Join(dir, 
 "requirements.txt")); err == nil { return 
 "pip" } if _, err := 
 os.Stat(filepath.Join(dir, 
 "pyproject.toml")); err == nil { return 
 "pip" } if _, err := 
 os.Stat(filepath.Join(dir, 
 "Cargo.toml")); err == nil { return 
 "cargo" } return "go_modules"
# 3. verifica limpio}
grep -R "termux\|/data/data" 
--exclude-dir=.git -n || echo "LIMPIO ✅" 
func newUpdateCmd() *cobra.Command {
 cmd := &cobra.Command{ Use: "update 
  [ecosystem] [dir]", Short: "Run 
  dependabot update locally ", Args: 
  cobra.MaximumNArgs(2), RunE: func(cmd 
  *cobra.Command, args []string) error {
# 4. si da LIMPIO, push if 
# os.Getenv("DEPENDABOT_LOCAL_RUN")!= "1" 
# {
git add -A return errors.New("activa 
DEPENDABOT_LOCAL_RUN=1") git commit -m 
"chore: purga final refs termux en codigo 
y docs" }
git push   ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Minute)
   defer cancel()
   ecosystem, dir := "", "."
   if len(args) == 1 {
    if info, err := os.Stat(args[0]); err == nil && info.IsDir() { dir = args[0] } else if args[0]!= "." { ecosystem = args[0] }
   }
   if len(args) == 2 { ecosystem = args[0]; dir = args[1] }
   absDir, _ := filepath.Abs(dir)
   if ecosystem == "" { ecosystem = detectEcosystem(absDir) }
   verbose, _ := cmd.Flags().GetBool("verbose")
// relativo
   job := model.Job{PackageManager: ecosystem, Command: model.UpdateFilesCommand}
   return infra.RunLite(ctx, infra.RunParams{LocalDir: absDir, Job: &job})
  },
 }
 return cmd
}
