import { ISequence, SequenceState } from '../../shared/interfaces/sequence';
import { IServerSequenceStage } from '../interfaces/sequence-stage';

export class Sequence implements ISequence {
  name!: string;
  project!: string;
  service!: string;
  shkeptncontext!: string;
  state!: SequenceState;
  time!: string;
  stages!: IServerSequenceStage[];
  problemTitle?: string;

  public static fromJSON(data: unknown): Sequence {
    return Object.assign(new this(), data);
  }

  public reduceToStage(stageName: string): void {
    this.stages = this.stages.filter((stage) => stage.name === stageName);
  }
}
