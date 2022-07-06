import { SecretScope } from './secret-scope';

export interface SecretKeyValuePair {
  key: string;
  value: string;
}

interface ISecret {
  name: string;
  scope: SecretScope;
}

export interface IClientSecret extends ISecret {
  keys: string[];
}

export interface IServiceSecret extends ISecret {
  data: SecretKeyValuePair[];
}
