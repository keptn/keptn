import { SecretScope } from './secret-scope';

export interface SecretKeyValuePair {
  key: string;
  value: string;
}

export interface ISecret {
  name: string;
  scope: SecretScope;
  data?: SecretKeyValuePair[];
  keys?: string[];
}
