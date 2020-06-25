package main

import (
	"time"

	"github.com/Azure/azure-container-networking/npm/metrics"
)

// file for testing metrics visually

func main() {
	messWithMetrics()
	metrics.StartHTTP(false)
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
				metrics.Inc(metrics.NumIPSets)
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
				metrics.Observe(metrics.AddIPSetExecTime, float64(2*k))
				time.Sleep(time.Second * time.Duration((k+1)/2))
			}
			for k := 0; k < 3; k++ {
				metrics.Observe(metrics.AddIPSetExecTime, float64(-k))
				time.Sleep(time.Second * time.Duration(k+1))
			}
		}
	}()
}
