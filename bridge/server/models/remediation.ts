import { Sequence, SequenceStage } from './sequence.js';
import { ResultTypes } from './result-types.js';
import { EventState } from './event-state';

export class Remediation extends Sequence {
  stages: (SequenceStage & {
    actions: {
      action: string;
      description: string;
      name: string;
      state: EventState,
      result?: ResultTypes
    }[]
  })[] = [];

  public static fromJSON(data: unknown): Remediation {
    return Object.assign(new this(), data);
  }
}
