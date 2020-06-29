package promutil

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Azure/azure-container-networking/npm/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const delayAfterHTTPStart = 10

// GetValue is used for validation. It returns a gaugeMetric's value as shown in the HTML Prometheus endpoint.
func GetValue(gaugeMetric prometheus.Collector) (int, error) {
	return getMetricValue(metrics.GetMetricName(gaugeMetric))
}

// GetCountValue is used for validation. It returns the number of times a summaryMetric has recorded an observation as shown in the HTML Prometheus endpoint.
func GetCountValue(summaryMetric prometheus.Collector) (int, error) {
	return getMetricValue(metrics.GetMetricName(summaryMetric) + "_count")
}

func getMetricValue(metricName string) (int, error) {
	metrics.StartHTTP(true, delayAfterHTTPStart)
	regex := regexp.MustCompile(metricName + " [0-9]+")
	if regex == nil {
		return 0, fmt.Errorf("Couldn't compile regular expression for metric: " + metricName)
	}
	text, err := getMetricsText()
	if err != nil {
		return 0, err
	}
	locations := regex.FindStringIndex(text)
	if locations == nil {
		return 0, fmt.Errorf("Couldn't find a match for metric: " + metricName)
	}
	start := locations[0]
	end := locations[1]
	value := text[start+len(metricName)+1 : end]
	return strconv.Atoi(value)
}

func getMetricsText() (string, error) {
	response, err := http.Get("http://localhost" + metrics.HTTPPort + metrics.MetricsPath)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
