import {
  IWebhookConfigYamlResultV1Alpha1,
  IWebhookConfigYamlResultV1Beta1,
  WebhookApiVersions,
} from '../interfaces/webhook-config-yaml-result';
import {
  mapBridgeSecretsToYamlSecrets,
  mapYamlSecretsToBridgeSecrets,
  migrateWebhook,
  parseToClientWebhookRequest,
  stringifyPayload,
} from './webhook-config.utils';
import { WebhookConfigYaml } from './webhook-config-yaml';
import { IWebhookConfigClient } from '../../shared/interfaces/webhook-config';
import { IWebhookSecret } from '../interfaces/webhook-config';
import { ISecret } from '../../shared/interfaces/secret';
import { SecretScopeDefault } from '../../shared/interfaces/secret-scope';

describe('Test webhook-config-yaml', () => {
  it('should correctly migrate from v1alpha1 to v1beta1', () => {
    // given
    const yamlConfig = getDefaultV1Alpha1WebhookYaml(false, true);

    // when
    const result = migrateWebhook(yamlConfig);

    // then
    const webhookConfig: IWebhookConfigYamlResultV1Beta1 = {
      apiVersion: WebhookApiVersions.V1BETA1,
      kind: 'WebhookConfig',
      metadata: {
        name: 'webhook-configuration',
      },
      spec: {
        webhooks: [
          {
            envFrom: [
              {
                name: 'mySecret',
                secretRef: {
                  name: 'myName',
                  key: 'myKey',
                },
              },
            ],
            type: 'sh.keptn.event.deployment.started',
            sendFinished: false,
            sendStarted: true,
            subscriptionID: 'myID',
            requests: [
              {
                url: 'https://keptn.sh/asdf asdf',
                headers: [
                  {
                    key: 'content-type',
                    value: 'application/json',
                  },
                ],
                options: '--proxy https://keptn.sh/proxy',
                payload: '{\n  "data": "myData"\n}',
                method: 'GET',
              },
            ],
          },
        ],
      },
    };
    expect(result).toEqual(webhookConfig);
  });

  it('should not change or migrate v1beta1 webhook-yaml', () => {
    // given
    const yamlConfig = getDefaultV1Beta1WebhookYaml(false, true);

    // when
    const result = migrateWebhook(yamlConfig);

    // then
    // check if it hasn't changed
    expect(result).toEqual(getDefaultV1Beta1WebhookYaml(false, true));
    // check if it is the same reference
    expect(yamlConfig).toBe(result);
  });

  it('should parse config correctly for client', () => {
    // given
    const yamlConfigTrue = getDefaultWebhookYaml(true, true);

    // when
    const webhookConfigTrue = parseToClientWebhookRequest(yamlConfigTrue, 'myID');

    // then
    expect(webhookConfigTrue).toEqual({
      webhookConfig: {
        header: [
          {
            key: 'content-type',
            value: 'application/json',
          },
        ],
        url: 'https://keptn.sh/asdf asdf',
        sendFinished: true,
        sendStarted: true,
        type: 'sh.keptn.event.deployment.started',
        method: 'GET',
        payload: '{\n  "data": "myData"\n}',
        proxy: 'https://keptn.sh/proxy',
      },
      secrets: [
        {
          name: 'mySecret',
          secretRef: {
            name: 'myName',
            key: 'myKey',
          },
        },
      ],
    } as { webhookConfig: IWebhookConfigClient; secrets: IWebhookSecret[] });
  });

  it('should set sendFinished to the value specified', () => {
    // given
    const yamlConfigTrue = getDefaultWebhookYaml(true, true);
    const yamlConfigFalse = getDefaultWebhookYaml(false, false);
    const yamlConfigUndefined = getDefaultWebhookYaml();

    // when
    const webhookConfigTrue = parseToClientWebhookRequest(yamlConfigTrue, 'myID');
    const webhookConfigFalse = parseToClientWebhookRequest(yamlConfigFalse, 'myID');
    const webhookConfigUndefined = parseToClientWebhookRequest(yamlConfigUndefined, 'myID');

    // then
    expect(webhookConfigTrue?.webhookConfig.sendFinished).toBe(true);
    expect(webhookConfigFalse?.webhookConfig.sendFinished).toBe(false);
    expect(webhookConfigUndefined?.webhookConfig.sendFinished).toBe(false);
  });

  it('should format payload', () => {
    // given
    const obj = '{\n\r   "myProp": \t"myVal"\n}';

    // when
    const formatted = stringifyPayload(obj);

    // then
    expect(formatted).toBe('{"myProp":"myVal"}');
  });

  it('should not format payload', () => {
    // given
    const obj = '<meta>\n\rmyCustomData\n</meta>';

    // when
    const formatted = stringifyPayload(obj);

    // then
    expect(formatted).toBe('<meta>myCustomData</meta>');
  });

  it('should return undefined if subscriptionID does not exist', () => {
    // given
    const yamlConfig = getDefaultWebhookYaml(true);
    // when
    const webhookConfig = parseToClientWebhookRequest(yamlConfig, 'notExistingID');
    // then
    expect(webhookConfig).toBeUndefined();
  });

  it('should correctly map secrets of the yaml file to secrets of bridge UI', () => {
    // given
    const config = parseToClientWebhookRequest(getDefaultVWebhookConfigWithSecrets(), 'myID') as {
      webhookConfig: IWebhookConfigClient;
      secrets: IWebhookSecret[];
    };

    // when
    mapYamlSecretsToBridgeSecrets(config.webhookConfig, config.secrets);

    // then
    expect(config.webhookConfig.url).toBe('abc{{.secret.myName1.myKey1}}cde{{.secret.myName2.myKey2}}fg');
    expect(config.webhookConfig.payload).toBe('1{{.secret.myName1.myKey1}}2{{.secret.myName2.myKey2}}3');
    expect(config.webhookConfig.header[0].value).toBe(
      '.env.secret_myName1_myKey1{{.secret.myName1.myKey1}}.env.secret_myName2_myKey2{{.secret.myName2.myKey2}}.env.secret_myName1_myKey1'
    );
    expect(config.webhookConfig.header[1].value).toBe('{{.env.myNotExistingSecret}}');
  });

  it('should correctly map secrets of the client to secrets of the yaml file', () => {
    // given
    const config = getDefaultClientWebhookConfig();
    const secrets = getDefaultSecrets();

    // when
    mapBridgeSecretsToYamlSecrets(config, secrets);

    // then
    expect(config.url).toBe('abc{{.env.secret_myName1_myKey1}}cde{{.env.secret_myName2_myKey2}}fg');
    expect(config.payload).toBe('1{{.env.secret_myName1_myKey1}}2{{.env.secret_myName2_myKey2}}3');
    expect(config.header[0].value).toBe(
      '.secret.myName1.myKey1{{.env.secret_myName1_myKey1}}.secret.myName2.myKey2{{.env.secret_myName2_myKey2}}.secret.myName1.myKey1'
    );
    expect(config.header[1].value).toBe('{{.secret.myName1.myNotExistingKey}}');
  });

  function getDefaultV1Beta1WebhookYaml(
    sendFinished?: boolean,
    sendStarted?: boolean
  ): IWebhookConfigYamlResultV1Beta1 {
    return migrateWebhook(getDefaultV1Alpha1WebhookYaml(sendFinished, sendStarted));
  }

  function getDefaultWebhookYaml(sendFinished?: boolean, sendStarted?: boolean): WebhookConfigYaml {
    return WebhookConfigYaml.fromJSON(getDefaultV1Beta1WebhookYaml(sendFinished, sendStarted));
  }

  function getDefaultV1Alpha1WebhookYaml(
    sendFinished?: boolean,
    sendStarted?: boolean
  ): IWebhookConfigYamlResultV1Alpha1 {
    return {
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
    };
  }

  function getDefaultVWebhookConfigWithSecrets(): WebhookConfigYaml {
    return WebhookConfigYaml.fromJSON({
      kind: 'WebhookConfig',
      apiVersion: WebhookApiVersions.V1BETA1,
      spec: {
        webhooks: [
          {
            subscriptionID: 'myID',
            type: 'sh.keptn.event.deployment.started',
            requests: [
              {
                url: 'abc{{.env.secret_myName1_myKey1}}cde{{.env.secret_myName2_myKey2}}fg',
                options: '',
                headers: [
                  {
                    key: 'myKey',
                    value:
                      '.env.secret_myName1_myKey1{{.env.secret_myName1_myKey1}}.env.secret_myName2_myKey2{{.env.secret_myName2_myKey2}}.env.secret_myName1_myKey1',
                  },
                  {
                    key: 'myKey',
                    value: '{{.env.myNotExistingSecret}}',
                  },
                ],
                method: 'PUT',
                payload: '1{{.env.secret_myName1_myKey1}}2{{.env.secret_myName2_myKey2}}3',
              },
            ],
            envFrom: [
              {
                name: 'secret_myName1_myKey1',
                secretRef: {
                  name: 'myName1',
                  key: 'myKey1',
                },
              },
              {
                name: 'secret_myName2_myKey2',
                secretRef: {
                  name: 'myName2',
                  key: 'myKey2',
                },
              },
            ],
            sendStarted: true,
            sendFinished: true,
          },
        ],
      },
      metadata: {
        name: 'webhook-configuration',
      },
    });
  }

  function getDefaultSecrets(): ISecret[] {
    return [
      {
        name: 'myName1',
        scope: SecretScopeDefault.WEBHOOK,
        keys: ['myKey1'],
      },
      {
        name: 'myName2',
        scope: SecretScopeDefault.WEBHOOK,
        keys: ['myKey2'],
      },
    ];
  }

  function getDefaultClientWebhookConfig(): IWebhookConfigClient {
    return {
      url: 'abc{{.secret.myName1.myKey1}}cde{{.secret.myName2.myKey2}}fg',
      payload: '1{{.secret.myName1.myKey1}}2{{.secret.myName2.myKey2}}3',
      header: [
        {
          key: 'myKey1',
          value:
            '.secret.myName1.myKey1{{.secret.myName1.myKey1}}.secret.myName2.myKey2{{.secret.myName2.myKey2}}.secret.myName1.myKey1',
        },
        {
          key: 'myKey2',
          value: '{{.secret.myName1.myNotExistingKey}}',
        },
      ],
      type: 'sh.keptn.event.evaluation.triggered',
      method: 'PUT',
      sendStarted: true,
      sendFinished: true,
    };
  }
});
