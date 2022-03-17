import { WebhookConfig } from '../../shared/models/webhook-config';
import { generateWebhookConfigCurl, parseCurl } from './curl.utils';

describe('Test curl-parser', () => {
  it('should generate webhook curl', () => {
    const webhookConfig: WebhookConfig = Object.assign(new WebhookConfig(), {
      url: 'http://keptn.sh/asdf asdf',
      type: 'sh.keptn.event.deployment.started',
      header: [
        {
          name: 'content-type',
          value: 'application/json',
        },
        {
          name: 'Authorization',
          value: 'Bearer myToken',
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
    expect(generateWebhookConfigCurl(webhookConfig)).toBe(
      `curl --header 'content-type: application/json' --header 'Authorization: Bearer myToken' --request GET --proxy http://keptn.sh/proxy --data '{"data":"myData"}' http://keptn.sh/asdf asdf`
    );
  });

  it('should not add header if not provided', () => {
    const webhookConfig: WebhookConfig = Object.assign(new WebhookConfig(), {
      url: 'http://keptn.sh/asdf asdf',
      type: 'sh.keptn.event.deployment.started',
      proxy: 'http://keptn.sh/proxy',
      payload: '{\n  "data": "myData"\n}',
      method: 'GET',
      secrets: [],
      sendFinished: false,
      sendStarted: false,
    });
    expect(generateWebhookConfigCurl(webhookConfig)).toBe(
      `curl --request GET --proxy http://keptn.sh/proxy --data '{"data":"myData"}' http://keptn.sh/asdf asdf`
    );
  });

  it('should not add proxy if not provided', () => {
    const webhookConfig: WebhookConfig = Object.assign(new WebhookConfig(), {
      url: 'http://keptn.sh/asdf asdf',
      type: 'sh.keptn.event.deployment.started',
      proxy: '',
      payload: '{\n  "data": "myData"\n}',
      method: 'GET',
      secrets: [],
      sendFinished: false,
      sendStarted: false,
    });
    expect(generateWebhookConfigCurl(webhookConfig)).toBe(
      `curl --request GET --data '{"data":"myData"}' http://keptn.sh/asdf asdf`
    );
  });

  it('should not add data if not provided', () => {
    const webhookConfig: WebhookConfig = Object.assign(new WebhookConfig(), {
      url: 'http://keptn.sh/asdf',
      type: 'sh.keptn.event.deployment.started',
      proxy: '',
      payload: '',
      method: 'GET',
      secrets: [],
      sendFinished: false,
      sendStarted: false,
    });
    expect(generateWebhookConfigCurl(webhookConfig)).toBe(`curl --request GET http://keptn.sh/asdf`);
  });

  it('should remove new lines in payload', () => {
    const webhookConfig: WebhookConfig = Object.assign(new WebhookConfig(), {
      url: 'http://keptn.sh/asdf',
      type: 'sh.keptn.event.deployment.started',
      proxy: '',
      payload: '\r\n<meta>myText</meta>\n<meta2>myOtherText</meta2>',
      method: 'GET',
      secrets: [],
      sendFinished: false,
      sendStarted: false,
    });
    expect(generateWebhookConfigCurl(webhookConfig)).toBe(
      `curl --request GET --data '<meta>myText</meta><meta2>myOtherText</meta2>' http://keptn.sh/asdf`
    );
  });

  it('should parse curl correctly', () => {
    expect(
      parseCurl(
        `curl http://keptn.sh/asdf asdf --request GET --header 'content-type: application/json' --header 'Authorization: Bearer myToken' --proxy http://keptn.sh/proxy --data '{"data": "myData"}'`
      )
    ).toEqual({
      _: ['http://keptn.sh/asdf', 'asdf'],
      request: ['GET'],
      header: ['content-type: application/json', 'Authorization: Bearer myToken'],
      proxy: ['http://keptn.sh/proxy'],
      data: ['{"data": "myData"}'],
    });
  });
});
