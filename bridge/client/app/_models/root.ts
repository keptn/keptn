import {Trace} from "./trace";
import {EventTypes} from "./event-types";

export class Root extends Trace {
  traces: Trace[] = [];

  isFaulty(stageName?: string): boolean {
    // a Sequence is faulty, if there is a sequence that is faulty, but no other sequence that is successful on the same stage
    let result = false;
    const stages = stageName ? [stageName] : this.getStages();
    if (stages.length > 0) {
      stages.forEach(stage => {
        result ||= this.areTracesFaulty(stage);
      });
    }
    else {
      result = this.areTracesFaulty();
    }
    return result;
  }

  private areTracesFaulty(stageName?: string): boolean {
    const stageTraces = stageName ? this.traces.filter(t => t.getStage() === stageName) : this.traces;
    return stageTraces.some(t => t.isFaulty()) && !stageTraces.some(t => t.isSuccessful());
  }

  isStarted(): boolean {
    return this.traces.length === 0 ? false : this.traces[this.traces.length-1].isStarted();
  }

  hasFailedEvaluation(): string {
    let result: string = null;
    if(this.traces) {
      let failedEvaluation = this.findTrace(t => t.isEvaluation() && t.isFailedEvaluation());
      if(failedEvaluation)
        result = failedEvaluation.getStage();
    }
    return result;
  }

  getDeploymentTrace(stage: string): Trace {
    return this.findTrace(trace => trace.isDeployment() === stage);
  }

  isDeployment(): string {
    return this.traces.reduce((result: string, trace: Trace) => result ? result : trace.isDeployment(), null);
  }

  isWarning(stageName?: string): boolean {
    return this.traces.reduce((result: boolean, trace: Trace) => trace.isWarning(stageName), false);
  }

  isSuccessful(stageName?: string): boolean {
    return this.traces.reduce((result: boolean, trace: Trace) => {
      if (result)
        return !trace.isFaulty(stageName);
      else
        return trace.isSuccessful(stageName);
    }, false);
  }

  hasPendingApproval(stage: string): boolean {
    const tracesOfStage = this.getTracesOfStage(stage);
    let pending = undefined;

    for(let i = 0; i < tracesOfStage.length && pending === undefined; ++i) {
      if(tracesOfStage[i].getLastTrace().isApproval()) {
        pending = tracesOfStage[i].isApprovalPending();
      }
    }
    return pending === undefined ? false : pending;
  }

  getPendingApproval(stageName?: string): Trace {
    return this.findTrace(trace => trace.isApproval() && trace.isApprovalPending() && (!stageName || trace.getStage() == stageName));
  }

  getTracesOfStage(stage: string): Trace[] {
    return this.traces?.filter(trace => trace.data.stage === stage);
  }

  getFirstTraceOfStage(stage: string): Trace {
    return this.getTracesOfStage(stage)?.[0];
  }

  getLastTraceOfStage(stage: string): Trace {
    let traces = this.getTracesOfStage(stage);
    return traces ? traces[traces.length-1].getLastTrace() : null;
  }

  getLastSequenceOfStage(stage: string): Trace {
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
      this.data.project = this.findTrace(trace => !!trace.data.project).data.project;
    return this.data.project;
  }

  getService(): string {
    if(!this.data.service)
      this.data.service = this.findTrace(trace => !!trace.data.project).data.service;
    return this.data.service;
  }

  getEvaluation(stageName: String): Trace {
    return this.findLastTrace(trace => trace.type == EventTypes.EVALUATION_TRIGGERED && trace.data.stage == stageName && trace.traces.some(t => t.type == EventTypes.EVALUATION_STARTED));
  }

  getDeploymentDetails(stage: string): Trace {
    return this.findTrace(t => t.type == EventTypes.DEPLOYMENT_TRIGGERED && t.data.stage == stage)?.getFinishedEvent();
  }

  getRemediationActions(): Trace[] {
    // return remediation sequences
    return this.traces;
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
        if(this.getPendingApproval() != null)
          return "waiting for approval";
        else
          return "started";
    }
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
