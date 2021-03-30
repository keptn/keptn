import {Root} from "./root";
import {Trace} from "./trace";
import { Deployment } from './Deployment';

export class Service {
  serviceName: string;
  deployedImage: string;
  stage: string;
  allDeploymentsLoaded = false;
  deployments: Deployment[] = [];

  roots: Root[] = [];
  openApprovals: Trace[] = [];

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
    return Object.assign(new this, data);
  }
}
