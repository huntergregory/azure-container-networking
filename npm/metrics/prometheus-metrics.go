package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "npm"

var (
	NumPolicies            = createGauge(numPoliciesLabel, "The number of current network policies for this node")
	AddPolicyExecTime      = createSummary(addPolicyExecTimeLabel, "Execution time for adding a network policy")
	NumIPTableRules        = createGauge(numIPTableRules, "The number of current IPTable rules for this node")
	AddIPTableRuleExecTime = createSummary(addIPTableRuleExecTimeLabel, "Execution time for adding an IPTable rule to a chain")
	NumIPSets              = createGauge(numIPSetsLabel, "The number of current IP sets for this node")
	AddIPSetExecTime       = createSummary(addIPSetExecTimeLabel, "Execution time for creating an IP set")
)

const (
	numPoliciesLabel            = "num_policies"
	addPolicyExecTimeLabel      = "add_policy_exec_time"
	numIPTableRules             = "num_iptables_rules"
	addIPTableRuleExecTimeLabel = "add_iptables_rule_exec_time"
	numIPSetsLabel              = "num_ipsets"
	addIPSetExecTimeLabel       = "add_ipset_exec_time"
)

var allMetricNames = map[prometheus.Collector]string{
	NumPolicies:            numPoliciesLabel,
	AddPolicyExecTime:      addPolicyExecTimeLabel,
	NumIPTableRules:        numIPTableRules,
	AddIPTableRuleExecTime: addIPTableRuleExecTimeLabel,
	NumIPSets:              numIPSetsLabel,
	AddIPSetExecTime:       addIPSetExecTimeLabel,
}

func createGauge(name string, helpMessage string) prometheus.Gauge {
	return prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      helpMessage,
		},
		//[]string{"node"}, // include labels in a slice like this if creating Vectors
	)
}

func createSummary(name string, helpMessage string) prometheus.Summary {
	return prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       name,
			Help:       helpMessage,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, //quantiles
		},
	)
}

func init() {
	for metric := range allMetricNames {
		prometheus.DefaultRegisterer.MustRegister(metric)
	}
}

// Observe records a value in the given summary
func Observe(summary prometheus.Summary, value float64) {
	summary.Observe(value)
	// if changed to a vector, use summary.WithLabelValues(hostName).Observe(value)
}

// Inc increases a gauge by 1
func Inc(gauge prometheus.Gauge) {
	gauge.Inc()
}

// Dec decreases a gauge by 1
func Dec(gauge prometheus.Gauge) {
	gauge.Dec()
}

// GetMetricName is for validation purposes. It returns the name representation of any metric registered in this file.
// Returns an empty string if the metric is not declared and exported in this file.
func GetMetricName(collector prometheus.Collector) string {
	return allMetricNames[collector]
}
