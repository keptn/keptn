export class Deployment {
  public version: string;
  public stages: string[];
  public service: string;
  public gitCommit: string;

  static fromJSON(data: any): Deployment {
    return Object.assign(new this, data);
  }
}
