import {Sequence} from './sequence';

export class DeploymentStage {
  public stageName: string;
  public remediations: Sequence[] = [];
  public config: string = null;

  constructor(stageName: string) {
    this.stageName = stageName;
  }
}
