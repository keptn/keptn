export class Deployment {
  public version: string;
  public stages: string[];
  public service: string;
  public gitCommit: string;

  constructor(version: string, service: string, stage: string) {
    this.version = version;
    this.service = service;
    this.stages = [stage];
  }

  static fromJSON(data: any): Deployment {
    return Object.assign(new this, data);
  }
}
