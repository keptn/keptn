import { Trace } from '../_models/trace';

export interface EvaluationHistory {
  type: 'evaluationHistory' | 'invalidateEvaluation';
  triggerEvent: Trace;
  traces?: Trace[];
}
