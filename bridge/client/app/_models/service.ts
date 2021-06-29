import {Root} from './root';
import {Trace} from './trace';
import { Deployment } from './deployment';
import {EventTypes} from './event-types';
import {Sequence} from './sequence';
import {EvaluationResult} from './evaluation-result';

export class Service {
  serviceName: string;
  deployedImage: string;
  stage: string;
  allDeploymentsLoaded = false;
  deployments: Deployment[] = [];
  lastEventTypes: {[key: string]: {eventId: string, keptnContext: string, time: number}};

  sequences: Sequence[] = [];
  roots: Root[] = [];
  openApprovals: Trace[] = [];

  static fromJSON(data: any) {
    return Object.assign(new this(), data);
  }

  get deploymentContext(): string {
    return this.lastEventTypes?.[EventTypes.DEPLOYMENT_FINISHED]?.keptnContext ?? this.evaluationContext;
  }

  get deploymentTime(): number {
    return this.lastEventTypes?.[EventTypes.DEPLOYMENT_FINISHED]?.time || this.lastEventTypes?.[EventTypes.EVALUATION_FINISHED]?.time;
  }

  get evaluationContext(): string {
    return this.lastEventTypes?.[EventTypes.EVALUATION_FINISHED]?.keptnContext;
  }
  public getShortImageName() {
    return this.deployedImage?.split('/').pop().split(':').find(() => true);
  }

  getImageName(): string {
    return this.deployedImage?.split('/').pop();
  }

  getImageVersion(): string {
    return this.deployedImage?.split(':').pop();
  }

  getOpenApprovals(): Trace[] {
    return this.openApprovals || [];
  }

  getOpenProblems(): Trace[] {
    // show running remediation or last faulty remediation
    return this.roots?.filter((root, index) => root.isRemediation() && (!root.isFinished() || root.isFaulty() && index === 0)) || [];
  }

  getRecentRoot(): Root {
    return this.roots[0];
  }

  getRecentEvaluation(): Trace {
    return this.getRecentRoot()?.getEvaluation(this.stage);
  }

  public hasRemediations(): boolean {
    return this.deployments.some(d => d.stages.some(s => s.remediations.length !== 0));
  }
}
