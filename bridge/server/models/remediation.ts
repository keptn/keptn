import { Sequence } from './sequence';
import { SequenceStage } from '../../shared/models/sequence';
import { IRemediationAction } from '../../shared/models/remediation-action';

export class Remediation extends Sequence {
  stages: (SequenceStage & {
    actions: IRemediationAction[]
  })[] = [];

  public static fromJSON(data: unknown): Remediation {
    return Object.assign(new this(), data);
  }
}
