import { Trace } from './trace.js';
import { EvaluationResult } from '../interfaces/evaluation-result.js';

export type SequenceEvent = {
  id: string,
  time: string,
  type: string
};

export type SequenceStage = {
  image?: string,
  latestEvaluation?: EvaluationResult,
  latestEvent?: SequenceEvent,
  latestFailedEvent?: SequenceEvent,
  name: string,
  latestEvaluationTrace?: Trace;
};

export class Sequence {
  name!: string;
  project!: string;
  service!: string;
  shkeptncontext!: string;
  stages!: SequenceStage[];
  state!: 'triggered' | 'finished' | 'waiting';
  time!: Date;
  problemTitle?: string;

  public static fromJSON(data: unknown): Sequence {
    return Object.assign(new this(), data);
  }

  public reduceToStage(stageName: string): void {
    this.stages = this.stages.filter(stage => stage.name === stageName);
  }
}
