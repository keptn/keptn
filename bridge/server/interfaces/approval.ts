import { Approval as ap } from '../../shared/interfaces/approval';
import { Trace } from '../models/trace';

export interface Approval extends ap {
  trace: Trace;
  evaluationTrace?: Trace;
}
