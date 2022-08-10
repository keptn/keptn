import { Trace } from '../../_models/trace';
import {
  ChartItem,
  ChartItemPoint,
  ChartItemPointInfo,
  DrawType,
  FuncDateToDict,
  FuncDateToString,
  FuncEvaluationToChartItemPoint,
  FuncMapIndicatorResult,
  FuncMetricToChartItem,
} from '../../_interfaces/chart';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';

export function createChartPoints(evaluationHistory: Trace[]): ChartItem[] {
  const mapEvaluationToScoreCharItemPoint =
    (includeColor: boolean): FuncEvaluationToChartItemPoint =>
    (evaluation: Trace, index: number): ChartItemPoint => ({
      x: index,
      y: evaluation.data.evaluation?.score ?? 0,
      color: includeColor ? getScoreColor(evaluation) : undefined,
    });
  const mapIndicatorResultChartItemPointsToChartItem =
    (chartInfoDict: ChartItemPointInfo): FuncMetricToChartItem =>
    (metric: string): ChartItem => ({
      label: chartInfoDict[metric]?.label || metric,
      points: chartInfoDict[metric]?.points ?? [],
      type: 'metric-line',
      invisible: true,
    });
  const mapChartItemPointsToChartItem = (chartPoints: ChartItemPoint[], type: DrawType): ChartItem => ({
    label: 'score',
    points: chartPoints,
    type,
  });
  const addIndicatorResultChartItemToDict =
    (chartPoints: ChartItemPointInfo, index: number): FuncMapIndicatorResult =>
    (indicatorResult: IndicatorResult): void => {
      const metricChartItemPoints = (chartPoints[indicatorResult.value.metric] ??= { points: [] });
      metricChartItemPoints.points.push({
        x: index,
        y: indicatorResult.value.value,
      });
      metricChartItemPoints.label ||= indicatorResult.displayName;
    };

  const reduceEvaluationToIndicatorResultChartPoint = (
    chartPoints: ChartItemPointInfo,
    evaluation: Trace,
    index: number
  ): ChartItemPointInfo => {
    evaluation.data.evaluation?.indicatorResults?.forEach(addIndicatorResultChartItemToDict(chartPoints, index));
    return chartPoints;
  };

  const scoreBarChartPoints: ChartItemPoint[] = evaluationHistory.map(mapEvaluationToScoreCharItemPoint(true));
  const scoreLineChartPoints: ChartItemPoint[] = evaluationHistory.map(mapEvaluationToScoreCharItemPoint(false));
  const indicatorResultChartPoints = evaluationHistory.reduce(
    reduceEvaluationToIndicatorResultChartPoint,
    {} as ChartItemPointInfo
  );

  const chartScoreBarItems = mapChartItemPointsToChartItem(scoreBarChartPoints, 'score-bar');
  const chartScoreLineItems = mapChartItemPointsToChartItem(scoreLineChartPoints, 'score-line');
  const chartMetricItems: ChartItem[] = Object.keys(indicatorResultChartPoints).map(
    mapIndicatorResultChartItemPointsToChartItem(indicatorResultChartPoints)
  );

  return [chartScoreBarItems, chartScoreLineItems, ...chartMetricItems];
}

export function createChartXLabels(evaluationHistory: Trace[]): Record<number, string> {
  const mapEvaluationToLabel = (evaluation: Trace): string => evaluation.getChartLabel();
  const mapDatesToDuplicateCounterAndSetOccurrences =
    (dict: Record<string, number | undefined>): FuncDateToDict =>
    (date): number | undefined => {
      dict[date] = (dict[date] ?? 0) + 1;
      return dict[date];
    };
  const mapDateToUniqueDate =
    (dict: Record<string, number | undefined>, counter: (number | undefined)[]): FuncDateToString =>
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
