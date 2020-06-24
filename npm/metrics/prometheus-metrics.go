package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// var networkingRegistry *prometheus.Registery
// var hostName = os.Getenv("HOSTNAME")

const tempHelp = "temporary help description" //TODO unique for each metric
const namespace = "npm"

// TODO add quantiles for summaries? remove quantiles?

var (
	NumPolicies = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "num_policies",
			// Help:      tempHelp,
		},
		//[]string{"node"},
		// include labels in a slice like above if a vector
	)

	AddPolicyExecTime = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "add_policy_exec_time",
			// Help:       tempHelp,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, // TODO remove?
		},
	)

	NumIpTableRules = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "num_iptables_rules",
			// Help:      tempHelp,
		},
	)

	AddIpTableRuleExecTime = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "add_iptables_rule_exec_time",
			// Help:      tempHelp,
		},
	)

	NumIpSets = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "num_ipsets",
			// Help:      tempHelp,
		},
	)

	AddIpSetExecTime = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "add_ipset_exec_time",
			// Help:      tempHelp,
		},
	)
)

var allMetrics = []prometheus.Collector{NumPolicies, AddPolicyExecTime, NumIpTableRules, AddIpTableRuleExecTime, NumIpSets, AddIpSetExecTime}
var handler http.Handler

func init() {
	// networkingRegistry = prometheus.NewRegistry()
	for _, metric := range allMetrics {
		err := prometheus.DefaultRegisterer.Register(metric)
		if err != nil {
			fmt.Printf("While registering a certain prometheus metric, an error occurred: %s", err)
		}
	}
}

func GetHandler() http.Handler {
	if handler == nil {
		handler = promhttp.Handler()
		// 	handler = promhttp.HandlerFor(networkingRegistry, promhttp.HandlerOpts{}) // promhttp.Handler()
	}
	return handler
}

func Observe(summary prometheus.Summary, value float64) {
	summary.Observe(value)
	// if changed to a vector, use summary.WithLabelValues(hostName).Observe(value)
}

func Inc(gauge prometheus.Gauge) {
	gauge.Inc()
}

func Dec(gauge prometheus.Gauge) {
	gauge.Dec()
}
