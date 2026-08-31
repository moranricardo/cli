package infra

type Collector struct{}

const CollectorImageName = "otel/opentelemetry-collector:latest"

func NewCollector() *Collector {
    return &Collector{}
}
