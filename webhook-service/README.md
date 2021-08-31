# Webhook Service 

## Overview

The **webhook service** is used to define webhooks - in the form of `curl` commands - for executing tasks of a task sequence.

## Configuring webhooks

To configure a webhook for a certain task within a sequence, a `webhook.yaml` file needs to be present in the 
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
      envFrom: 
        - name: "secret-key"
          secretRef:
            name: "my-k8s-secret"
            key: "my-key"
      requests:
        - "curl --header 'x-token: {{.env.secret-key}}' http://shipyard-controller:8080/v1/project/{{.data.project}}"
        - "curl http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}"
```

The example above will configure a webhook that should be executed whenever a `sh.keptn.event.mytask.triggered` event is received by the webhook service.
In this case, the following two `curl` requests will be executed:

```
curl --header 'x-token: {{.env.secret-key}}' http://shipyard-controller:8080/v1/project/{{.data.project}}
curl http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}
```

the responses of those requests will be stored in the `data["sh.keptn.event.mytask.triggered"].responses` property of the correlating 
`sh.keptn.event.mytask.finished` event that will be sent by the webhook service once the two requests have been finished:

```
{
  "data": {
    "labels": null,
    "project": "webhooks",
    "result": "pass",
    "service": "myservice",
    "sh.keptn.event.mytask.triggered": {
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
