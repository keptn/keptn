import { Secret as scrt, SecretKeyValuePair, SecretScope } from '../../../shared/interfaces/secret';

export class Secret implements scrt {
  name!: string;
  scope!: SecretScope;
  data: SecretKeyValuePair[];

  constructor() {
    this.scope = SecretScope.DEFAULT;
    this.data = [];
  }

  static fromJSON(data: unknown): Secret {
    return Object.assign(new this(), data);
  }

  setName(name: string): void {
    this.name = name;
  }

  addData(key: string, value: string): void {
    this.data?.push({key, value});
  }

  getData(index: number): SecretKeyValuePair {
    return this.data[index];
  }

  removeData(index: number): void {
    this.data?.splice(index, 1);
  }
}
