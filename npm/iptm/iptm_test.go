package iptm

import (
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-container-networking/npm/metrics"

	"github.com/Azure/azure-container-networking/npm/util"
)

const testPrometheusToo = true
const prometheusErrorMessage = "You can turn off Prometheus testing by flipping the boolean constant testPrometheusToo."

func printPrometheusError(t *testing.T, message string) {
	t.Errorf(message + ". " + prometheusErrorMessage)
}

func TestSave(t *testing.T) {
	iptMgr := &IptablesManager{}
	if err := iptMgr.Save(util.IptablesTestConfigFile); err != nil {
		t.Errorf("TestSave failed @ iptMgr.Save")
	}
}

func TestRestore(t *testing.T) {
	iptMgr := &IptablesManager{}
	if err := iptMgr.Save(util.IptablesTestConfigFile); err != nil {
		t.Errorf("TestRestore failed @ iptMgr.Save")
	}

	if err := iptMgr.Restore(util.IptablesTestConfigFile); err != nil {
		t.Errorf("TestRestore failed @ iptMgr.Restore")
	}
}

func TestInitNpmChains(t *testing.T) {
	iptMgr := &IptablesManager{}

	if err := iptMgr.Save(util.IptablesTestConfigFile); err != nil {
		t.Errorf("TestInitNpmChains failed @ iptMgr.Save")
	}

	defer func() {
		if err := iptMgr.Restore(util.IptablesTestConfigFile); err != nil {
			t.Errorf("TestInitNpmChains failed @ iptMgr.Restore")
		}
	}()

	if err := iptMgr.InitNpmChains(); err != nil {
		t.Errorf("TestInitNpmChains @ iptMgr.InitNpmChains")
	}
}

func TestUninitNpmChains(t *testing.T) {
	iptMgr := &IptablesManager{}

	if err := iptMgr.Save(util.IptablesTestConfigFile); err != nil {
		t.Errorf("TestUninitNpmChains failed @ iptMgr.Save")
	}

	defer func() {
		if err := iptMgr.Restore(util.IptablesTestConfigFile); err != nil {
			t.Errorf("TestUninitNpmChains failed @ iptMgr.Restore")
		}
	}()

	if err := iptMgr.InitNpmChains(); err != nil {
		t.Errorf("TestUninitNpmChains @ iptMgr.InitNpmChains")
	}

	if err := iptMgr.UninitNpmChains(); err != nil {
		t.Errorf("TestUninitNpmChains @ iptMgr.UninitNpmChains")
	}
}

func TestExists(t *testing.T) {
	iptMgr := &IptablesManager{}
	if err := iptMgr.Save(util.IptablesTestConfigFile); err != nil {
		t.Errorf("TestExists failed @ iptMgr.Save")
	}

	defer func() {
		if err := iptMgr.Restore(util.IptablesTestConfigFile); err != nil {
			t.Errorf("TestExists failed @ iptMgr.Restore")
		}
	}()

	iptMgr.OperationFlag = util.IptablesCheckFlag
	entry := &IptEntry{
		Chain: util.IptablesForwardChain,
		Specs: []string{
			util.IptablesJumpFlag,
			util.IptablesAccept,
		},
	}
	if _, err := iptMgr.Exists(entry); err != nil {
		t.Errorf("TestExists failed @ iptMgr.Exists")
	}
}

func TestAddChain(t *testing.T) {
	iptMgr := &IptablesManager{}
	if err := iptMgr.Save(util.IptablesTestConfigFile); err != nil {
		t.Errorf("TestAddChain failed @ iptMgr.Save")
	}

	defer func() {
		if err := iptMgr.Restore(util.IptablesTestConfigFile); err != nil {
			t.Errorf("TestAddChain failed @ iptMgr.Restore")
		}
	}()

	if err := iptMgr.AddChain("TEST-CHAIN"); err != nil {
		t.Errorf("TestAddChain failed @ iptMgr.AddChain")
	}
}

func TestDeleteChain(t *testing.T) {
	iptMgr := &IptablesManager{}
	if err := iptMgr.Save(util.IptablesTestConfigFile); err != nil {
		t.Errorf("TestDeleteChain failed @ iptMgr.Save")
	}

	defer func() {
		if err := iptMgr.Restore(util.IptablesTestConfigFile); err != nil {
			t.Errorf("TestDeleteChain failed @ iptMgr.Restore")
		}
	}()

	if err := iptMgr.AddChain("TEST-CHAIN"); err != nil {
		t.Errorf("TestDeleteChain failed @ iptMgr.AddChain")
	}

	if err := iptMgr.DeleteChain("TEST-CHAIN"); err != nil {
		t.Errorf("TestDeleteChain failed @ iptMgr.DeleteChain")
	}
}

func TestAdd(t *testing.T) {
	iptMgr := &IptablesManager{}
	if err := iptMgr.Save(util.IptablesTestConfigFile); err != nil {
		t.Errorf("TestAdd failed @ iptMgr.Save")
	}

	defer func() {
		if err := iptMgr.Restore(util.IptablesTestConfigFile); err != nil {
			t.Errorf("TestAdd failed @ iptMgr.Restore")
		}
	}()

	entry := &IptEntry{
		Chain: util.IptablesForwardChain,
		Specs: []string{
			util.IptablesJumpFlag,
			util.IptablesReject,
		},
	}

	var (
		val    = 0
		newVal = 0
		err    error
	)
	if testPrometheusToo {
		val, err = metrics.GetValue("num_iptables_rules")
		if err != nil {
			printPrometheusError(t, "Problem getting http metrics")
		}
	}

	if err := iptMgr.Add(entry); err != nil {
		t.Errorf("TestAdd failed @ iptMgr.Add")
	}

	if testPrometheusToo {
		newVal, err = metrics.GetValue("num_iptables_rules")
		if err != nil {
			printPrometheusError(t, "Problem getting http metrics")
		}
		if newVal != val+1 {
			printPrometheusError(t, "Add iptable rule didn't register in prometheus")
		}
	}
}

func TestDelete(t *testing.T) {
	iptMgr := &IptablesManager{}
	if err := iptMgr.Save(util.IptablesTestConfigFile); err != nil {
		t.Errorf("TestDelete failed @ iptMgr.Save")
	}

	defer func() {
		if err := iptMgr.Restore(util.IptablesTestConfigFile); err != nil {
			t.Errorf("TestDelete failed @ iptMgr.Restore")
		}
	}()

	entry := &IptEntry{
		Chain: util.IptablesForwardChain,
		Specs: []string{
			util.IptablesJumpFlag,
			util.IptablesReject,
		},
	}
	if err := iptMgr.Add(entry); err != nil {
		t.Errorf("TestDelete failed @ iptMgr.Add")
	}

	var (
		val    = 0
		newVal = 0
		err    error
	)
	if testPrometheusToo {
		val, err = metrics.GetValue("num_iptables_rules")
		if err != nil {
			printPrometheusError(t, "Problem getting http metrics")
		}
	}

	if err := iptMgr.Delete(entry); err != nil {
		t.Errorf("TestDelete failed @ iptMgr.Delete")
	}

	if testPrometheusToo {
		newVal, err = metrics.GetValue("num_iptables_rules")
		if err != nil {
			printPrometheusError(t, "Problem getting http metrics")
		}
		if newVal != val-1 {
			printPrometheusError(t, "Delete iptable rule didn't register in prometheus")
		}
	}
}

func TestRun(t *testing.T) {
	iptMgr := &IptablesManager{}
	if err := iptMgr.Save(util.IptablesTestConfigFile); err != nil {
		t.Errorf("TestRun failed @ iptMgr.Save")
	}

	defer func() {
		if err := iptMgr.Restore(util.IptablesTestConfigFile); err != nil {
			t.Errorf("TestRun failed @ iptMgr.Restore")
		}
	}()

	iptMgr.OperationFlag = util.IptablesChainCreationFlag
	entry := &IptEntry{
		Chain: "TEST-CHAIN",
	}
	if _, err := iptMgr.Run(entry); err != nil {
		t.Errorf("TestRun failed @ iptMgr.Run")
	}
}

func TestMain(m *testing.M) {
	iptMgr := NewIptablesManager()
	iptMgr.Save(util.IptablesConfigFile)

	exitCode := m.Run()

	iptMgr.Restore(util.IptablesConfigFile)

	// messWithMetrics()

	os.Exit(exitCode)
}

func messWithMetrics() {
	go func() {
		for {
			metrics.Inc(metrics.NumPolicies)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for k := 0; k < 25; k++ {
			for j := 0; j < 2*k; j++ {
				metrics.Inc(metrics.NumIpSets)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for j := 0; j < 500; j += 2 {
			for k := 0; k < 2; k++ {
				metrics.Observe(metrics.AddPolicyExecTime, float64(2*k*j))
				time.Sleep(time.Second * time.Duration((k+1)/2))
			}
			for k := 0; k < 3; k++ {
				metrics.Observe(metrics.AddPolicyExecTime, float64(-k+j))
				time.Sleep(time.Second * time.Duration(k/3))
			}
		}
	}()

	go func() {
		for {
			for k := 0; k < 2; k++ {
				metrics.Observe(metrics.AddIpSetExecTime, float64(2*k))
				time.Sleep(time.Second * time.Duration((k+1)/2))
			}
			for k := 0; k < 3; k++ {
				metrics.Observe(metrics.AddIpSetExecTime, float64(-k))
				time.Sleep(time.Second * time.Duration(k+1))
			}
		}
	}()
}
