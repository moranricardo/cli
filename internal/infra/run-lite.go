package infra

import (
"bytes"
"context"
"encoding/json"
"fmt"
"log"
"net/http"
"os"
"os/exec"
"path/filepath"
"strings"

"github.com/dependabot/cli/internal/server"
)

type UpdatedDependencyFile struct {
Name    string `json:"name"`
Content string `json:"content"`
}

type PRPayload struct {
Data struct {
BaseCommitSHA          string                  `json:"base-commit-sha"`
UpdatedDependencyFiles []UpdatedDependencyFile `json:"updated-dependency-files"`
} `json:"data"`
}

func RunLite(params RunParams) error {
if err := params.Validate(); err != nil {
return err
}

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

api := server.NewAPI(params.Expected, params.Writer)
defer api.Stop()

apiUrl := params.ApiUrl
if apiUrl == "" {
apiUrl = fmt.Sprintf("http://127.0.0.1:%v", api.Port())
}

log.Printf("local runner v0.5: api=%s repo=%s pm=%s", apiUrl, params.LocalDir, params.Job.PackageManager)

repoDir := params.LocalDir
if repoDir == "" {
repoDir = "."
}
absRepo, _ := filepath.Abs(repoDir)

selectivePkg := ""
for _, a := range os.Args {
if strings.Contains(a, "dependabot") {
continue
}
if a == "update" || a == "go_modules" || strings.HasPrefix(a, "-") {
continue
}
if a == "." || a == absRepo || a == repoDir {
continue
}
if a == "./..." {
selectivePkg = a
}
}
if v := os.Getenv("PACKAGE"); v != "" {
selectivePkg = v
}

workDir := absRepo
env := append(os.Environ(),
fmt.Sprintf("DEPENDABOT_API_URL=%s", apiUrl),
fmt.Sprintf("DEPENDABOT_PACKAGE_MANAGER=%s", params.Job.PackageManager),
fmt.Sprintf("DEPENDABOT_REPO_DIR=%s", absRepo),
)

switch params.Job.PackageManager {
case "go_modules":
if selectivePkg != "" {
fmt.Printf(">> [Native Go v0.5] Selective %s in %s\n", selectivePkg, workDir)
} else {
fmt.Printf(">> [Native Go v0.5] Full check in %s\n", workDir)
}

runCmd(ctx, workDir, env, "go", "list", "-m", "-u", "all")

if os.Getenv("DRY_RUN") != "1" {
if selectivePkg != "" && selectivePkg != "./..." {
runCmd(ctx, workDir, env, "go", "get", "-u", selectivePkg)
} else {
runCmd(ctx, workDir, env, "go", "get", "-u", "./...")
}
runCmd(ctx, workDir, env, "go", "mod", "tidy")
reportPR(apiUrl, workDir, absRepo)
}
default:
runCmd(ctx, workDir, env, "ls", "-la")
}

fmt.Println("local run complete v0.5")
return nil
}

func runCmd(ctx context.Context, dir string, env []string, name string, args ...string) {
cmd := exec.CommandContext(ctx, name, args...)
cmd.Dir = dir
cmd.Env = env
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
if err := cmd.Run(); err != nil {
fmt.Printf(">> [warn] %s %s failed: %v\n", name, strings.Join(args, " "), err)
}
}

func reportPR(apiUrl, workDir, repoDir string) {
goModPath := filepath.Join(workDir, "go.mod")
goMod, err := os.ReadFile(goModPath)
if err != nil {
fmt.Printf(">> [warn] could not read go.mod: %v\n", err)
return
}

files := []UpdatedDependencyFile{
{Name: "go.mod", Content: string(goMod)},
}

goSumPath := filepath.Join(workDir, "go.sum")
if goSum, err := os.ReadFile(goSumPath); err == nil {
files = append(files, UpdatedDependencyFile{Name: "go.sum", Content: string(goSum)})
}

commitSHA := "main"
cmdSHA := exec.Command("git", "rev-parse", "HEAD")
cmdSHA.Dir = repoDir
if out, err := cmdSHA.Output(); err == nil {
commitSHA = string(bytes.TrimSpace(out))
}

payload := PRPayload{}
payload.Data.BaseCommitSHA = commitSHA
payload.Data.UpdatedDependencyFiles = files

body, err := json.Marshal(payload)
if err != nil {
fmt.Printf(">> [warn] failed to marshal payload: %v\n", err)
return
}

fmt.Printf(">> Reporting PR to API Proxy (%s)...\n", apiUrl)
resp, err := http.Post(apiUrl+"/create_pull_request", "application/json", bytes.NewBuffer(body))
if err != nil {
fmt.Printf(">> Warning: PR failed: %v\n", err)
return
}
defer resp.Body.Close()
fmt.Printf(">> PR status: %s\n", resp.Status)
}
