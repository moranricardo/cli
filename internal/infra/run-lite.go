package infra

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type RunParams struct {
	LocalDir     string
	ExcludePaths []string
	Job          interface{}
}

func isExcluded(path string, excludes []string) bool {
	base := filepath.Base(path)
	for _, exclude := range excludes {
		if matched, _ := filepath.Match(exclude, base); matched {
			return true
		}
	}
	return false
}

func RunLite(ctx context.Context, p RunParams) error {
	workflowsDir := ".github/workflows"
	if p.LocalDir != "" {
		workflowsDir = filepath.Join(p.LocalDir, workflowsDir)
	}

	if _, err := os.Stat(workflowsDir); os.IsNotExist(err) {
		return nil
	}

	excludes := p.ExcludePaths
	if len(excludes) == 0 {
		excludes = []string{"prisma-visor.yml"}
	}

	fmt.Printf("[Runner STABLE] Ecosistema github_actions en %s\n", workflowsDir)

	count := 0
	err := filepath.Walk(workflowsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yml" && filepath.Ext(path) != ".yaml" {
			return nil
		}

		if isExcluded(path, excludes) {
			fmt.Printf(">> Skip excluido: %s\n", path)
			return nil
		}

		fmt.Printf(">> Analizando workflow: %s\n", path)
		count++
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf(">> [github_actions] %d workflows analizados\n", count)
	return nil
}
