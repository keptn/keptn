import { Root } from './root';
import { DeploymentStage } from './deployment-stage';
import { Trace } from './trace';

export type DeploymentSelection = { deployment: Deployment; stage: string };

export class Deployment {
  public version?: string;
  public stages: DeploymentStage[];
  public service: string;
  public shkeptncontext: string;
  private _sequence?: Root;
  public name: string;

  constructor(version: string | undefined, service: string, stage: DeploymentStage, shkeptncontext: string) {
    this.version = version;
    this.service = service;
    this.name = version || service;
    this.stages = [stage];
    this.shkeptncontext = shkeptncontext;
  }

  set sequence(sequence: Root | undefined) {
    this._sequence = sequence;
    if (this._sequence) {
      for (const stage of this.stages) {
        stage.evaluation = this._sequence.getEvaluation(stage.stageName);
      }
    }
  }

  get sequence(): Root | undefined {
    return this._sequence;
  }

  public getEvaluation(stageName: string): Trace | undefined {
    return this.getStage(stageName)?.evaluation || this.sequence?.getEvaluation(stageName);
  }

  public getStage(stage: string): DeploymentStage | undefined {
    return this.stages.find((s) => s.stageName === stage);
  }

  public hasStage(stage: string): boolean {
    return this.stages.some((s) => s.stageName === stage);
  }

  public hasRemediation(stageName?: string): boolean {
    return stageName
      ? !!this.getStage(stageName)?.remediations.length
      : this.stages.some((s) => s.remediations.length !== 0);
  }

  public setEvaluation(evaluation: Trace | undefined) {
    if (evaluation?.stage) {
      const stage = this.getStage(evaluation.stage);
      if (stage) {
        stage.evaluation = evaluation;
      }
    }
  }
}
