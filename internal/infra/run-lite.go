
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
	log.Printf("local runner v0.4: api=%s repo=%s pm=%s", apiUrl, params.LocalDir, params.Job.PackageManager)

	repoDir := params.LocalDir
	if repoDir == "" {
		repoDir = "."
	}
	absRepo, _ := filepath.Abs(repoDir)

	// v0.4: detectar paquete selectivo desde os.Args
	selectivePkg := ""
	for _, a := range os.Args {
		if a == "go_modules" || a == "npm" || strings.HasPrefix(a, "-") || a == "update" || a == "dependabot" {
			continue
		}
		if strings.Contains(a, "/") || strings.Contains(a, ".") {
			if a != absRepo && a != repoDir && a != "." {
				selectivePkg = a
			}
		}
	}
	if envPkg := os.Getenv("PACKAGE"); envPkg != "" {
		selectivePkg = envPkg
	}

	env := append(os.Environ(),
		fmt.Sprintf("DEPENDABOT_API_URL=%s", apiUrl),
		fmt.Sprintf("DEPENDABOT_PACKAGE_MANAGER=%s", params.Job.PackageManager),
		fmt.Sprintf("DEPENDABOT_REPO_DIR=%s", absRepo),
	)

	switch params.Job.PackageManager {
	case "go_modules":
		if selectivePkg != "" {
			fmt.Printf(">> [Native Go v0.4] Selective update: %s in %s\n", selectivePkg, absRepo)
		} else {
			fmt.Printf(">> [Native Go v0.4] Checking updates in %s\n", absRepo)
		}
		runCmd(ctx, absRepo, env, "go", "list", "-m", "-u", "all")

		if os.Getenv("DRY_RUN") != "1" {
			if selectivePkg != "" {
				fmt.Printf(">> Updating only %s...\n", selectivePkg)
				runCmd(ctx, absRepo, env, "go", "get", "-u", selectivePkg)
			} else {
				fmt.Println(">> Updating all modules (go get -u ./...)...")
				runCmd(ctx, absRepo, env, "go", "get", "-u", "./...")
			}
			runCmd(ctx, absRepo, env, "go", "mod", "tidy")
			reportPR(apiUrl, absRepo)
		} else {
			fmt.Println(">> DRY_RUN=1 - skipping update")
		}
	default:
		fmt.Printf(">> [Fallback] No runner for %s\n", params.Job.PackageManager)
		runCmd(ctx, absRepo, env, "ls", "-la")
	}
	fmt.Println("local run complete v0.4")
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

func reportPR(apiUrl, repoDir string) {
	goMod, errMod := os.ReadFile(filepath.Join(repoDir, "go.mod"))
	goSum, _ := os.ReadFile(filepath.Join(repoDir, "go.sum"))
	if errMod != nil {
		return
	}
	commitSHA := "main"
	cmdSHA := exec.Command("git", "rev-parse", "HEAD")
	cmdSHA.Dir = repoDir
	if out, err := cmdSHA.Output(); err == nil {
		commitSHA = string(bytes.TrimSpace(out))
	}
	payload := PRPayload{}
	payload.Data.BaseCommitSHA = commitSHA
	payload.Data.UpdatedDependencyFiles = []UpdatedDependencyFile{
		{Name: "go.mod", Content: string(goMod)},
		{Name: "go.sum", Content: string(goSum)},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return
	}
	fmt.Printf(">> Reporting PR to API Proxy (%s)...\n", apiUrl)
	resp, err := http.Post(apiUrl+"/create_pull_request", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf(">> Warning: PR reporting failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf(">> PR status: %s\n", resp.Status)
}

