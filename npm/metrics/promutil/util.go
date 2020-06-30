package promutil

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

const delayAfterHTTPStart = 10

// NotifyIfErrors writes any non-nil errors to a testing utility
func NotifyIfErrors(t *testing.T, errors ...error) {
	allGood := true
	for _, err := range errors {
		if err != nil {
			allGood = false
			break
		}
	}
	if !allGood {
		t.Errorf("Encountered these errors while getting metric values: ")
		for _, err := range errors {
			if err != nil {
				t.Errorf("%v", err)
			}
		}
	}
}

// GetValue is used for validation. It returns a gaugeMetric's value.
func GetValue(gaugeMetric prometheus.Collector) (int, error) {
	dtoMetric, err := getDTOMetric(gaugeMetric)
	if err != nil {
		return 0, err
	}
	return int(dtoMetric.Gauge.GetValue()), nil
}

// GetCountValue is used for validation. It returns the number of times a summaryMetric has recorded an observation.
func GetCountValue(summaryMetric prometheus.Collector) (int, error) {
	dtoMetric, err := getDTOMetric(summaryMetric)
	if err != nil {
		return 0, err
	}
	return int(dtoMetric.Summary.GetSampleCount()), nil
}

func getDTOMetric(collector prometheus.Collector) (*dto.Metric, error) {
	channel := make(chan prometheus.Metric, 1)
	collector.Collect(channel)
	metric := &dto.Metric{}
	err := (<-channel).Write(metric)
	return metric, err
}
