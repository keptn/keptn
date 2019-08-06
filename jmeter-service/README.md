# keptn Jmeter Service
Service for triggering JMeter tests. Therefore, this service listens on `deployment-finished` events and then starts the tests. In case the tests succeeed, this service sends a `test-finished` event. In case the tests do not succeed (e.g. the error rate is too high), this service sends an `evaluation-done` event with the data `evaluationpassed=false`.
