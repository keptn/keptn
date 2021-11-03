import { Deployment as dp, IStageDeployment, SubSequence } from '../../../shared/interfaces/deployment';
import { EvaluationResult } from '../../../shared/interfaces/evaluation-result';
import { Trace } from './trace';
import { SequenceState } from '../../../shared/models/sequence';
import { ResultTypes } from '../../../shared/models/result-types';
import { Sequence } from './sequence';

export class StageDeployment implements IStageDeployment {
  name!: string;
  deploymentURL?: string;
  hasEvaluation!: boolean;
  evaluationResult?: EvaluationResult;
  latestEvaluation?: Trace;
  openRemediations!: Sequence[];
  approvalInformation?: {
    trace: Trace;
    latestImage?: string;
  };
  subSequences!: SubSequence[];

  public static fromJSON(data: IStageDeployment): StageDeployment {
    const stage: StageDeployment = Object.assign(new this(), data);
    if (stage.latestEvaluation) {
      stage.latestEvaluation = Trace.fromJSON(stage.latestEvaluation);
    }
    if (stage.approvalInformation?.trace) {
      stage.approvalInformation.trace = Trace.fromJSON(stage.approvalInformation.trace);
    }
    stage.openRemediations = stage.openRemediations.map((seq) => Sequence.fromJSON(seq));
    return stage;
  }

  public isSuccessful(): boolean {
    return this.subSequences.every((seq) => seq.state === SequenceState.FINISHED && seq.result === ResultTypes.PASSED);
  }

  public isWarning(): boolean {
    return this.subSequences.some((seq) => seq.result === ResultTypes.WARNING) && !this.isFaulty();
  }

  public isFaulty(): boolean {
    return this.subSequences.some((seq) => seq.result === ResultTypes.FAILED);
  }
}

export class Deployment implements dp {
  stages!: StageDeployment[];
  labels!: { [p: string]: string };
  state!: SequenceState;

  public static fromJSON(data: dp): Deployment {
    const deployment: Deployment = Object.assign(new this(), data);
    deployment.stages = deployment.stages.map((stage) => StageDeployment.fromJSON(stage));
    return deployment;
  }

  public getStage(stageName: string): StageDeployment | undefined {
    return this.stages.find((stage) => stage.name === stageName);
  }
}
