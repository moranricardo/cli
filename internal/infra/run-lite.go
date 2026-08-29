package infra

import (
  "context"
  "fmt"
  "log"
  "os"
  "os/exec"
  "path/filepath"
  "github.com/dependabot/cli/internal/server"
)

func RunLite(params RunParams) error {
  if err := params.Validate(); err != nil { return err }
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()
  api := server.NewAPI(params.Expected, params.Writer)
  defer api.Stop()
  apiUrl := params.ApiUrl
  if apiUrl == "" {
    apiUrl = fmt.Sprintf("http://127.0.0.1:%v", api.Port())
  }
  log.Printf("e6 lite - API %s repo %s pm %s", apiUrl, params.LocalDir, params.Job.PackageManager)
  repoDir := params.LocalDir
  if repoDir == "" { repoDir = "." }
  absRepo, _ := filepath.Abs(repoDir)
  cmd := exec.CommandContext(ctx, "bash", "updater-lite.sh")
  cmd.Dir = absRepo
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Env = append(os.Environ(), fmt.Sprintf("GITHUB_TOKEN=%s", os.Getenv("GITHUB_TOKEN")))
  if err := cmd.Run(); err != nil {
    log.Printf("updater-lite.sh error: %v", err)
  }
  fmt.Println("e6 lite complete")
  return nil
}
