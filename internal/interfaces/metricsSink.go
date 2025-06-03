package interfaces

type MetricsSink interface {
	IncRequest()
	IncMalicious()
}
