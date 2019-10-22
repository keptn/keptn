package cmd

import (
	"bytes"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"os"
	"testing"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestConfigureMonitoringCmd(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	mocking = true
	/*
	   	const tmpSIFileName = "service-indicators.yaml"
	   	const tmpSOFileName = "service-objectives.yaml"
	   	const tmpRemediationFileName = "remediation.yaml"

	   	siContent := `indicators:
	   - metric: cpu_usage_sockshop_carts
	     source: Prometheus
	     query: avg(rate(container_cpu_usage_seconds_total{namespace="sockshop-$ENVIRONMENT",pod_name=~"carts-primary-.*"}[$DURATION]))
	     queryObject: []
	   - metric: request_latency_seconds
	     source: Prometheus
	     query: rate(requests_latency_seconds_sum{job='carts-sockshop-$ENVIRONMENT'}[$DURATION])/rate(requests_latency_seconds_count{job='carts-sockshop-$ENVIRONMENT'}[$DURATION])
	     queryObject: []
	   - metric: request_latency_dt
	     source: Dynatrace
	     query: ""
	     queryObject:
	     - key: timeseriesId
	       value: com.dynatrace.builtin:service.responsetime
	     - key: aggregation
	       value: AVG`

	   	soContent := `pass: 90
	   warning: 75
	   objectives:
	   - metric: request_latency_seconds
	     threshold: 0.8
	     timeframe: 5m
	     score: 25
	   - metric: request_latency_dt
	     threshold: 1e+06
	     timeframe: 5m
	     score: 25
	   - metric: cpu_usage_sockshop_carts
	     threshold: 0.02
	     timeframe: 5m
	     score: 50`

	   	remediationContent := `remediations:
	   - name: cpu_usage_sockshop_carts
	     actions:
	     - action: scaling
	       value: "+1"`

	   	ioutil.WriteFile(tmpSIFileName, []byte(siContent), 0644)
	   	ioutil.WriteFile(tmpSOFileName, []byte(soContent), 0644)
	   	ioutil.WriteFile(tmpRemediationFileName, []byte(remediationContent), 0644)


	*/
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"configure",
		"monitoring",
		"prometheus",
		"--project=sockshop",
		"--service=carts",
		/*
			"--service-indicators=" + tmpSIFileName,
			"--service-objectives=" + tmpSOFileName,
			"--remediation=" + tmpRemediationFileName,

		*/
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	/*
		os.Remove(tmpSIFileName)
		os.Remove(tmpSOFileName)
		os.Remove(tmpRemediationFileName)

	*/

	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
}
