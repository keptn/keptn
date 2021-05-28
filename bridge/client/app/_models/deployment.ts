import {Root} from './root';
import {Sequence} from './sequence';

export class Deployment {
  public version: string;
  public stages: {stageName: string, remediations: Sequence[], config: string }[];
  public service: string;
  public shkeptncontext: string;
  public sequence: Root;
  public name: string;

  static fromJSON(data: any): Deployment {
    const deployment = Object.assign(new this(), data);
    deployment.name = deployment.version || deployment.service;
    return deployment;
  }

  public getStage(stage: string): {stageName: string, remediations: Sequence[], config: string } {
    return this.stages.find(s => s.stageName === stage);
  }

  public hasStage(stage: string): boolean {
    return this.stages.some(s => s.stageName === stage);
  }

  public hasRemediation(stageName?: string): boolean {
    return stageName ? this.getStage(stageName)?.remediations.length > 0 : this.stages.some(s => s.remediations.length !== 0);
  }
}
