import { Approval as ap } from '../../../shared/interfaces/approval';
import { Trace } from '../_models/trace';

export interface Approval extends ap {
  trace: Trace;
  evaluationTrace?: Trace;
}
