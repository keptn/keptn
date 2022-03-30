# Webhook Service 

## Overview

The **webhook service** is used to define webhooks - in the form of `curl` commands - for executing tasks of a task sequence.

## Configuring webhooks

To configure a webhook for a certain task within a sequence, a `webhook/webhook.yaml` file needs to be present in the 
configuration repository of the project/stage/service that should make use of the webhook. This file can be uploaded using the `keptn add-resource` command, 
or by using the `/resource` APIs of the `configuration-service`. An example for a webhook.yaml file is as follows:

```yaml
apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: my-subscription-id
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-k8s-secret"
            key: "my-key"
      requests:
        - "curl --header 'x-token: {{.env.secretKey}}' http://shipyard-controller:8080/v1/project/{{.data.project}}"
        - "curl http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}"
```

The example above will configure a webhook that should be executed whenever a `sh.keptn.event.mytask.triggered` event is received by the webhook service.
In this case, the following two `curl` requests will be executed:

```
curl --header 'x-token: {{.env.secret-key}}' http://shipyard-controller:8080/v1/project/{{.data.project}}
curl http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}
```

the responses of those requests will be stored in the `data.mytask.responses` property of the correlating 
`sh.keptn.event.mytask.finished` event that will be sent by the webhook service once the two requests have been finished:

```json
{
  "data": {
    "labels": null,
    "project": "webhooks",
    "result": "pass",
    "service": "myservice",
    "mytask": {
      "responses": [
        "{\"projects\":[]}",
        "{\"services\":[]}"
      ]
    },
    "stage": "dev",
    "status": "succeeded"
  },
  "id": "803151ae-98cd-49df-89a1-86d09581928a",
  "source": "webhook-service",
  "specversion": "1.0",
  "time": "2021-08-30T13:49:49.929Z",
  "type": "sh.keptn.event.mytask.finished",
  "shkeptncontext": "7a5a0757-9ccf-4f89-ad2b-0fd659eabefc",
  "triggeredid": "e6121823-81a5-43a7-a484-c456389ce88e"
}
```

As shown in the example above, the webhook.yaml file allows referencing secrets that are located in the same namespace as the 
Keptn control plane. Those secrets can then be used in the `curl` commands that should be executed for a certain task, using the `{{.env.<secret>}}` placeholder.
In addition to secrets, properties from incoming events, such as e.g. `{{.data.project}}`, `{{.shkeptncontext}}` etc. can be referenced using the template syntax.
Note that the execution of the defined requests will fail if any of the referenced values is not available.

### Disable automatic started/finished events

By default, the webhook service will send one `<task>.started` and one `<task>.finished` event for each received triggered event, where the `<task>.finished` event contains the aggregated responses 
of the executed requests. This behavior can be changed such that the responsibility of sending the `<task>.finished` events is moved to the services called by the webhook service. In this case,
the webhook service will send a `<task>.started` event for each of the curl requests that are to be executed. Afterwards, the requests are executed and, if they are successful, no `<task>.finished` event is sent by the webhook service.
If, however, one of the requests fails (e.g. due to an unknown environment variable, or if a disallowed curl command has been detected) the webhook service will send a `<task>.finished` event with `resutl=fail;status=errored` for this particular request and all requests that should have been executed afterwards. The remaining requests will not be executed in this case.
Sending the finished events can be disabled by setting `sendFinished` to `false` within the webhook configuration, e.g.:

```yaml
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.othertask.triggered"
      subscriptionID: my-subscription-id
      sendFinished: false
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - "curl http://shipyard-controller:8080/v1/project"
```

If, in addition to disabling the `<task>.finished` event also the `<task>.started` events should not be sent by the webhook service, the property `sendStarted` can be set to `false`:

```yaml
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.othertask.triggered"
      subscriptionID: my-subscription-id
      sendFinished: false
      sendStarted: false
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - "curl http://shipyard-controller:8080/v1/project"
```

### Enabling webhooks for a project, stage or service

If the same `webhook.yaml` file should be used across all stages and services within a project, the `webhook.yaml` file can be added as a project - resource:

```
keptn add-resource --project=my-project --resource=webhook.yaml --resourceUri=webhook/webhook.yaml
```

If a `webhook.yaml` should be used only for a certain stage, the optional `stage` parameter can be added to the `add-resource` command:

```
keptn add-resource --project=my-project --stage=my-stage --resource=webhook.yaml --resourceUri=webhook/webhook.yaml
```

Finally, if only a specific service should make use of the `webhook.yaml`, the `service` parameter has to be passed to the `add-resource` command:

```
keptn add-resource --project=my-project --stage=my-stage --service=my-service --resource=webhook.yaml --resourceUri=webhook/webhook.yaml
```
