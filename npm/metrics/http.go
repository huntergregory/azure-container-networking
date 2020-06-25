package metrics

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var started = false

func StartHTTP(asGoRoutine bool) {
	if started {
		return
	}
	started = true

	http.Handle("/metrics", GetHandler())
	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi!\n")
	})
	if asGoRoutine {
		go http.ListenAndServe(":8000", nil)
	} else {
		http.ListenAndServe(":8000", nil)
	}
}

func getMetricsText() (string, error) {
	response, err := http.Get("http://localhost:8000/metrics")
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

func GetValue(gaugeMetric prometheus.Collector) (int, error) {
	return getMetricValue(allMetrics[gaugeMetric])
}

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
