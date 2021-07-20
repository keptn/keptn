import {Sequence} from './sequence';
import {Trace} from './trace';

export class DeploymentStage {
  public stageName: string;
  public remediations: Sequence[] = [];
  public evaluation?: Trace;
  public evaluationContext?: string;
  public config?: string;

  constructor(stageName: string, evaluationContext: string | undefined) {
    this.stageName = stageName;
    this.evaluationContext = evaluationContext;
  }
}
