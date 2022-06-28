import { ISequence, SequenceStage, SequenceState } from '../interfaces/sequence';
import { IRemediationAction } from './remediation-action';

export class Remediation implements ISequence {
  name!: string;
  project!: string;
  service!: string;
  shkeptncontext!: string;
  state!: SequenceState;
  time!: string;

  stages: (SequenceStage & {
    actions: IRemediationAction[];
  })[] = [];

  public static fromJSON(data: unknown): Remediation {
    return Object.assign(new this(), data);
  }
}
