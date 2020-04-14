import {Root} from "./root";

export class Service {
  serviceName: string;
  roots: Root[];

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
