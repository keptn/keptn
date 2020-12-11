import {Root} from "./root";

export class Service {
  serviceName: string;
  deployedImage: string;

  roots: Root[];

  getShortImageName() {
    return this.deployedImage.split("/").pop();
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
