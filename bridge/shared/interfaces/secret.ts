export interface SecretKeyValuePair {
  key: string;
  value: string;
}

export enum SecretScope {
  DEFAULT = 'keptn-default',
  WEBHOOK = 'keptn-webhook-service'
}

export interface Secret {
  name: string;
  scope: SecretScope;
  data: SecretKeyValuePair[];
}
