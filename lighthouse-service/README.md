# lighthouse-service

This service is responsible for evaluating test results. 
Whenever an event of the type `sh.keptn.event.start-evaluation` (for the quality gates standalone use case) or `sh.keptn.events.tests-finished` is received,
it will take first determine the SLI (Service Level Indicator) data source (e.g., Prometheus or Dynatrace), 
as well as the required SLIs for the project. Afterwards, it will send out an event of the type `sh.keptn.internal.event.get-sli`, to inform
the respective data source services to retrieve the values for the SLIs defined in the `slo.yaml` file.

When a data source service is finished with the retrieval of the SLI values, and has sent them as an event of the type `sh.keptn.internal.event.get-sli.done`,
the lighthouse-service will evaluate the SLI values based on the evaluation strategy that has been defined in the  `slo.yaml` file.

# Configuring behavior

The service supports the environment variable `SLO_REQUIRED`, which accepts a boolean value: 

```yaml
env:
  - name: SLO_REQUIRED
    value: 'false'
```

* If the environment variable `SLO_REQUIRED` is not configured, it is automatically set to *true*. 

* If `SLO_REQUIRED` is set to `true` and no SLO file for the service is available, the lighthouse-service sends a `sh.keptn.events.evaluation-done` event with result **failed**. 

* If `SLO_REQUIRED` is set to `false` and an no SLO file for the service is available, the lighthouse-service sends a `sh.keptn.events.evaluation-done` event with result **warning**. 

# Configuring a data source
For each project, one data source (e.g., Prometheus or Dynatrace) can be defined. To tell Keptn which data source should be used, 
a config map with the name `lighthouse-config-<project-name>` and the following format 
needs to be deployed in the `keptn` namespace:

```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: lighthouse-config-<project-name>>
  namespace: keptn
data:
  sli-provider: "<name of the sli provider>"
```

## Example 1: Using Prometheus as a data source:
```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: lighthouse-config-sockshop
  namespace: keptn
data:
  sli-provider: "prometheus"
```

## Example 2: Using Dynatrace as a data source:
```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: lighthouse-config-sockshop
  namespace: keptn
data:
  sli-provider: "dynatrace"
```

# Defining Service Level Objectives (SLOs)

The required SLOs for a project can be defined by adding a file called `slo.yaml` to a service within a Keptn project, using the `keptn add-resource` command:

```
keptn add-resource --project=sockshop --stage=staging --service=carts --resource=examples/slo.yaml --resourceUri=slo.yaml
```

Note that the name of the file needs to be `slo.yaml`. If your file is called differently, you can ensure the correct naming of the 
file by using the `--resourceUri` parameter of the `keptn add-resource` command (as in the example above).

## Example SLO file content

The following SLO file is an extensive example of a SLO specification, using both fixed thresholds and comparison of previous evaluation results
as an evaluation strategy for the defined objectives:

```yaml

---
spec_version: '1.0'
# filter is optional
# specifies selection criteria for service in the SLI provider; project, stage,
# and service can be overwritten, if needed
filter:
  handler: "HealthCheckController.getHealth"
# comparison is mandatory
comparison:
  # compare_with is mandatory
  # possible values:
  # - single_result: only compare with one previous result
  # - several_results: compare with several previous results
  #   this option requires ‘number_of_comparison_results’
  compare_with: "single_result"
  # include_result_with_score is optional
  # default value: all
  # possible values:
  # - pass: only use previous results that had a ‘pass’ result in comparison
  # - pass_or_warn: only use previous results that had a ‘pass’ or a ‘warning’
  #   result in the comparison
  # - all: all previous values are used in the comparison
  include_result_with_score: "pass"
  # number_of_comparison_results is optional
  # default value: 3
  # possible values are positive integers greater than zero
  # if less than 3 values are available for comparison the evaluation fails
  number_of_comparison_results: 3
  # aggregate_function is optional
  # decides on the aggregate function which is applied to the previous results
  # before comparison
  # default value: avg
  # possible values:
  # - avg: average
  # - p90: 90th percentile
  # - p95: 95th percentile
  aggregate_function: avg
# objectives is mandatory
# describes the objectives for SLIs
objectives:
  # sli is mandatory
  # can be specified several times, if sli is specified without further attributes
  # the values are fetched and stored but are not taken into account for the
  # evaluation
  - sli: request_latency_p50
    # pass is optional
    # it defines the pass criteria for the SLI values
    pass:        # pass if (relative change <= 10% OR absolute value is < 200)
      # e.g.: If response time changes by more than 10%, it should still
      #       be considered as a pass if it is less than 200 ms
      - criteria:
          - "<=+10%" # relative values require a prefixed sign (plus or minus)
          - "<1000"   # absolute values only require a logical operator
    warning:     # allow small relative changes, and response time has to be < 500 ms
      - criteria:  # criteria connected by AND
          - "<=800"
  - sli: error_rate
    weight: 2   # default weight: 1
    pass:       # do not allow any security vulnerabilities
      - criteria:
          - "=0"
total_score:  # maximum score = sum of weights
  pass: "90%" # by default this is interpreted as ">="
  warning: "75%"
```
