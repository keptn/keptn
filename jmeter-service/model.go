package main

/*
 * Example workload config
spec_version: '0.1.0'
workloads:
  - teststrategy: performance
    vuser: 200
    loopcount: 50
    script: load.jmx
  - teststrategy: performance_light
    vuser: 50
    loopcount: 10
    script: load.jmx
  - teststrategy: functional
    vuser: 1
    loopcount: 1
    script: func.jmx
*/
const (
	JMeterConfFilename       = "jmeter/jmeter.conf.yaml"
	TestStrategy_Performance = "performance"
	TestStrategy_Functional  = "functional"
	TestStrategy_HealthCheck = "healthcheck"
	TestStrategy_RealUser    = "real-user"
)

type JMeterConf struct {
	SpecVersion string      `json:"spec_version" yaml:"spec_version"`
	Workloads   []*Workload `json:"workloads" yaml:"workloads"`
}

type Workload struct {
	TestStrategy      string            `json:"teststrategy" yaml:"teststrategy"`
	VUser             int               `json:"vuser" yaml:"vuser"`
	LoopCount         int               `json:"loopcount" yaml:"loopcount"`
	ThinkTime         int               `json:"thinktime" yaml:"thinktime"`
	Script            string            `json:"script" yaml:"script"`
	AcceptedErrorRate float32           `json:"acceptederrorrate" yaml:"acceptederrorrate"`
	AvgRtValidation   int               `json:"avgrtvalidation" yaml:"avgrtvalidation"`
	Properties        map[string]string `json:"properties" yaml:"properties"`
}

var defaultWorkloads = []Workload{
	Workload{
		TestStrategy:      TestStrategy_HealthCheck,
		VUser:             1,
		LoopCount:         1,
		ThinkTime:         250,
		Script:            "jmeter/basiccheck.jmx",
		AcceptedErrorRate: 0.0,
	},
	Workload{
		TestStrategy:      TestStrategy_Performance,
		VUser:             10,
		LoopCount:         500,
		ThinkTime:         250,
		Script:            "jmeter/load.jmx",
		AcceptedErrorRate: 0.1,
	},
	Workload{
		TestStrategy:      TestStrategy_Functional,
		VUser:             1,
		LoopCount:         1,
		ThinkTime:         250,
		Script:            "jmeter/load.jmx",
		AcceptedErrorRate: 0.1,
	},
}
