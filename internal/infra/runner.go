package infra

import (
"context"
"fmt"

"github.com/moranricardo/cli/internal/model"
)

type RunParams struct {
LocalDir string
Job      *model.Job
}

func RunLite(ctx context.Context, p RunParams) error {
select {
case <-ctx.Done():
return ctx.Err()
default:
fmt.Printf("[Runner] Ejecutando ecosistema %s en %s\n", p.Job.PackageManager, p.LocalDir)
return nil
}
}
