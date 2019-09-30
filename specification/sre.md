# Specifications for Site Reliability Engineering with Keptn

To support site reliability engineering with Keptn and to enable the self-healing use case, Keptn relies on the specification of three file types:
* [Service Level Indicators (SLI)](#service-level-indicators-(sli))
* [Service Level Objectives (SLO)](#service-level-objectives-(slo))
* [Remediation Action](#remediation-action)

---

## Service Level Indicators (SLI)
The `service-indicators.yaml` file specifies the indicators that can be used to describe service objectives. These indicators are metrics gathered from monitoring sources and they are defined by a query. The query to obtain the metric is source-specific.

*Example:*
```yaml
indicators:
- metric: cpu_usage_sockshop_carts
  source: Prometheus
  query: avg(rate(container_cpu_usage_seconds_total{namespace="sockshop-$ENVIRONMENT",pod_name=~"carts-primary-.*"}[5m]))
- metric: request_latency_seconds
  source: Prometheus
  query: rate(requests_latency_seconds_sum{job='carts-sockshop-$ENVIRONMENT'}[$DURATION_MINUTESm])/rate(requests_latency_seconds_count{job='carts-sockshop-$ENVIRONMENT'}[$DURATION_MINUTESm])
```

([&uarr; up to index](#specifications-for-site-reliability-engineering-with-keptn))

## Service Level Objectives (SLO)
The `service-objectives.yaml` file specifies the service level objectives for one or more services. Therefore, this file first defines thresholds that express the fullfillment of the objectives by the `pass` and `warning` property. An evaluated objective that achives a score above the `pass` limit is considered to be fullfilled, between `warning` and `pass` it is in an acceptable range, and below `warning` it is not fullfilled. The `objectives` property lists all service level indicators by their metric name that are considered for this objective. Besides, each indicator is augmented by a `threshold` and `timeframe`. While the `threshold` defines the acceptance criteria of this indicator, the timeframe indicates the duration in which the metrics is evaluated. Finally, the `score` specifies the max number of points that can be achieved by this indicator. 

*Example:*
```yaml
pass: 90
warning: 75
objectives:
- metric: request_latency_seconds
  threshold: 0.8
  timeframe: 5m
  score: 50
- metric: cpu_usage_sockshop_carts
  threshold: 0.2
  timeframe: 5m
  score: 50
```

([&uarr; up to index](#specifications-for-site-reliability-engineering-with-keptn))

## Remediation Action
The `remediation.yaml` file defines remediation actions to execute in response to a problem related to the defined problem pattern / service objective. This action is interpreted by Keptn to trigger the proper remediation. 

*Example:*
```yaml
remediations:
- name: cpu_usage_sockshop_carts
  actions:
  - action: scaling
    value: +1
```

([&uarr; up to index](#specifications-for-site-reliability-engineering-with-keptn))