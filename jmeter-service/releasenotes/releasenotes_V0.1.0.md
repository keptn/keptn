# Release Notes 0.1.0

*jmeter-extended-service* builds on the status of the core *jenkins-service* and its capabilities to execute JMeter tests for functional and performance test strategy. 

## New Features

* Implemented support for custom teststragy specific workload definitions. For details please see readme.md

## Fixed Issues

* No fixed issues just new capabilities

## Known Limitations

* When defining large workloads, e.g: high number of Virtual Users you might need to adjust your pod memory & cpu limits as otherwise JMeter doesnt get enough resources and test results will be skewed
