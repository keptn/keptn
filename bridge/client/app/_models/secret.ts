export class Secret {
  name: string;
  scope: string;
  data: [{
    key: string;
    value: string;
  }];

  constructor() {
    this.scope = "keptn-default";
  }

  addData() {
    if(!this.data)
      this.data = [{ key: "", value: "" }];
    else
      this.data.push({ key: "", value: "" });
  }

  removeData(index) {
    this.data.splice(index, 1);
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
