# Wait Service

The *wait service* is a keptn component that gets triggered by a `sh.keptn.events.deployment-finished` event and then waits for a certain duration. The duration is set by the environment variable `WAIT_DURATION`. The value of this variable must follow following pattern: **[duration][unit]**, e.g., 1h, 5m, 30s. 

After sleeping for this time, a `sh.keptn.events.tests-finished` event will be sent to keptn's eventbroker.