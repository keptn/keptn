import { SequenceStage } from '../../shared/interfaces/sequence';
import { Trace } from '../../shared/models/trace';

export interface IServerSequenceStage extends SequenceStage {
  latestEvaluationTrace?: Trace;
}
