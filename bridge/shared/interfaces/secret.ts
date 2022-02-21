import { SecretScope } from './secret-scope';

export interface SecretKeyValuePair {
  key: string;
  value: string;
}

export interface Secret {
  name: string;
  scope: SecretScope | string;
  data?: SecretKeyValuePair[];
  keys?: string[];
}
