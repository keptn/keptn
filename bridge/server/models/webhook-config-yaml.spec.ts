import { IWebhookV1Beta1, WebhookApiVersions } from '../interfaces/webhook-config-yaml-result';
import { migrateWebhook } from './webhook-config.utils';
import { WebhookConfigYaml } from './webhook-config-yaml';

const config: WebhookConfigYaml = WebhookConfigYaml.fromJSON({
  kind: 'WebhookConfig',
  apiVersion: WebhookApiVersions.V1BETA1,
  spec: {
    webhooks: [
      {
        subscriptionID: 'myID',
        type: 'sh.keptn.event.deployment.started',
        requests: [
          {
            url: 'https://keptn.sh/asdf',
            options: '--proxy https://keptn.sh/proxy',
            payload: '{"data": "myData"}',
            method: 'GET',
          },
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
const configWithLongCurl: WebhookConfigYaml = WebhookConfigYaml.fromJSON({
  kind: 'WebhookConfig',
  apiVersion: WebhookApiVersions.V1BETA1,
  spec: {
    webhooks: [
      {
        subscriptionID: 'myID',
        type: 'sh.keptn.event.deployment.started',
        requests: [
          {
            url: 'https://keptn.sh/asdf',
            method: 'GET',
            options: '--proxy https://keptn.sh/proxy',
            payload: '{"data": "myData"}',
            headers: [
              {
                key: 'content-type',
                value: 'application/json',
              },
              {
                key: 'Authorization',
                value:
                  'myVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryLongHeader',
              },
            ],
          },
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
const configWithoutSecrets: WebhookConfigYaml = WebhookConfigYaml.fromJSON({
  kind: 'WebhookConfig',
  apiVersion: WebhookApiVersions.V1BETA1,
  spec: {
    webhooks: [
      {
        subscriptionID: 'myID',
        type: 'sh.keptn.event.deployment.started',
        requests: [
          {
            url: 'https://keptn.sh/asdf',
            method: 'GET',
            headers: [
              {
                key: 'content-type',
                value: 'application/json',
              },
              {
                key: 'Authorization',
                value:
                  'myVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryLongHeader',
              },
            ],
            options: '--proxy https://keptn.sh/proxy',
            payload: '{"data": "myData"}',
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

describe('Test webhook-config-yaml', () => {
  it('should add another webhook', () => {
    // given
    const yamlConfig = getDefaultWebhookYaml(true);

    // when
    yamlConfig.addWebhook(
      {
        type: 'sh.keptn.events.approval.triggered',
        sendFinished: true,
        sendStarted: true,
        header: [],
        url: 'https://keptn.sh',
        method: 'GET',
        payload: '',
        proxy: 'asdf',
      },
      'mySecondID',
      []
    );

    // then
    expect(yamlConfig.spec.webhooks.length).toBe(2);
    expect(yamlConfig.spec.webhooks[1]).toEqual({
      type: 'sh.keptn.events.approval.triggered',
      sendFinished: true,
      sendStarted: true,
      subscriptionID: 'mySecondID',
      requests: [
        {
          method: 'GET',
          payload: '',
          headers: [],
          url: 'https://keptn.sh',
          options: '--proxy asdf',
        },
      ],
    } as IWebhookV1Beta1);
  });

  it('should update an existing webhook', () => {
    // given
    const yamlConfig = getDefaultWebhookYaml(false);
    expect(yamlConfig.spec.webhooks.length).toBe(1);

    // when
    yamlConfig.addWebhook(
      {
        type: 'sh.keptn.events.approval.started',
        sendFinished: true,
        sendStarted: true,
        url: 'https://keptn.sh/second',
        method: 'GET',
        proxy: '',
        payload: '',
        header: [],
      },
      'myID',
      [
        {
          name: 'mySecret',
          secretRef: {
            name: 'myName',
            key: 'myKey',
          },
        },
      ]
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
      requests: [
        {
          url: 'https://keptn.sh/second',
          method: 'GET',
          headers: [],
          payload: '',
        },
      ],
      subscriptionID: 'myID',
    } as IWebhookV1Beta1);
  });

  it('should update an existing webhook and remove secrets', () => {
    // given
    const yamlConfig = getDefaultWebhookYaml(false);

    // when
    yamlConfig.addWebhook(
      {
        type: 'sh.keptn.events.approval.started',
        url: 'https://keptn.sh/second',
        method: 'GET',
        proxy: 'https://keptn.sh/proxy',
        payload: '',
        header: [],
        sendStarted: true,
        sendFinished: false,
      },
      'myID',
      []
    );

    // then
    expect(yamlConfig.spec.webhooks.length).toBe(1);
    expect(yamlConfig.spec.webhooks[0]).toEqual({
      type: 'sh.keptn.events.approval.started',
      sendFinished: false,
      sendStarted: true,
      subscriptionID: 'myID',
      requests: [
        {
          url: 'https://keptn.sh/second',
          method: 'GET',
          options: '--proxy https://keptn.sh/proxy',
          headers: [],
          payload: '',
        },
      ],
    } as IWebhookV1Beta1);
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
    expect(result).toBe(`apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - subscriptionID: myID
      type: sh.keptn.event.deployment.started
      requests:
        - url: https://keptn.sh/asdf
          options: --proxy https://keptn.sh/proxy
          payload: >-
            {"data": "myData"}
          method: GET
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
    expect(result).toBe(`apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - subscriptionID: myID
      type: sh.keptn.event.deployment.started
      requests:
        - url: https://keptn.sh/asdf
          method: GET
          options: --proxy https://keptn.sh/proxy
          payload: >-
            {"data": "myData"}
          headers:
            - key: content-type
              value: application/json
            - key: Authorization
              value: myVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryLongHeader
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
    expect(result).toBe(`apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - subscriptionID: myID
      type: sh.keptn.event.deployment.started
      requests:
        - url: https://keptn.sh/asdf
          method: GET
          headers:
            - key: content-type
              value: application/json
            - key: Authorization
              value: myVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryLongHeader
          options: --proxy https://keptn.sh/proxy
          payload: >-
            {"data": "myData"}
      sendFinished: true
      sendStarted: true
`);
  });

  function getDefaultWebhookYaml(sendFinished?: boolean, sendStarted?: boolean): WebhookConfigYaml {
    return WebhookConfigYaml.fromJSON(
      migrateWebhook({
        kind: 'WebhookConfig',
        apiVersion: WebhookApiVersions.V1ALPHA1,
        spec: {
          webhooks: [
            {
              subscriptionID: 'myID',
              type: 'sh.keptn.event.deployment.started',
              requests: [
                `curl https://keptn.sh/asdf asdf --request GET --header 'content-type: application/json' --proxy https://keptn.sh/proxy --data '{"data": "myData"}'`,
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
      })
    );
  }
});
