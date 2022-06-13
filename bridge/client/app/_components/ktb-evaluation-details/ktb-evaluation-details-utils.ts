import { EvaluationResultTypeExtension, IDataPoint, IHeatmapTooltipType } from '../../_interfaces/heatmap';
import { ResultTypes } from '../../../../shared/models/result-types';
import { Trace } from '../../_models/trace';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';
import { parse as parseYaml } from 'yaml';
import { SloConfig } from '../../../../shared/interfaces/slo-config';
import { IEvaluationData } from '../../../../shared/models/trace';

export type SliInfo = {
  score: number;
  warningCount: number;
  failedCount: number;
  passCount: number;
};

export const filterUnparsedEvaluations = (trace: Trace): trace is Trace & { data: { evaluation: IEvaluationData } } =>
  !!trace?.data?.evaluation?.sloFileContent && !trace.data.evaluation.sloFileContentParsed;

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

export function getTotalScore(evaluation: Trace): number {
  return evaluation.data.evaluation?.indicatorResults
    ? evaluation.data.evaluation?.indicatorResults.reduce((total: number, ir: IndicatorResult) => total + ir.score, 0)
    : 1;
}

export function indicatorResultToDataPoint(
  evaluation: Trace,
  scoreValue: number
): (indicatorResult: IndicatorResult) => IDataPoint {
  return (indicatorResult: IndicatorResult): IDataPoint => {
    const totalScore: number = getTotalScore(evaluation);
    const metricScore = totalScore === 0 ? 0 : (indicatorResult.score / totalScore) * (scoreValue ?? 1);
    const color = indicatorResult.value.success ? indicatorResult.status : EvaluationResultTypeExtension.INFO;
    return {
      xElement: evaluation.getHeatmapLabel(),
      yElement: indicatorResult.displayName || indicatorResult.value.metric,
      color,
      identifier: evaluation.id,
      comparedIdentifier: evaluation.data.evaluation?.comparedEvents ?? [],
      tooltip: {
        type: IHeatmapTooltipType.SLI,
        value: indicatorResult.value.value,
        score: metricScore,
        keySli: indicatorResult.keySli,
        passTargets: indicatorResult.passTargets ?? [],
        warningTargets: indicatorResult.warningTargets ?? [],
      },
    };
  };
}

export function evaluationToDataPoint(evaluation: Trace, scoreValue: number): IDataPoint {
  const resultInfo = getSliResultInfo(evaluation.data.evaluation?.indicatorResults ?? []);
  return {
    xElement: evaluation.getHeatmapLabel(),
    yElement: 'Score',
    color: evaluation.data.evaluation?.result ?? EvaluationResultTypeExtension.INFO,
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
}

const addEvaluationToDataPoints = (points: IDataPoint[], evaluation: Trace): IDataPoint[] => {
  const scoreValue = evaluation.data.evaluation?.score ?? 0;
  const results: IDataPoint[] = evaluation.data.evaluation?.indicatorResults
    ? evaluation.data.evaluation?.indicatorResults.map(indicatorResultToDataPoint(evaluation, scoreValue))
    : [];
  const score: IDataPoint = evaluationToDataPoint(evaluation, scoreValue);
  return [...points, ...results, score];
};

export function createDataPoints(evaluationHistory: Trace[]): IDataPoint[] {
  const dataPoints = evaluationHistory.reduce(addEvaluationToDataPoints, []);
  const scores = dataPoints.filter((dp) => dp.tooltip.type === IHeatmapTooltipType.SCORE);
  const sortedResults = dataPoints
    .filter((dp) => dp.tooltip.type == IHeatmapTooltipType.SLI)
    .sort((a, b) => a.yElement.localeCompare(b.yElement));
  return [...scores, ...sortedResults];
}

export function parseSloOfEvaluations(evaluationTraces: Trace[]): void {
  evaluationTraces.filter(filterUnparsedEvaluations).forEach((e) => parseSloFile(e.data.evaluation));
}

function parseSloFile(evaluation: IEvaluationData): void {
  try {
    const sloFileContentParsed: SloConfig = parseYaml(atob(evaluation.sloFileContent));
    evaluation.score_pass = sloFileContentParsed.total_score?.pass?.split('%')[0];
    evaluation.score_warning = sloFileContentParsed.total_score?.warning?.split('%')[0] ?? '';
    evaluation.compare_with = sloFileContentParsed.comparison.compare_with ?? '';
    evaluation.include_result_with_score = sloFileContentParsed.comparison.include_result_with_score;
    evaluation.sloObjectives = sloFileContentParsed.objectives;
    evaluation.sloFileContentParsed = true;
  } catch {}
}
