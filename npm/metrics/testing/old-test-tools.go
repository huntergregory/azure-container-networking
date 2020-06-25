package metrics

import (
	"testing"
)

// DIDN'T WORK
func GaugeIncTest(t *testing.T, metricName string, action func()) {
	// if !started {
	// 	StartHTTP(true)
	// 	time.Sleep(2 * time.Second)
	// }
	// val, err := GetValue(metricName)
	// action()

	// if err != nil {
	// 	t.Errorf("Problem getting http prometheus metrics for metric: " + metricName)
	// 	return
	// }

	// newVal, err := GetValue(metricName)
	// fmt.Println(val)
	// fmt.Println(newVal)
	// if err != nil {
	// 	t.Errorf("Problem getting http prometheus metrics for metric: " + metricName)
	// }
	// if newVal != val+1 {
	// 	t.Errorf("Metric adjustment didn't register in prometheus for metric: " + metricName)
	// }
}
