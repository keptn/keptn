import { ResultTypes } from '../models/result-types.js';

export interface EvaluationResult {
  result: ResultTypes;
  score: number;
}
