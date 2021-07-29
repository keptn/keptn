import { Sequence as sq, SequenceStage } from '../../shared/models/sequence';
import { Trace } from '../../shared/models/trace';

export class Sequence extends sq {
  stages!: (SequenceStage &
    {
      latestEvaluationTrace?: Trace
    })[]
  ;
  problemTitle?: string;

  public static fromJSON(data: unknown): Sequence {
    return Object.assign(new this(), data);
  }

  public reduceToStage(stageName: string): void {
    this.stages = this.stages.filter(stage => stage.name === stageName);
  }
}
