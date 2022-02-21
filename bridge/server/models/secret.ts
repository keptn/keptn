import { Secret as scrt, SecretKeyValuePair } from '../../shared/interfaces/secret';

export class Secret implements scrt {
  name!: string;
  scope!: string;
  data?: SecretKeyValuePair[];
  keys?: string[];

  public static fromJSON(data: unknown): Secret {
    return Object.assign(new this(), data);
  }
}
