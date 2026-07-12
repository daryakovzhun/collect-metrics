package handler

type MetricsList struct {
	Items []MetricView
}

type MetricView struct {
	Name  string
	Type  string
	Value string
}
