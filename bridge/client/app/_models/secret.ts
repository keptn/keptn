class KeyValuePair {
  key: string;
  value: string;
}

export class Secret {
  name: string;
  scope: string;
  data: KeyValuePair[];

  constructor() {
    this.scope = 'keptn-default';
    this.data = [];
  }

  static fromJSON(data: unknown) {
    return Object.assign(new this(), data);
  }

  setName(name: string): void {
    this.name = name;
  }

  addData(key: string, value: string): void {
    this.data?.push({ key, value });
  }

  getData(index: number): KeyValuePair {
    return this.data[index];
  }

  removeData(index: number): void {
    this.data?.splice(index, 1);
  }
}
