import {Root} from './root';
import {Trace} from './trace';
import { Deployment } from './deployment';
import {EventTypes} from './event-types';

export class Service {
  serviceName: string;
  deployedImage: string;
  stage: string;
  allDeploymentsLoaded = false;
  deployments: Deployment[] = [];
  lastEventTypes: {[key: string]: {eventId: string, keptnContext: string, time: Date}};

  roots: Root[] = [];
  openApprovals: Trace[] = [];

  get deploymentContext(): string {
    return this.lastEventTypes?.[EventTypes.DEPLOYMENT_FINISHED]?.keptnContext ?? this.lastEventTypes?.[EventTypes.EVALUATION_FINISHED]?.keptnContext;
  }

  getShortImageName(): string {
    return this.deployedImage?.split('/').pop();
  }

  getImageVersion(): string {
    return this.deployedImage?.split(':').pop();
  }

  getOpenApprovals(): Trace[] {
    return this.openApprovals || [];
  }

  getOpenProblems(): Trace[] {
    return this.roots?.filter(root => root.isProblem() && !root.isProblemResolvedOrClosed()) || [];
  }

  getRecentSequence(): Root {
    return this.roots[0];
  }

  getRecentEvaluation(): Trace {
    return this.getRecentSequence()?.getEvaluation(this.stage);
  }

  static fromJSON(data: any) {
    return Object.assign(new this(), data);
  }
}
