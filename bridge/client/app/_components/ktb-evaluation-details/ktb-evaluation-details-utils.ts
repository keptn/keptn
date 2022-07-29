import { EvaluationResultTypeExtension, IDataPoint, IHeatmapTooltipType } from '../../_interfaces/heatmap';
import { ResultTypes } from '../../../../shared/models/result-types';
import { Trace } from '../../_models/trace';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';
import { parse as parseYaml } from 'yaml';
import { SloConfig } from '../../../../shared/interfaces/slo-config';
import { IEvaluationData } from '../../../../shared/models/trace';
import { ChartItem, ChartItemPoint, IChartItemPointInfo } from '../../_interfaces/chart';

export type SliInfo = {
  score: number;
  warningCount: number;
  failedCount: number;
  passCount: number;
  keySliCount: number;
  keySliFailedCount: number;
};

export const filterUnparsedEvaluations = (trace: Trace): trace is Trace & { data: { evaluation: IEvaluationData } } =>
  !!trace?.data?.evaluation?.sloFileContent && !trace.data.evaluation.sloFileContentParsed;

export function getSliResultInfo(indicatorResults: IndicatorResult[]): SliInfo {
  return indicatorResults.reduce(
    (acc, result) => {
      const warning = result.status === ResultTypes.WARNING ? 1 : 0;
      const failed = result.status === ResultTypes.FAILED ? 1 : 0;
      const keySli = result.keySli ? 1 : 0;
      const keySliFailed = result.keySli && result.status === ResultTypes.FAILED ? 1 : 0;
      return {
        score: acc.score + result.score,
        warningCount: acc.warningCount + warning,
        failedCount: acc.failedCount + failed,
        passCount: acc.passCount + 1 - warning - failed,
        keySliCount: acc.keySliCount + keySli,
        keySliFailedCount: acc.keySliFailedCount + keySliFailed,
      };
    },
    { score: 0, warningCount: 0, failedCount: 0, passCount: 0, keySliCount: 0, keySliFailedCount: 0 } as SliInfo
  );
}

export function getTotalScore(evaluation: Trace): number {
  return evaluation.data.evaluation?.indicatorResults
    ? evaluation.data.evaluation?.indicatorResults.reduce((total: number, ir: IndicatorResult) => total + ir.score, 0)
    : 1;
}

export function getScoreState(evaluation: IEvaluationData): string {
  if (evaluation.score_pass && evaluation.score >= evaluation.score_pass) {
    return 'pass';
  } else if (evaluation.score_warning && evaluation.score >= evaluation.score_warning) {
    return 'warning';
  } else {
    return 'fail';
  }
}

export function getScoreInfo(evaluation: IEvaluationData): string {
  if (evaluation.score_state === 'pass') {
    return ` >= ${evaluation.score_pass}`;
  } else if (evaluation.score_state === 'warning') {
    return ` >= ${evaluation.score_warning}`;
  } else {
    return ` < ${evaluation.score_warning || evaluation.score_pass}`;
  }
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

export function evaluationToScoreDataPoint(evaluation: Trace, scoreValue: number): IDataPoint {
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
      thresholdPass: evaluation.data.evaluation?.score_pass ?? 0,
      thresholdWarn: evaluation.data.evaluation?.score_warning ?? 0,
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
  const score: IDataPoint = evaluationToScoreDataPoint(evaluation, scoreValue);
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
    evaluation.score_pass = +(sloFileContentParsed.total_score?.pass?.split('%')[0] ?? 0);
    evaluation.score_warning = +(sloFileContentParsed.total_score?.warning?.split('%')[0] ?? 0);
    evaluation.score_state = getScoreState(evaluation);
    evaluation.score_info = getScoreInfo(evaluation);
    evaluation.compare_with = sloFileContentParsed.comparison.compare_with ?? '';
    evaluation.include_result_with_score = sloFileContentParsed.comparison.include_result_with_score;
    evaluation.sloObjectives = sloFileContentParsed.objectives;
    evaluation.sloFileContentParsed = true;
  } catch {}
}

export function createChartPoints(evaluationHistory: Trace[]): ChartItem[] {
  const mapEvaluationToScoreCharItemPoint = (evaluation: Trace, index: number): ChartItemPoint => ({
    x: index,
    y: evaluation.data.evaluation?.score ?? 0,
    color: getScoreColor(evaluation),
    identifier: evaluation.id,
  });
  const mapIndicatorResultChartItemPointsToChartItem =
    (chartInfoDict: IChartItemPointInfo): ((metric: string) => ChartItem) =>
    (metric: string): ChartItem => ({
      label: chartInfoDict[metric]?.label || metric,
      points: chartInfoDict[metric]?.points ?? [],
      type: 'metric-line',
      identifier: metric,
      invisible: true,
    });
  const mapChartItemPointsToChartItem =
    (chartPoints: ChartItemPoint[]): ((type: ChartItem['type']) => ChartItem) =>
    (type: ChartItem['type']): ChartItem => ({
      label: 'score',
      points: chartPoints,
      identifier: 'score',
      type,
    });
  const addIndicatorResultChartItemToDict =
    (
      chartPoints: IChartItemPointInfo,
      identifier: string,
      index: number
    ): ((indicatorResult: IndicatorResult) => void) =>
    (indicatorResult: IndicatorResult): void => {
      const metricChartItemPoints = (chartPoints[indicatorResult.value.metric] ??= { points: [] });
      metricChartItemPoints.points.push({
        x: index,
        y: indicatorResult.value.value,
        identifier,
      });
      metricChartItemPoints.label ||= indicatorResult.displayName;
    };

  const reduceEvaluationToIndicatorResultChartPoint = (
    chartPoints: IChartItemPointInfo,
    evaluation: Trace,
    index: number
  ): IChartItemPointInfo => {
    evaluation.data.evaluation?.indicatorResults?.forEach(
      addIndicatorResultChartItemToDict(chartPoints, evaluation.id, index)
    );
    return chartPoints;
  };

  const scoreChartPoints: ChartItemPoint[] = evaluationHistory.map(mapEvaluationToScoreCharItemPoint);
  const indicatorResultChartPoints = evaluationHistory.reduce(
    reduceEvaluationToIndicatorResultChartPoint,
    {} as IChartItemPointInfo
  );

  const chartScoreItems = (['score-bar', 'score-line'] as ChartItem['type'][]).map(
    mapChartItemPointsToChartItem(scoreChartPoints)
  );
  const chartMetricItems: ChartItem[] = Object.keys(indicatorResultChartPoints).map(
    mapIndicatorResultChartItemPointsToChartItem(indicatorResultChartPoints)
  );

  return [...chartScoreItems, ...chartMetricItems];
}

export function createChartXLabels(evaluationHistory: Trace[]): Record<number, string> {
  const mapEvaluationToLabel = (evaluation: Trace): string => evaluation.getHeatmapLabel();
  const mapDatesToDuplicateCounterAndSetOccurrences =
    (dict: Record<string, number | undefined>): ((date: string) => number | undefined) =>
    (date) => {
      dict[date] = (dict[date] ?? 0) + 1;
      return dict[date];
    };
  const mapDateToUniqueDate =
    (
      dict: Record<string, number | undefined>,
      counter: (number | undefined)[]
    ): ((date: string, index: number) => string) =>
    (date: string, index: number): string =>
      dict[date] === 1 ? date : `${date} (${counter[index]})`;
  const reduceArrayToDict = (labels: Record<number, string>, date: string, index: number): Record<number, string> => {
    labels[index] = date;
    return labels;
  };
  const occurrencesDict: Record<string, number | undefined> = {};
  const dates = evaluationHistory.map(mapEvaluationToLabel);
  const duplicateCounter = dates.map(mapDatesToDuplicateCounterAndSetOccurrences(occurrencesDict));

  return dates
    .map(mapDateToUniqueDate(occurrencesDict, duplicateCounter))
    .reduce(reduceArrayToDict, {} as Record<number, string>);
}

export function createChartTooltipLabels(
  evaluationHistory: Trace[],
  formatTime: (date: string) => string
): Record<number, string> {
  const mapEvaluationToTime = (evaluation: Trace): string | undefined => evaluation.time;
  const reduceArrayToDict = (
    tooltipLabels: Record<number, string>,
    time: string | undefined,
    index: number
  ): Record<number, string> => {
    tooltipLabels[index] = time ? `SLO evaluation of test from ${formatTime(time)}` : '';
    return tooltipLabels;
  };

  return evaluationHistory.map(mapEvaluationToTime).reduce(reduceArrayToDict, {} as Record<number, string>);
}

/***
 * Maps the evaluation result to a color (pass, warning, error)
 * @param evaluation
 */
function getScoreColor(evaluation: Trace): string {
  if (evaluation.isWarning()) {
    return '#e6be00';
  }
  if (evaluation.isFaulty()) {
    return '#dc172a';
  }
  return '#7dc540';
}
