package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "npm"

var (
	NumPolicies            = createGauge(numPoliciesLabel, "The number of current network policies for this node")
	AddPolicyExecTime      = createSummary(addPolicyExecTimeLabel, "Execution time in milliseconds for adding a network policy")
	NumIPTableRules        = createGauge(numIPTableRules, "The number of current IPTable rules for this node")
	AddIPTableRuleExecTime = createSummary(addIPTableRuleExecTimeLabel, "Execution time in milliseconds for adding an IPTable rule to a chain")
	NumIPSets              = createGauge(numIPSetsLabel, "The number of current IP sets for this node")
	AddIPSetExecTime       = createSummary(addIPSetExecTimeLabel, "Execution time in milliseconds for creating an IP set")
)

const (
	numPoliciesLabel            = "num_policies"
	addPolicyExecTimeLabel      = "add_policy_exec_time"
	numIPTableRules             = "num_iptables_rules"
	addIPTableRuleExecTimeLabel = "add_iptables_rule_exec_time"
	numIPSetsLabel              = "num_ipsets"
	addIPSetExecTimeLabel       = "add_ipset_exec_time"
)

func createGauge(name string, helpMessage string) prometheus.Gauge {
	gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      helpMessage,
		},
		//[]string{"node"}, // include labels in a slice like this if creating Vectors
	)
	prometheus.DefaultRegisterer.MustRegister(gauge)
	return gauge
}

func createSummary(name string, helpMessage string) prometheus.Summary {
	summary := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       name,
			Help:       helpMessage,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, //quantiles
		},
	)
	prometheus.DefaultRegisterer.MustRegister(summary)
	return summary
}
