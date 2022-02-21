import { SecretScope } from '../../../../../shared/interfaces/secret-scope';
import { Secret } from '../../../_models/secret';

const secrets = {
  Secrets: [
    Secret.fromJSON({
      name: 'keptn',
      scope: SecretScope.DEFAULT,
      keys: ['API_TOKEN'],
    }),
    Secret.fromJSON({
      name: 'webhook',
      scope: SecretScope.WEBHOOK,
      keys: ['API_TOKEN'],
    }),
    Secret.fromJSON({
      name: 'dynatrace',
      scope: 'dynatrace-service',
      keys: ['DT_TOKEN', 'DT_TENANT'],
    }),
  ],
};

export { secrets as SecretsResponseMock };
