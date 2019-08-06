# Gatekeeper Service

Service which implements the quality gate, i.e., depending on the evaluation result it either promotes an artifact to the next stage or not. Therefore, this service listens on `evaluation-done` events, which contains the result of the evaluation. In case the evaluation result is positive, this service sends a `new-artifact` event. In case the evaluation result is negative and the service is deployed with a b/g strategy, this service changes the configuration back to the old version and sends a `configuration-changed` event.
