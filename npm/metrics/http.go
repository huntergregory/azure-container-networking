package metrics

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var started = false

const httpPort = ":8000"

// StartHTTP starts a HTTP endpoint on port 8000. Metrics are exposed on the endpoint /metrics.
// Set asGoRoutine to true if you want to be able to effectively run other code after calling this.
func StartHTTP(asGoRoutine bool) {
	if started {
		return
	}
	started = true

	http.Handle("/metrics", getHandler())
	if asGoRoutine {
		go http.ListenAndServe(httpPort, nil)
	} else {
		http.ListenAndServe(httpPort, nil)
	}
}

// getHandler returns the HTTP handler for the metrics endpoint
func getHandler() http.Handler {
	if handler == nil {
		handler = promhttp.Handler()
		// 	handler = promhttp.HandlerFor(networkingRegistry, promhttp.HandlerOpts{}) // promhttp.Handler()
	}
	return handler
}

func getMetricsText() (string, error) {
	response, err := http.Get("http://localhost" + httpPort + "/metrics")
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

// GetValue returns a gaugeMetric's value as shown in the HTML Prometheus endpoint.
func GetValue(gaugeMetric prometheus.Collector) (int, error) {
	return getMetricValue(allMetrics[gaugeMetric])
}

// GetCountValue returns the number of times a summaryMetric has recorded an observation as shown in the HTML Prometheus endpoint.
func GetCountValue(summaryMetric prometheus.Collector) (int, error) {
	return getMetricValue(allMetrics[summaryMetric] + "_count")
}

func getMetricValue(metricName string) (int, error) {
	if !started {
		StartHTTP(true)
		time.Sleep(2 * time.Second)
	}
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
