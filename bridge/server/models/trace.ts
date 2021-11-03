import { Trace as ts } from '../../shared/models/trace';
import { ResultTypes } from '../../shared/models/result-types';

export class Trace extends ts {
  traces: Trace[] = [];

  static fromJSON(data: unknown): Trace {
    return Object.assign(new this(), data);
  }

  public static traceMapper(traces: Trace[]): Trace[] {
    traces = traces.map((trace) => Trace.fromJSON(trace));
    return ts.traceMapperGlobal(traces);
  }

  public getMessage(): string {
    let message = '';
    const finishedEvent = this.getFinishedEvent();
    if (finishedEvent?.data.message) {
      message = finishedEvent.data.message;
    } else {
      const failedEvent = this.findTrace((t) => t.data.result === ResultTypes.FAILED);
      let eventState;

      if (failedEvent) {
        if (!failedEvent.isFinished() && !failedEvent.isChanged()) {
          eventState = 'started';
        } else if (failedEvent.isChanged()) {
          eventState = 'changed';
        } else if (failedEvent.isFinished()) {
          eventState = `finished with result ${failedEvent.data.result}`;
        } else {
          eventState = '';
        }
        message = `${failedEvent.source} ${eventState}`;
      }
    }
    return message;
  }
}
