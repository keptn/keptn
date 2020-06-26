import {Service} from "./service";

export class Stage {
  stageName: string;
  services: Service[];

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }

}
