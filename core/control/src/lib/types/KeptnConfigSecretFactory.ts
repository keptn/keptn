import { KeptnConfig } from './KeptnConfig';
import { KeptnConfigSecret } from './KeptnConfigSecret';

export class KeptnConfigSecretFactory {
  constructor() {}
  createKeptnConfigSecret(keptnConfig: KeptnConfig): KeptnConfigSecret {
    const secret = {
      apiVersion: 'v1',
      kind: 'Secret',
      metadata: {
        name: 'github-credentials',
        namespace: 'keptn',
      },
      type: 'Opaque',
      data: keptnConfig,
    };
    return secret;
  }
}
