import { Trace } from '../models/trace';

export interface Approval {
  trace: Trace;
  evaluationTrace?: Trace;
}
