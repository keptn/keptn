import { Deployment as dp, IStageDeployment, SubSequence } from '../../../shared/interfaces/deployment';
import { EvaluationResult } from '../../../shared/interfaces/evaluation-result';
import { Trace } from './trace';
import { SequenceState } from '../../../shared/models/sequence';
import { ResultTypes } from '../../../shared/models/result-types';
import { Sequence } from './sequence';
import { ServiceRemediationInformation } from '../_interfaces/service-remediation-information';

export class StageDeployment implements IStageDeployment {
  name!: string;
  deploymentURL?: string;
  hasEvaluation!: boolean;
  lastTimeUpdated!: number;
  evaluationResult?: EvaluationResult;
  latestEvaluation?: Trace;
  openRemediations!: Sequence[];
  remediationConfig?: string;
  approvalInformation?: {
    trace: Trace;
    deployedImage?: string;
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

  public removeApproval(): void {
    this.approvalInformation = undefined;
    for (const subSequence of this.subSequences) {
      subSequence.hasPendingApproval = false;
    }
  }
}

export class Deployment implements dp {
  stages!: StageDeployment[];
  image?: string;
  keptnContext!: string;
  service!: string;
  labels!: { [p: string]: string };
  state!: SequenceState;

  public static fromJSON(data: unknown): Deployment {
    const deployment: Deployment = Object.assign(new this(), data);
    deployment.stages = deployment.stages.map((stage) => StageDeployment.fromJSON(stage));
    return deployment;
  }

  public getStage(stageName: string): StageDeployment | undefined {
    return this.stages.find((stage) => stage.name === stageName);
  }

  public get latestTimeUpdated(): number {
    return this.stages.reduce(
      (maxTime, stage) => (maxTime > stage.lastTimeUpdated ? maxTime : stage.lastTimeUpdated),
      0
    );
  }

  public updateRemediations(remediationInfo: ServiceRemediationInformation): void {
    for (const stage of this.stages) {
      const newStage = remediationInfo.stages.find((st) => st.name === stage.name);
      if (newStage && stage.deploymentURL) {
        stage.openRemediations = newStage.remediations;
        stage.remediationConfig = newStage.config;
      } else {
        stage.openRemediations = [];
        stage.remediationConfig = undefined;
      }
    }
  }

  public update(deployment: Deployment): void {
    this.image = deployment.image;
    this.labels = deployment.labels;
    this.state = deployment.state;
    for (const stage of deployment.stages) {
      const originalStage = this.stages.find((s) => s.name === stage.name);
      if (!originalStage) {
        // new stage
        this.stages.push(stage);
      } else {
        // update existing stage
        originalStage.lastTimeUpdated = stage.lastTimeUpdated;
        originalStage.approvalInformation = stage.approvalInformation;
        originalStage.remediationConfig = stage.remediationConfig;
        originalStage.openRemediations = stage.openRemediations;
        originalStage.deploymentURL ??= stage.deploymentURL;
        originalStage.evaluationResult ??= stage.evaluationResult;

        if ((stage.hasEvaluation && !originalStage.hasEvaluation) || !originalStage.latestEvaluation) {
          originalStage.latestEvaluation = stage.latestEvaluation;
        }
        originalStage.hasEvaluation = stage.hasEvaluation;

        this.updateSubSequences(originalStage, stage);
      }
    }
  }

  private updateSubSequences(originalStage: StageDeployment, newStage: StageDeployment): void {
    if (!originalStage.subSequences.length) {
      originalStage.subSequences = newStage.subSequences;
    } else {
      for (let i = newStage.subSequences.length - 1; i >= 0; i--) {
        const subSequence = newStage.subSequences[i];
        const originalSubSequence = originalStage.subSequences.find((seq) => seq.id === subSequence.id);
        if (originalSubSequence) {
          // update existing subSequence
          originalSubSequence.state = subSequence.state;
          originalSubSequence.result = subSequence.result;
          originalSubSequence.hasPendingApproval = subSequence.hasPendingApproval;
          originalSubSequence.message = subSequence.message;
        } else {
          // add new subSequences
          originalStage.subSequences.unshift(subSequence);
        }
      }
    }
  }
}
