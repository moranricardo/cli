package infra
type Collector struct{ url string }
const CollectorImageName = "otel/opentelemetry-collector:latest"
func (c *Collector) TailLogs() error { return nil }
func (c *Collector) Close() error { return nil }
