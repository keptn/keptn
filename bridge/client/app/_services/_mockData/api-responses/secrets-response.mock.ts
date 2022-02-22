import { SecretScopeDefault } from '../../../../../shared/interfaces/secret-scope';
import { Secret } from '../../../_models/secret';

const secrets = {
  Secrets: [
    Secret.fromJSON({
      name: 'keptn',
      scope: SecretScopeDefault.DEFAULT,
      keys: ['API_TOKEN'],
    }),
    Secret.fromJSON({
      name: 'webhook',
      scope: SecretScopeDefault.WEBHOOK,
      keys: ['API_TOKEN'],
    }),
    Secret.fromJSON({
      name: 'dynatrace',
      scope: SecretScopeDefault.DYNATRACE,
      keys: ['DT_TOKEN', 'DT_TENANT'],
    }),
  ],
};

export { secrets as SecretsResponseMock };
