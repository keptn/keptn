import { ISequence, SequenceStage, SequenceState } from '../../shared/interfaces/sequence';
import { Trace } from '../../shared/models/trace';

export class Sequence implements ISequence {
  name!: string;
  project!: string;
  service!: string;
  shkeptncontext!: string;
  state!: SequenceState;
  time!: string;
  stages!: (SequenceStage & {
    latestEvaluationTrace?: Trace;
  })[];
  problemTitle?: string;

  public static fromJSON(data: unknown): Sequence {
    return Object.assign(new this(), data);
  }

  public reduceToStage(stageName: string): void {
    this.stages = this.stages.filter((stage) => stage.name === stageName);
  }
}
