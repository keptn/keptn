import {Root} from './root';
import {DeploymentStage} from './deployment-stage';

export class Deployment {
  public version: string;
  public stages: DeploymentStage[];
  public service: string;
  public shkeptncontext: string;
  private _sequence: Root;
  public name: string;

  static fromJSON(data: any): Deployment {
    const deployment = Object.assign(new this(), data);
    deployment.name = deployment.version || deployment.service;
    return deployment;
  }

  set sequence(sequence: Root) {
    this._sequence = sequence;
    for (const stage of this.stages) {
      stage.evaluation = this.sequence.getEvaluation(stage.stageName);
    }
  }

  get sequence(): Root {
    return this._sequence;
  }

  public getEvaluation(stageName: string) {
    return this.getStage(stageName)?.evaluation || this.sequence.getEvaluation(stageName);
  }

  public getStage(stage: string): DeploymentStage {
    return this.stages.find(s => s.stageName === stage);
  }

  public hasStage(stage: string): boolean {
    return this.stages.some(s => s.stageName === stage);
  }

  public hasRemediation(stageName?: string): boolean {
    return stageName ? this.getStage(stageName)?.remediations.length > 0 : this.stages.some(s => s.remediations.length !== 0);
  }
}
