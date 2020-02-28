import {Service} from "./service";
import {Trace} from "./trace";

export class Stage {
  stageName: string;
  services: Service[];

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }

}
