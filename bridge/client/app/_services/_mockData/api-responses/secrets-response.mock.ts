import { SecretScopeDefault } from '../../../../../shared/interfaces/secret-scope';

const secrets = {
  Secrets: [
    {
      name: 'keptn',
      scope: SecretScopeDefault.DEFAULT,
      keys: ['API_TOKEN'],
    },
    {
      name: 'webhook',
      scope: SecretScopeDefault.WEBHOOK,
      keys: ['API_TOKEN'],
    },
    {
      name: 'dynatrace',
      scope: SecretScopeDefault.DYNATRACE,
      keys: ['DT_TOKEN', 'DT_TENANT'],
    },
  ],
};

export { secrets as SecretsResponseMock };
