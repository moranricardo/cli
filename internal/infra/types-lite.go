package infra

import (
"github.com/moranricardo/cli/internal/model"
)

type RunParams struct {
LocalDir string
Job      *model.Job
}
