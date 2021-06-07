import {Sequence} from './sequence';
import {Trace} from './trace';

export class DeploymentStage {
  public stageName: string;
  public remediations: Sequence[] = [];
  public evaluation: Trace;
  public evaluationContext: string;
  public config: string = null;

  constructor(stageName: string, evaluationContext: string) {
    this.stageName = stageName;
    this.evaluationContext = evaluationContext;
  }
}
