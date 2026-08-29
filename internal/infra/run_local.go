// Copyright (c) 2026 Ricardo Moran (moranricardo) - Local Runner
// Original implementation for docker-less environments (Termux, etc)
// Licensed under MIT same as github.com/dependabot/cli

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

// RunLocal executes update without Docker using updater-lite.sh
// Author: Ricardo Moran
func RunLocal(params RunParams) error {
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
	log.Printf("local runner: api=%s repo=%s pm=%s", apiUrl, params.LocalDir, params.Job.PackageManager)
	repoDir := params.LocalDir
	if repoDir == "" {
		repoDir = "."
	}
	absRepo, _ := filepath.Abs(repoDir)
	script := filepath.Join(absRepo, "updater-lite.sh")
	if _, err := os.Stat(script); os.IsNotExist(err) {
		script = "updater-lite.sh"
	}
	cmd := exec.CommandContext(ctx, "bash", script)
	cmd.Dir = absRepo
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("local updater failed: %w", err)
	}
	fmt.Println("local run complete")
	return nil
}
