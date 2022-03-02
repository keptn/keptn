import { WebhookConfigYaml } from './webhook-config-yaml';

const config = WebhookConfigYaml.fromJSON({
  kind: 'WebhookConfig',
  apiVersion: 'webhookconfig.keptn.sh/v1alpha1',
  spec: {
    webhooks: [
      {
        subscriptionID: 'myID',
        type: 'sh.keptn.event.deployment.started',
        requests: [
          `curl http://keptn.sh/asdf asdf --request GET --proxy http://keptn.sh/proxy --data '{"data": "myData"}'`,
        ],
        envFrom: [
          {
            name: 'mySecret',
            secretRef: {
              name: 'myName',
              key: 'myKey',
            },
          },
        ],
        sendFinished: true,
      },
    ],
  },
  metadata: {
    name: 'webhook-configuration',
  },
});

const configWithLongCurl = WebhookConfigYaml.fromJSON({
  kind: 'WebhookConfig',
  apiVersion: 'webhookconfig.keptn.sh/v1alpha1',
  spec: {
    webhooks: [
      {
        subscriptionID: 'myID',
        type: 'sh.keptn.event.deployment.started',
        requests: [
          `curl http://keptn.sh/asdf asdf --request GET --header 'content-type: application/json' --header 'Authorization: myVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryLongHeader' --proxy http://keptn.sh/proxy --data '{"data": "myData"}'`,
        ],
        envFrom: [
          {
            name: 'mySecret',
            secretRef: {
              name: 'myName',
              key: 'myKey',
            },
          },
        ],
        sendFinished: true,
      },
    ],
  },
  metadata: {
    name: 'webhook-configuration',
  },
});

describe('Test webhook-config-yaml', () => {
  it('should generate yaml correctly', () => {
    // when
    const result = config.toYAML();
    // then
    expect(result).toBe(
      `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - subscriptionID: myID
      type: sh.keptn.event.deployment.started
      requests:
        - "curl http://keptn.sh/asdf asdf --request GET --proxy http://keptn.sh/proxy --data '{\\"data\\": \\"myData\\"}'"
      envFrom:
        - name: mySecret
          secretRef:
            name: myName
            key: myKey
      sendFinished: true
`
    );
  });

  it('should generate yaml correctly with long curl without line breaks', () => {
    // when
    const result = configWithLongCurl.toYAML();

    // then
    expect(result).toBe(
      `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - subscriptionID: myID
      type: sh.keptn.event.deployment.started
      requests:
        - "curl http://keptn.sh/asdf asdf --request GET --header 'content-type: application/json' --header 'Authorization: myVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryLongHeader' --proxy http://keptn.sh/proxy --data '{\\"data\\": \\"myData\\"}'"
      envFrom:
        - name: mySecret
          secretRef:
            name: myName
            key: myKey
      sendFinished: true
`
    );
  });
});
