export class Secret {
  name: string;
  scope: string;
  data: {
    key: string;
    value: string;
  }?];

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

  removeData(index: number): void {
    this.data?.splice(index, 1);
  }
}
