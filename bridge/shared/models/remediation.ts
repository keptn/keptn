import { ISequenceState, SequenceStage, SequenceStatus } from '../interfaces/sequence';
import { IRemediationAction } from './remediation-action';

export class Remediation implements ISequenceState {
  name!: string;
  project!: string;
  service!: string;
  shkeptncontext!: string;
  state!: SequenceStatus;
  time!: string;

  stages: (SequenceStage & {
    actions: IRemediationAction[];
  })[] = [];

  public static fromJSON(data: unknown): Remediation {
    return Object.assign(new this(), data);
  }
}
