package infra

import (
 "time"
 "github.com/moranricardo/cli/internal/model"
 "io"
 "context"
)

type RunParams struct {
 Input         string
 Job           *model.Job
 Expected      []model.Output
 LocalDir      string
 Creds         []model.Credential
 CacheDir      string
 Output        string
 ProxyCertPath string
 PullImages    bool
 Debug         bool
 Flamegraph    bool
 Volumes       []string
 Timeout       time.Duration
 ExtraHosts    []string
 UpdaterImage  string
 ApiUrl        string
 Writer        io.Writer
}

func (r RunParams) Validate() error { return nil }
type Networks struct{}
type Proxy struct{ url string }

func NewProxy(ctx context.Context, params *RunParams, nets *Networks) (*Proxy, error) {
 return &Proxy{url: "http://127.0.0.1:44569"}, nil
}
func (p *Proxy) TailLogs(ctx context.Context) error { return nil }
func (p *Proxy) Close() error { return nil }
