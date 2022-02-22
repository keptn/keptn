export enum SecretScopeDefault {
  DEFAULT = 'keptn-default',
  WEBHOOK = 'keptn-webhook-service',
  DYNATRACE = 'dynatrace-service',
}

export type SecretScope = SecretScopeDefault | string;
