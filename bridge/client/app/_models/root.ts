import {Trace} from "./trace";

export class Root extends Trace {
  traces: Trace[];

  isFaulty(): string {
    return this.traces.reduce((result: string, trace: Trace) => trace.isFaulty() ? trace.data.stage : result, null);
  }

  isProblem(): boolean {
    return this.traces.reduce((result: boolean, trace: Trace) => trace.isProblem() && !trace.isProblemResolvedOrClosed() ? true : result, false);
  }

  isWarning(): string {
    return this.traces.reduce((result: string, trace: Trace) => trace.isWarning() ? trace.data.stage : result, null);
  }

  isSuccessful(): string {
    return this.traces.reduce((result: string, trace: Trace) => trace.isSuccessful() ? trace.data.stage : result, null);
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

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
