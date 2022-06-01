import { Trace } from '../../_models/trace';
import {
  createDataPoints,
  evaluationToDataPoint,
  getTotalScore,
  indicatorResultToDataPoint,
} from './ktb-evaluation-details-utils';
import { EvaluationsMock } from '../../_services/_mockData/evaluations.mock';
import { IDataPoint, IHeatmapScoreTooltip, IHeatmapSliTooltip, IHeatmapTooltipType } from '../../_interfaces/heatmap';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';

describe('KtbEvaluationDetailsUtils', () => {
  it('should calculate the right total score', () => {
    // given
    const traces = EvaluationsMock.data.evaluationHistory as Trace[];

    // when
    const totalScores = traces.map((t) => getTotalScore(t));

    // then
    expect(totalScores).toStrictEqual([7, 0.5, 1, 1, 1, 1, 0, 1, 0, 1, 1, 1]);
  });
  it('should transform an evaluation score to a data point', () => {
    // given
    const traces = EvaluationsMock.data.evaluationHistory as Trace[];
    const evaluation = traces[0];
    const scoreValue = evaluation.data.evaluation?.score ?? 0;

    // when
    const dataPoint = evaluationToDataPoint(evaluation, scoreValue);

    // then
    expect(dataPoint.xElement).toBe('2020-11-10 12:12');
    expect(dataPoint.yElement).toBe('Score');
    expect(dataPoint.identifier).toBe('04266cc2-eeea-485b-85b3-f1dea50890ce');
    expect(dataPoint.color).toBe('pass');
    expect(dataPoint.comparedIdentifier).toStrictEqual(['cfa408ce-f552-43c4-aff2-673b1e0548d2']);
    expect(dataPoint.tooltip.type).toBe(IHeatmapTooltipType.SCORE);
    const tooltip = dataPoint.tooltip as IHeatmapScoreTooltip;
    expect(tooltip.value).toBe(scoreValue);
    expect(tooltip.passCount).toBe(3);
    expect(tooltip.failedCount).toBe(0);
    expect(tooltip.warningCount).toBe(0);
    expect(tooltip.thresholdPass).toBe(0);
    expect(tooltip.thresholdWarn).toBe(0);
    expect(tooltip.fail).toBe(false);
    expect(tooltip.warn).toBe(false);
  });
  it('should transform an indicator result to a data point', () => {
    // given
    const traces = EvaluationsMock.data.evaluationHistory as Trace[];
    const evaluation = traces[0];
    const scoreValue = evaluation.data.evaluation?.score ?? 0;

    // when
    const mapper = indicatorResultToDataPoint(evaluation, scoreValue);
    const dataPoint = mapper(<IndicatorResult>evaluation.data.evaluation?.indicatorResults[0]);

    // then
    expect(dataPoint.xElement).toBe('2020-11-10 12:12');
    expect(dataPoint.yElement).toBe('response_time_p95');
    expect(dataPoint.identifier).toBe('04266cc2-eeea-485b-85b3-f1dea50890ce');
    expect(dataPoint.color).toBe('pass');
    expect(dataPoint.comparedIdentifier).toStrictEqual(['cfa408ce-f552-43c4-aff2-673b1e0548d2']);
    expect(dataPoint.tooltip.type).toBe(IHeatmapTooltipType.SLI);
    const tooltip = dataPoint.tooltip as IHeatmapSliTooltip;
    expect(tooltip.value).toBe(299.18637492576534);
    expect(tooltip.keySli).toBe(false);
    expect(tooltip.score).toBe(14.285714285714285);
    expect(tooltip.warningTargets).toStrictEqual([]);
    expect(tooltip.passTargets).toStrictEqual([]);
  });
  it('should create data points from an evaluation history', () => {
    // given
    const traces = EvaluationsMock.data.evaluationHistory as Trace[];

    // when
    const dataPoints: IDataPoint[] = createDataPoints(traces);

    // then
    expect(dataPoints.length).toBe(26);

    // Scores are before the SLIs
    // and sorting of them is NOT changed
    traces.forEach((value, index) => {
      const dataPointScore = dataPoints[index];
      expect(dataPointScore.tooltip.type).toBe(IHeatmapTooltipType.SCORE);
      expect(dataPointScore.identifier).toBe(value.id);
    });

    // Check if metrics are sorted
    const startOfMetricsIndex = traces.length;
    for (let i = startOfMetricsIndex + 1; i < dataPoints.length; i++) {
      const prev = dataPoints[i - 1];
      const cur = dataPoints[i];
      expect(prev.tooltip.type).toBe(IHeatmapTooltipType.SLI);
      expect(cur.tooltip.type).toBe(IHeatmapTooltipType.SLI);
      expect(prev.yElement.localeCompare(cur.yElement)).toBeLessThanOrEqual(0);
    }
  });
});
