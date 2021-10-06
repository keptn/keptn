import { Trace } from '../_models/trace';

export interface EvaluationHistory {
  type: string;
  triggerEvent: Trace;
  traces?: Trace[];
}
