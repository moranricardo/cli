package infra

import (
"context"

docker "github.com/docker/docker/client"
)

type Collector struct {
url string
}

const CollectorImageName = "otel/opentelemetry-collector:latest"

func NewCollector(ctx context.Context, dockerClient *docker.Client, networks *Networks, runParams *RunParams, proxy *Proxy) (*Collector, error) {
return &Collector{url: ""}, nil
}

func (c *Collector) TailLogs(ctx context.Context, dockerClient *docker.Client) error {
return nil
}

func (c *Collector) Close() error {
return nil
}
