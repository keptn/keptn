import {Trace} from "./trace";
import {EventTypes} from "./event-types";
import {Stage} from "./stage";

export class Root extends Trace {
  traces: Trace[] = [];

  isFaulty(): string {
    return this.traces.reduce((result: string, trace: Trace) => trace.isFaulty() ? trace.data.stage : result, null);
  }

  isStarted(): boolean {
    return this.traces.length === 0 ? false : this.traces[this.traces.length-1].isStarted();
  }

  isFailedEvaluation(): string {
    let result: string = null;
    if(this.traces) {
      this.traces.forEach((trace) => {
        if(trace.isFailedEvaluation()) {
          result = trace.data.stage;
        }
      });
    }
    return result;
  }

  isDeployment(): string {
    return this.traces.reduce((result: string, trace: Trace) => result ? result : trace.isDeployment(), null);
  }

  isWarning(): string {
    return this.traces.reduce((result: string, trace: Trace) => trace.isWarning() ? trace.data.stage : result, null);
  }

  isSuccessful(): string {
    return this.traces.reduce((result: string, trace: Trace) => {
      if(result)
        return trace.isFaulty() ? null : result;
      else
        return trace.isSuccessful() ? trace.data.stage : result
    }, null);
  }

  hasPendingApproval(stage: string): boolean {
    const tracesOfStage = this.getTracesOfStage(stage);
    let pending = undefined;

    for(let i = 0; i < tracesOfStage.length && pending === undefined; ++i){
      if(tracesOfStage[i].isApproval()){
        pending = tracesOfStage[i].isApprovalPending();
      }
    }
    return pending === undefined ? false : pending;
  }

  getPendingApprovals(stageName?: string): Trace[] {
    return this.traces.filter(trace => trace.isApproval() && trace.isApprovalPending() && (!stageName || trace.getStage() == stageName));
  }

  getLastTrace(): Trace {
    return this.traces ? this.traces[this.traces.length - 1] : null;
  }

  getTracesOfStage(stage: string): Trace[] {
    return this.traces?.filter(trace => trace.data.stage === stage);
  }

  getFirstTraceOfStage(stage: string): Trace {
    return this.getTracesOfStage(stage)?.[0];
  }

  getLastTraceOfStage(stage: string): Trace {
    let traces = this.getTracesOfStage(stage);
    return traces ? traces[traces.length-1] : null;
  }

  getStages(): string[] {
    let result: string[] = [];
    if(this.traces) {
      this.traces.forEach((trace) => {
        if(trace.data.stage && result.indexOf(trace.data.stage) == -1)
          result.push(trace.data.stage);
      });
    }
    return result;
  }

  getProject(): string {
    if(!this.data.project)
      this.data.project = this.traces.find(trace => !!trace.data.project).data.project;
    return this.data.project;
  }

  getService(): string {
    if(!this.data.service)
      this.data.service = this.traces.find(trace => !!trace.data.project).data.service;
    return this.data.service;
  }

  getEvaluation(stageName: String): Trace {
    return this.traces.find(t => t.type == EventTypes.EVALUATION_TRIGGERED && t.data.stage == stageName);
  }

  getDeploymentDetails(stage: Stage): Trace {
    return this.traces.find(t => t.type == EventTypes.DEPLOYMENT_TRIGGERED && t.data.stage == stage.stageName)?.getFinishedEvent();
  }

  getRemediationActions(): Root[] {
    // create chunks of Remediations and start new chunk at REMEDIATION_TRIGGERED event
    return this.traces.reduce((result, trace: Trace) => {
      if(trace.type == EventTypes.ACTION_TRIGGERED)
        result.push(Root.fromJSON(JSON.parse(JSON.stringify(trace))));
      else if(result.length)
        result[result.length-1].traces = [...result[result.length-1].traces||[], trace];
      return result;
    }, []);
  }

  isFinished() {
    return this.traces.every(t => t.isFinished());
  }

  getSequenceName() {
    return this.type;
  }

  getStatus() {
    if(this.isProblem() && this.isProblemResolvedOrClosed())
      return "resolved";
    else if(this.isProblem())
      return "opened";
    else if(this.isFinished() && this.isFaulty())
      return "failed";
    else if(this.isFinished())
      return "succeeded";
    else
      return "active";
  }

  getStatusLabel() {
    switch(this.getStatus()) {
      case "resolved":
        return "resolved";
      case "opened":
        return "opened";
      case "failed":
        return "failed";
      case "succeeded":
        return "succeeded";
      case "active":
        if(this.getPendingApprovals().length > 0)
          return "waiting for approval";
        else
          return "started";
    }
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
