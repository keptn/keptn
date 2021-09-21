import { WebhookConfig } from '../../../../shared/models/webhook-config';

const config = new WebhookConfig();
config.method = 'POST';
config.url = 'https://keptn.sh';
config.payload = '{}';
config.header = [{name: 'Content-Type', value: 'application/json'}];
export { config as WebhookConfigMock };
