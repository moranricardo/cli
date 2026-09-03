package infra

import (
"context"
"fmt"
)

func RunLite(ctx context.Context, p RunParams) error {
select {
case <-ctx.Done():
return ctx.Err()
default:
if p.Job != nil {
fmt.Printf("[Runner STABLE] Ejecutando ecosistema %s en %s\n", p.Job.PackageManager, ".")
} else {
fmt.Printf("[Runner STABLE] Ejecutando en %s\n", ".")
}
return nil
}
}
