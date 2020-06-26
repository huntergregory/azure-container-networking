package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// HTTPPort is the port used by the HTTP server (includes a preceding colon)
	HTTPPort = ":8000"

	//MetricsPath is the path for the Prometheus metrics endpoint (includes preceding slash)
	MetricsPath = "/metrics"
)

var started = false
var handler http.Handler

// StartHTTP starts a HTTP server with endpoint on port 8000. Metrics are exposed on the endpoint /metrics.
// Set asGoRoutine to true if you want to be able to effectively run other code after calling this.
// The function will pause for delayAmountAfterStart seconds after starting the HTTP server for the first time.
func StartHTTP(asGoRoutine bool, delayAmountAfterStart int) {
	if started {
		return
	}
	started = true

	http.Handle(MetricsPath, getHandler())
	if asGoRoutine {
		go http.ListenAndServe(HTTPPort, nil)
	} else {
		http.ListenAndServe(HTTPPort, nil)
	}
	time.Sleep(time.Second * time.Duration(delayAmountAfterStart))
}

// getHandler returns the HTTP handler for the metrics endpoint
func getHandler() http.Handler {
	if handler == nil {
		handler = promhttp.Handler()
		// 	handler = promhttp.HandlerFor(networkingRegistry, promhttp.HandlerOpts{}) // promhttp.Handler()
	}
	return handler
}
