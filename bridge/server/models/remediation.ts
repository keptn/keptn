import { Sequence } from './sequence';
import { SequenceStage } from '../../shared/models/sequence';
import { EventState } from '../../shared/models/event-state';
import { ResultTypes } from '../../shared/models/result-types';

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
