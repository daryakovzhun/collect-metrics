package agent

type IAgent interface {
	GetCounterMetrics()
	GetGaugeMetrics()
}
