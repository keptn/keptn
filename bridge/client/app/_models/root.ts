import { Trace } from './trace';

export class Root extends Trace {
  traces: Trace[] = [];

  static fromJSON(data: unknown): Root {
    return Object.assign(new this(), data);
  }

  isFaulty(stageName?: string): boolean {
    // a Sequence is faulty, if there is a sequence that is faulty, but no other sequence that is successful on the same stage
    let result = false;
    const stages = stageName ? [stageName] : this.getStages();
    if (stages.length > 0) {
      stages.forEach((stage) => {
        result ||= this.areTracesFaulty(stage);
      });
    } else {
      result = this.areTracesFaulty();
    }
    return result;
  }

  private areTracesFaulty(stageName?: string): boolean {
    const stageTraces = stageName ? this.traces.filter((t) => t.stage === stageName) : this.traces;
    return stageTraces.some((t) => t.isFaulty()) && !stageTraces.some((t) => t.isSuccessful());
  }

  getStages(): string[] {
    const result: string[] = [];
    if (this.traces) {
      this.traces.forEach((trace) => {
        if (trace.data.stage && result.indexOf(trace.data.stage) === -1) {
          result.push(trace.data.stage);
        }
      });
    }
    return result;
  }
}
