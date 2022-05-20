import {
  EvaluationResultType,
  EvaluationResultTypeExtension,
  IDataPoint,
  IHeatmapTooltipType,
} from '../../_interfaces/heatmap';
import { ResultTypes } from '../../../../shared/models/result-types';
import { Trace } from '../../_models/trace';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';

export type SliInfo = {
  score: number;
  warningCount: number;
  failedCount: number;
  passCount: number;
};

const _resultType: { [key: string]: EvaluationResultType } = {
  pass: ResultTypes.PASSED,
  warning: ResultTypes.WARNING,
  fail: ResultTypes.FAILED,
  failed: ResultTypes.FAILED,
  info: EvaluationResultTypeExtension.INFO,
};

export function getSliResultInfo(indicatorResults: IndicatorResult[]): {
  score: number;
  warningCount: number;
  failedCount: number;
  passCount: number;
} {
  return indicatorResults.reduce(
    (acc, result) => {
      const warning = result.status === ResultTypes.WARNING ? 1 : 0;
      const failed = result.status === ResultTypes.FAILED ? 1 : 0;
      return {
        score: acc.score + result.score,
        warningCount: acc.warningCount + warning,
        failedCount: acc.failedCount + failed,
        passCount: acc.passCount + 1 - warning - failed,
      };
    },
    { score: 0, warningCount: 0, failedCount: 0, passCount: 0 } as SliInfo
  );
}

const getSlisAndScore = (points: IDataPoint[], evaluation: Trace): IDataPoint[] => {
  const totalScore: number = evaluation.data.evaluation?.indicatorResults
    ? evaluation.data.evaluation?.indicatorResults.reduce((total: number, ir: IndicatorResult) => total + ir.score, 0)
    : 1;
  const scoreValue = evaluation.data.evaluation?.score ?? 0;

  const indicatorResultToDataPoint = (indicatorResult: IndicatorResult): IDataPoint => {
    const color = indicatorResult.value.success
      ? _resultType[indicatorResult.status]
      : EvaluationResultTypeExtension.INFO;
    return {
      xElement: evaluation.getHeatmapLabel(),
      yElement: indicatorResult.value.metric,
      color,
      identifier: evaluation.id,
      comparedIdentifier: [],
      tooltip: {
        type: IHeatmapTooltipType.SLI,
        value: indicatorResult.value.value,
        score: (indicatorResult.score / totalScore) * (scoreValue ?? 1),
        keySli: indicatorResult.keySli,
        passTargets: indicatorResult.passTargets ?? [],
        warningTargets: indicatorResult.warningTargets ?? [],
      },
    };
  };

  const slis: IDataPoint[] = evaluation.data.evaluation?.indicatorResults
    ? evaluation.data.evaluation?.indicatorResults
        .map(indicatorResultToDataPoint)
        .sort((a, b) => b.yElement.localeCompare(a.yElement))
    : [];
  const resultInfo = getSliResultInfo(evaluation.data.evaluation?.indicatorResults ?? []);
  const score: IDataPoint = {
    xElement: evaluation.getHeatmapLabel(),
    yElement: 'Score',
    color: _resultType[evaluation.data.evaluation?.result ?? EvaluationResultTypeExtension.INFO],
    identifier: evaluation.id,
    comparedIdentifier: evaluation.data.evaluation?.comparedEvents ?? [],
    tooltip: {
      type: IHeatmapTooltipType.SCORE,
      value: scoreValue,
      passCount: resultInfo.passCount,
      warningCount: resultInfo.warningCount,
      failedCount: resultInfo.failedCount,
      thresholdPass: +(evaluation.data.evaluation?.score_pass ?? 0),
      thresholdWarn: +(evaluation.data.evaluation?.score_warning ?? 0),
      fail: evaluation.isFailed(),
      warn: evaluation.isWarning(),
    },
  };

  return [...points, ...slis, score];
};

export function createDataPoints(evaluationHistory: Trace[]): IDataPoint[] {
  const slisAndScores = evaluationHistory.reduce(getSlisAndScore, []);
  const scores = slisAndScores
    .filter((v) => v.tooltip.type === IHeatmapTooltipType.SCORE)
    .sort((a, b) => a.xElement.localeCompare(b.xElement));
  const slis = slisAndScores.filter((v) => v.tooltip.type === IHeatmapTooltipType.SLI);
  return [...slis, ...scores];
}
