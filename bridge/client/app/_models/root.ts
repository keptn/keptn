import {Trace} from "./trace";

export class Root extends Trace {
  traces: Trace[];

  isFaulty(): string {
    let result: string = null;
    if(this.traces) {
      this.traces.forEach((trace) => {
        if(trace.isFaulty()) {
          result = trace.data.stage;
        }
      });
    }
    return result;
  }

  isWarning(): string {
    let result: string = null;
    if(this.traces) {
      this.traces.forEach((trace) => {
        if(trace.isWarning()) {
          result = trace.data.stage;
        }
      });
    }
    return result;
  }

  isSuccessful(): boolean {
    let result: boolean = false;
    return !this.isFaulty() && result;
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
