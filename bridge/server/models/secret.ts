import { Secret as scrt, SecretKeyValuePair } from '../../shared/interfaces/secret';
import { SecretScope } from '../../shared/interfaces/secret-scope';

export class Secret implements scrt {
  name!: string;
  scope!: SecretScope;
  data?: SecretKeyValuePair[];
  keys?: string[];

  public static fromJSON(data: unknown): Secret {
    return Object.assign(new this(), data);
  }
}
