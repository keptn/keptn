import { Secret as scrt, SecretKeyValuePair, SecretScope } from '../../shared/interfaces/secret';

export class Secret implements scrt {
  name!: string;
  scope!: SecretScope;
  data!: SecretKeyValuePair[];

  public static fromJSON(data: unknown): Secret {
    return Object.assign(new this(), data);
  }
}
