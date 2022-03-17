import { WebhookConfigYaml } from './webhook-config-yaml';
import { WebhookConfig } from '../../shared/models/webhook-config';
import { Webhook } from '../interfaces/webhook-config-yaml-result';

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
        sendStarted: true,
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
        sendStarted: true,
      },
    ],
  },
  metadata: {
    name: 'webhook-configuration',
  },
});
const configWithoutSecrets = WebhookConfigYaml.fromJSON({
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
        sendFinished: true,
        sendStarted: true,
      },
    ],
  },
  metadata: {
    name: 'webhook-configuration',
  },
});

describe('Test webhook-config-yaml', () => {
  it('should parse request correctly', () => {
    // given
    const yamlConfig = getDefaultWebhookYaml();

    // when
    const result = yamlConfig.parsedRequest('myID') as WebhookConfig;

    // then
    const webhookConfig: WebhookConfig = Object.assign(new WebhookConfig(), {
      url: 'http://keptn.sh/asdf asdf',
      type: 'sh.keptn.event.deployment.started',
      header: [
        {
          name: 'content-type',
          value: 'application/json',
        },
      ],
      proxy: 'http://keptn.sh/proxy',
      payload: '{\n  "data": "myData"\n}',
      method: 'GET',
      secrets: [
        {
          name: 'mySecret',
          secretRef: {
            name: 'myName',
            key: 'myKey',
          },
        },
      ],
      sendFinished: false,
      sendStarted: false,
    });
    expect(result).toEqual(webhookConfig);
  });

  it('should set sendFinished to the value specified', () => {
    // given
    const yamlConfigTrue = getDefaultWebhookYaml(true, true);
    const yamlConfigFalse = getDefaultWebhookYaml(false, false);
    const yamlConfigUndefined = getDefaultWebhookYaml();

    // when
    const webhookConfigTrue = yamlConfigTrue.parsedRequest('myID') as WebhookConfig;
    const webhookConfigFalse = yamlConfigFalse.parsedRequest('myID') as WebhookConfig;
    const webhookConfigUndefined = yamlConfigUndefined.parsedRequest('myID') as WebhookConfig;

    // then
    expect(webhookConfigTrue.sendFinished).toBe(true);
    expect(webhookConfigFalse.sendFinished).toBe(false);
    expect(webhookConfigUndefined.sendFinished).toBe(false);
  });

  it('should set sendStarted to the value specified', () => {
    // given
    const yamlConfigTrue = getDefaultWebhookYaml(true, true);
    const yamlConfigFalse = getDefaultWebhookYaml(false, false);
    const yamlConfigUndefined = getDefaultWebhookYaml();

    // when
    const webhookConfigTrue = yamlConfigTrue.parsedRequest('myID') as WebhookConfig;
    const webhookConfigFalse = yamlConfigFalse.parsedRequest('myID') as WebhookConfig;
    const webhookConfigUndefined = yamlConfigUndefined.parsedRequest('myID') as WebhookConfig;

    // then
    expect(webhookConfigTrue.sendStarted).toBe(true);
    expect(webhookConfigFalse.sendStarted).toBe(false);
    expect(webhookConfigUndefined.sendStarted).toBe(true);
  });

  it('should return undefined if subscriptionID does not exist', () => {
    // given
    const yamlConfig = getDefaultWebhookYaml(true);
    // when
    const webhookConfig = yamlConfig.parsedRequest('notExistingID');
    // then
    expect(webhookConfig).toBeUndefined();
  });

  it('should add another webhook', () => {
    // given
    const yamlConfig = getDefaultWebhookYaml(true);

    // when
    yamlConfig.addWebhook(
      'sh.keptn.events.approval.triggered',
      'curl http://keptn.sh --request GET',
      'mySecondID',
      [],
      true,
      true
    );

    // then
    expect(yamlConfig.spec.webhooks.length).toBe(2);
    expect(yamlConfig.spec.webhooks[1]).toEqual({
      type: 'sh.keptn.events.approval.triggered',
      sendFinished: true,
      sendStarted: true,
      subscriptionID: 'mySecondID',
      requests: ['curl http://keptn.sh --request GET'],
    } as Webhook);
  });

  it('should update an existing webhook', () => {
    // given
    const yamlConfig = getDefaultWebhookYaml(false);
    expect(yamlConfig.spec.webhooks.length).toBe(1);

    // when
    yamlConfig.addWebhook(
      'sh.keptn.events.approval.started',
      'curl http://keptn.sh/second --request GET',
      'myID',
      [
        {
          name: 'mySecret',
          secretRef: {
            name: 'myName',
            key: 'myKey',
          },
        },
      ],
      true,
      true
    );

    // then
    expect(yamlConfig.spec.webhooks.length).toBe(1);
    expect(yamlConfig.spec.webhooks[0]).toEqual({
      type: 'sh.keptn.events.approval.started',
      sendFinished: true,
      sendStarted: true,
      envFrom: [
        {
          name: 'mySecret',
          secretRef: {
            name: 'myName',
            key: 'myKey',
          },
        },
      ],
      subscriptionID: 'myID',
      requests: ['curl http://keptn.sh/second --request GET'],
    } as Webhook);
  });

  it('should update an existing webhook and remove secrets', () => {
    // given
    const yamlConfig = getDefaultWebhookYaml(false);

    // when
    yamlConfig.addWebhook(
      'sh.keptn.events.approval.started',
      'curl http://keptn.sh/second --request GET',
      'myID',
      [],
      true,
      true
    );

    // then
    expect(yamlConfig.spec.webhooks.length).toBe(1);
    expect(yamlConfig.spec.webhooks[0]).toEqual({
      type: 'sh.keptn.events.approval.started',
      sendFinished: true,
      sendStarted: true,
      subscriptionID: 'myID',
      requests: ['curl http://keptn.sh/second --request GET'],
    } as Webhook);
  });

  it('should remove webhooks', () => {
    // given
    const yamlConfig = getDefaultWebhookYaml(false);
    expect(yamlConfig.spec.webhooks.length).toBe(1);
    expect(yamlConfig.hasWebhooks()).toBe(true);

    // when
    yamlConfig.removeWebhook('myID');

    // then
    expect(yamlConfig.hasWebhooks()).toBe(false);
    expect(yamlConfig.spec.webhooks.length).toBe(0);
  });

  it('should generate yaml correctly', () => {
    // when
    const result = config.toYAML();
    // then
    expect(result).toBe(`apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - subscriptionID: myID
      type: sh.keptn.event.deployment.started
      requests:
        - >-
          curl http://keptn.sh/asdf asdf --request GET --proxy http://keptn.sh/proxy --data '{"data": "myData"}'
      envFrom:
        - name: mySecret
          secretRef:
            name: myName
            key: myKey
      sendFinished: true
      sendStarted: true
`);
  });

  it('should generate yaml correctly with long curl without line breaks or escape characters', () => {
    // when
    const result = configWithLongCurl.toYAML();

    // then
    expect(result).toBe(`apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - subscriptionID: myID
      type: sh.keptn.event.deployment.started
      requests:
        - >-
          curl http://keptn.sh/asdf asdf --request GET --header 'content-type: application/json' --header 'Authorization: myVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryLongHeader' --proxy http://keptn.sh/proxy --data '{"data": "myData"}'
      envFrom:
        - name: mySecret
          secretRef:
            name: myName
            key: myKey
      sendFinished: true
      sendStarted: true
`);
  });

  it('should generate yaml without secrets', () => {
    // when
    const result = configWithoutSecrets.toYAML();

    // then
    expect(result).toBe(`apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - subscriptionID: myID
      type: sh.keptn.event.deployment.started
      requests:
        - >-
          curl http://keptn.sh/asdf asdf --request GET --header 'content-type: application/json' --header 'Authorization: myVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryLongHeader' --proxy http://keptn.sh/proxy --data '{"data": "myData"}'
      sendFinished: true
      sendStarted: true
`);
  });

  function getDefaultWebhookYaml(sendFinished?: boolean, sendStarted?: boolean): WebhookConfigYaml {
    return WebhookConfigYaml.fromJSON({
      kind: 'WebhookConfig',
      apiVersion: 'webhookconfig.keptn.sh/v1alpha1',
      spec: {
        webhooks: [
          {
            subscriptionID: 'myID',
            type: 'sh.keptn.event.deployment.started',
            requests: [
              `curl http://keptn.sh/asdf asdf --request GET --header 'content-type: application/json' --proxy http://keptn.sh/proxy --data '{"data": "myData"}'`,
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
            ...(sendFinished !== undefined && { sendFinished }),
            ...(sendStarted !== undefined && { sendStarted }),
          },
        ],
      },
      metadata: {
        name: 'webhook-configuration',
      },
    });
  }
});
