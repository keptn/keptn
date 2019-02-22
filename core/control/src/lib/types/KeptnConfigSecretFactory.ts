import { KeptnGithubCredentials } from './KeptnGithubCredentials';
import { KeptnGithubCredentialsSecret } from './KeptnGithubCredentialsSecret';

import { base64encode, base64decode } from 'nodejs-base64';

export class KeptnConfigSecretFactory {

  constructor() { }

  createKeptnConfigSecret(creds: KeptnGithubCredentials): KeptnGithubCredentialsSecret {
    creds.org = base64encode(creds.org);
    creds.token = base64encode(creds.token);
    creds.user = base64encode(creds.user);
    const secret = {
      apiVersion: 'v1',
      kind: 'Secret',
      metadata: {
        name: 'github-credentials',
        namespace: 'keptn',
      },
      type: 'Opaque',
      data: creds,
    };

    return secret;
  }
}
