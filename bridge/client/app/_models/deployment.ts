import {Root} from './root';

export class Deployment {
  public version: string;
  public stages: string[];
  public service: string;
  public shkeptncontext: string;
  public sequence: Root;

  static fromJSON(data: any): Deployment {
    return Object.assign(new this(), data);
  }
}
