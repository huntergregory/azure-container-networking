package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "npm"

// Gauge metrics have methods Inc(), Dec(), and Set(float64)
// Summary metrics has the method Observe(float64)
// For any Vector metric, you can call WithLabelValues(...string) before the above methods e.g. SomeGaugeVec.WithLabelValues("label1", "label2").Dec()

var (
	NumPolicies            = createGauge(numPoliciesLabel, numPoliciesHelp)
	AddPolicyExecTime      = createSummary(addPolicyExecTimeLabel, addPolicyExecTimeHelp)
	NumIPTableRules        = createGauge(numIPTableRulesLabel, numIPTableRulesHelp)
	AddIPTableRuleExecTime = createSummary(addIPTableRuleExecTimeLabel, addIPTableRuleExecTimeHelp)
	NumIPSets              = createGauge(numIPSetsLabel, numIPSetsHelp)
	AddIPSetExecTime       = createSummary(addIPSetExecTimeLabel, addIPSetExecTimeHelp)
	IPSetInventory         = createGaugeVec(ipsetInventoryLabel, ipsetInventoryHelp)
)

const (
	numPoliciesLabel = "num_policies"
	numPoliciesHelp  = "The number of current network policies for this node"

	addPolicyExecTimeLabel = "add_policy_exec_time"
	addPolicyExecTimeHelp  = "Execution time in milliseconds for adding a network policy"

	numIPTableRulesLabel = "num_iptables_rules"
	numIPTableRulesHelp  = "The number of current IPTable rules for this node"

	addIPTableRuleExecTimeLabel = "add_iptables_rule_exec_time"
	addIPTableRuleExecTimeHelp  = "Execution time in milliseconds for adding an IPTable rule to a chain"

	numIPSetsLabel = "num_ipsets"
	numIPSetsHelp  = "The number of current IP sets for this node"

	addIPSetExecTimeLabel = "add_ipset_exec_time"
	addIPSetExecTimeHelp  = "Execution time in milliseconds for creating an IP set"

	ipsetCountsLabel = "ipset_counts"
	ipsetCountsHelp  = "Number of entries in each individual IPSet"
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
