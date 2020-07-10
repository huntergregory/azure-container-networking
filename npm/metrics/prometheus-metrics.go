package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "npm"

// Prometheus Metrics
// Gauge metrics have the methods Inc(), Dec(), and Set(float64)
// Summary metrics has the method Observe(float64)
// For any Vector metric, you can call With(prometheus.Labels) before the above methods
//   e.g. SomeGaugeVec.With(prometheus.Labels{label1: val1, label2: val2, ...).Dec()
var (
	NumPolicies            = createGauge(numPoliciesName, numPoliciesHelp)
	AddPolicyExecTime      = createSummary(addPolicyExecTimeName, addPolicyExecTimeHelp)
	NumIPTableRules        = createGauge(numIPTableRulesName, numIPTableRulesHelp)
	AddIPTableRuleExecTime = createSummary(addIPTableRuleExecTimeName, addIPTableRuleExecTimeHelp)
	NumIPSets              = createGauge(numIPSetsName, numIPSetsHelp)
	AddIPSetExecTime       = createSummary(addIPSetExecTimeName, addIPSetExecTimeHelp)
	IPSetInventory         = createGaugeVec(ipsetInventoryName, ipsetInventoryHelp, SetNameLabel)
)

// Constants for metric names and descriptions as well as exported labels for Vector metrics
const (
	numPoliciesName = "num_policies"
	numPoliciesHelp = "The number of current network policies for this node"

	addPolicyExecTimeName = "add_policy_exec_time"
	addPolicyExecTimeHelp = "Execution time in milliseconds for adding a network policy"

	numIPTableRulesName = "num_iptables_rules"
	numIPTableRulesHelp = "The number of current IPTable rules for this node"

	addIPTableRuleExecTimeName = "add_iptables_rule_exec_time"
	addIPTableRuleExecTimeHelp = "Execution time in milliseconds for adding an IPTable rule to a chain"

	numIPSetsName = "num_ipsets"
	numIPSetsHelp = "The number of current IP sets for this node"

	addIPSetExecTimeName = "add_ipset_exec_time"
	addIPSetExecTimeHelp = "Execution time in milliseconds for creating an IP set"

	ipsetInventoryName = "ipset_counts"
	ipsetInventoryHelp = "Number of entries in each individual IPSet"
	SetNameLabel       = "set_name"
)

func createGauge(name string, helpMessage string) prometheus.Gauge {
	gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      helpMessage,
		},
	)
	prometheus.DefaultRegisterer.MustRegister(gauge)
	return gauge
}

func createGaugeVec(name string, helpMessage string, labels ...string) *prometheus.GaugeVec {
	gaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      helpMessage,
		},
		labels,
	)
	prometheus.DefaultRegisterer.MustRegister(gaugeVec)
	return gaugeVec
}

func createSummary(name string, helpMessage string) prometheus.Summary {
	summary := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace:  namespace,
			Name:       name,
			Help:       helpMessage,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			// quantiles e.g. the "0.5 quantile" will actually be the phi quantile for some phi in [0.5 - 0.05, 0.5 + 0.05]
		},
	)
	prometheus.DefaultRegisterer.MustRegister(summary)
	return summary
}
