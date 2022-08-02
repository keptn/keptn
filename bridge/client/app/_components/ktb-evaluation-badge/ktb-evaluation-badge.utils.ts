import { Trace } from '../../_models/trace';
import { TruncateNumberPipe } from '../../_pipes/truncate-number';
import { EvaluationResult } from '../../../../shared/interfaces/evaluation-result';
import { ResultTypes } from '../../../../shared/models/result-types';

export enum EvaluationBadgeVariant {
  BORDER,
  FILL,
  NONE,
}

export interface IEvaluationBadgeState {
  isError: boolean;
  isWarning: boolean;
  isSuccess: boolean;
  fillState: EvaluationBadgeVariant;
  score?: number;
}

const truncateNumberPipe = new TruncateNumberPipe();
const truncScoreDecimals = 0;

export function getEvaluationBadgeState(evaluation: Trace, fillState: EvaluationBadgeVariant): IEvaluationBadgeState {
  return {
    isError: evaluation.isFaulty(),
    isWarning: evaluation.isWarning(),
    isSuccess: evaluation.isSuccessful(),
    score: truncateNumberPipe.transform(
      evaluation.getEvaluationFinishedEvent()?.data?.evaluation?.score,
      truncScoreDecimals
    ),
    fillState,
  };
}

export function getEvaluationResultBadgeState(
  evaluationResult: EvaluationResult | undefined,
  fillState: EvaluationBadgeVariant
): IEvaluationBadgeState {
  return {
    isError: evaluationResult?.result === ResultTypes.FAILED,
    isWarning: evaluationResult?.result === ResultTypes.WARNING,
    isSuccess: evaluationResult?.result === ResultTypes.PASSED,
    score: truncateNumberPipe.transform(evaluationResult?.score, truncScoreDecimals),
    fillState,
  };
}
