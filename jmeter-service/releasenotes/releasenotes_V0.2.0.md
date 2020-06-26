# Release Notes 0.2.0

*jmeter-extended-service* builds on the status of the core *jenkins-service* and its capabilities to execute JMeter tests for functional and performance test strategy. 

## New Features

* Implemented Feature Request #3 allowing a user to store *.jmx files on either service, stage or project level. The JMeter service will first look at service, then at stage and last at project level.

## Fixed Issues

* No fixed issues just new capabilities

## Known Limitations

* (just as 0.1.0) When defining large workloads, e.g: high number of Virtual Users you might need to adjust your pod memory & cpu limits as otherwise JMeter doesnt get enough resources and test results will be skewed
