import {Trace} from "./trace";
import {EventTypes} from "./event-types";
import {Stage} from "./stage";

export class Sequence extends Trace {
  traces: Trace[];

  isFaulty(): string {
    return this.traces.reduce((result: string, trace: Trace) => trace.isFaulty() ? trace.data.stage : result, null);
  }

  isProblem(): boolean {
    return this.traces.reduce((result: boolean, trace: Trace) => trace.isProblem() && !trace.isProblemResolvedOrClosed() ? true : result, false);
  }

  isProblemResolvedOrClosed(): boolean {
    return this.traces.reduce((result: boolean, trace: Trace) => trace.isProblem() && trace.isProblemResolvedOrClosed() ? true : result, false);
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

  isApproval(): string {
    return this.getLastTrace().isApproval();
  }

  getLastTrace(): Trace {
    return this.traces ? this.traces[this.traces.length - 1] : null;
  }

  getStages(): String[] {
    let result: String[] = [];
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

  getEvaluation(stage: Stage): Trace {
    return this.traces.find(t => t.type == EventTypes.EVALUATION_DONE && t.data.stage == stage.stageName);
  }

  getDeploymentDetails(stage: Stage): Trace {
    return this.traces.find(t => t.type == EventTypes.DEPLOYMENT_FINISHED && t.data.stage == stage.stageName);
  }

  getRemediationActions(): Sequence[] {
    // create chunks of Remediations and start new chunk at REMEDIATION_TRIGGERED event
    return this.traces.reduce((result, trace: Trace) => {
      if(trace.type == EventTypes.ACTION_TRIGGERED)
        result.push(Sequence.fromJSON(JSON.parse(JSON.stringify(trace))));
      else if(result.length)
        result[result.length-1].traces = [...result[result.length-1].traces||[], trace];
      return result;
    }, []);
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
