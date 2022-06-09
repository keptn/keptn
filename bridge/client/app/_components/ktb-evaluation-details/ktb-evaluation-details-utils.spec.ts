import { Trace } from '../../_models/trace';
import {
  createDataPoints,
  evaluationToDataPoint,
  getTotalScore,
  indicatorResultToDataPoint,
  parseSloOfEvaluations,
} from './ktb-evaluation-details-utils';
import { EvaluationsMock } from '../../_services/_mockData/evaluations.mock';
import { IDataPoint, IHeatmapScoreTooltip, IHeatmapSliTooltip, IHeatmapTooltipType } from '../../_interfaces/heatmap';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';
import { EventTypes } from '../../../../shared/interfaces/event-types';

describe('KtbEvaluationDetailsUtils', () => {
  const validSLOFile =
    'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg';

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

  it('should parse all SLO files even if there is an invalid one', () => {
    // given
    const trace = getTraceWithSloContent('_myInvalidSLOFile_');
    const trace2 = getTraceWithSloContent(validSLOFile);
    const trace3 = getTraceWithSloContent(undefined);

    // when
    parseSloOfEvaluations([trace, trace2, trace3]);

    // then
    expect(trace.data.evaluation?.sloFileContentParsed).toBeUndefined();
    expect(trace2.data.evaluation?.sloFileContentParsed).toEqual({
      comparison: {
        aggregate_function: 'avg',
        compare_with: 'single_result',
        include_result_with_score: 'pass',
        number_of_comparison_results: 1,
      },
      filter: null,
      objectives: [
        {
          displayName: 'Response time P95',
          key_sli: false,
          pass: [
            {
              criteria: ['<=+10%', '<600'],
            },
          ],
          sli: 'response_time_p95',
          warning: [
            {
              criteria: ['<=800'],
            },
          ],
          weight: 1,
        },
      ],
      spec_version: '1.0',
      total_score: {
        pass: '90%',
        warning: '75%',
      },
    });
    expect(trace3.data.evaluation?.sloFileContentParsed).toBe(false);
  });

  it('should correctly set values after parsing the SLO file', () => {
    // given
    const trace = getTraceWithSloContent(validSLOFile, true);

    // when
    parseSloOfEvaluations([trace]);

    // then
    expect(trace.data.evaluation?.score_pass).toBe('90');
    expect(trace.data.evaluation?.score_warning).toBe('75');
    expect(trace.data.evaluation?.compare_with).toBe('single_result');
    expect(trace.data.evaluation?.include_result_with_score).toBe('pass');
    expect(trace.data.evaluation?.sloFileContentParsed).toBe(true);
    expect(trace.data.evaluation?.sloObjectives).toEqual([
      {
        displayName: 'Response time P95',
        key_sli: false,
        pass: [
          {
            criteria: ['<=+10%', '<600'],
          },
        ],
        sli: 'response_time_p95',
        warning: [
          {
            criteria: ['<=800'],
          },
        ],
        weight: 1,
      },
    ]);
  });

  function getTraceWithSloContent(sloContent: string | undefined, comparedEvents = false): Trace {
    return Trace.fromJSON({
      data: {
        evaluation: {
          sloFileContent: sloContent,
          ...(comparedEvents && { comparedEvents: ['myOtherId'] }),
        },
      },
      id: 'myId',
      type: EventTypes.EVALUATION_FINISHED,
      shkeptncontext: 'myContext',
      time: '2022-06-02T09:38:20.855Z',
    });
  }
});
