import { IWebhookConfigClient } from '../../../../shared/interfaces/webhook-config';

const config: IWebhookConfigClient = {
  method: 'POST',
  url: 'https://keptn.sh',
  payload: '{}',
  header: [{ key: 'Content-Type', value: 'application/json' }],
  sendFinished: false,
  sendStarted: false,
  type: 'sh.keptn.event.evaluation.triggered',
};
export { config as WebhookConfigMock };
