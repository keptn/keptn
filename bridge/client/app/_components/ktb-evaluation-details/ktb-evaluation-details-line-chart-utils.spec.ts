import { EvaluationsMock } from '../../_services/_mockData/evaluations.mock';
import {
  createChartPoints,
  createChartTooltipLabels,
  createChartXLabels,
} from './ktb-evaluation-detials-line-chart-utils';
import { Trace } from '../../_models/trace';

describe('Line chart utils', () => {
  describe('createChartPoints', () => {
    it('should correctly map evaluation to chartPoints', () => {
      // given
      const data = EvaluationsMock.data.evaluationHistory as Trace[];

      // when
      const points = createChartPoints([data[0]]);

      // then
      expect(points).toEqual([
        {
          identifier: 'score',
          label: 'score',
          points: [
            {
              color: '#7dc540',
              identifier: '04266cc2-eeea-485b-85b3-f1dea50890ce',
              x: 0,
              y: 100,
            },
          ],
          type: 'score-bar',
        },
        {
          identifier: 'score',
          label: 'score',
          points: [
            {
              color: undefined,
              identifier: '04266cc2-eeea-485b-85b3-f1dea50890ce',
              x: 0,
              y: 100,
            },
          ],
          type: 'score-line',
        },
        {
          identifier: 'response_time_p95',
          invisible: true,
          label: 'response_time_p95',
          points: [
            {
              identifier: '04266cc2-eeea-485b-85b3-f1dea50890ce',
              x: 0,
              y: 299.18637492576534,
            },
          ],
          type: 'metric-line',
        },
        {
          identifier: 'response_time_p90',
          invisible: true,
          label: 'response_time_p90',
          points: [
            {
              identifier: '04266cc2-eeea-485b-85b3-f1dea50890ce',
              x: 0,
              y: 250.18637492576534,
            },
          ],
          type: 'metric-line',
        },
        {
          identifier: 'response_time_p50',
          invisible: true,
          label: 'response_time_p50',
          points: [
            {
              identifier: '04266cc2-eeea-485b-85b3-f1dea50890ce',
              x: 0,
              y: 200.18637492576534,
            },
          ],
          type: 'metric-line',
        },
      ]);
    });

    it('should favor displayName over metric name', () => {
      // given
      const data = EvaluationsMock.data.evaluationHistory as Trace[];

      // when
      const points = createChartPoints([data[12]]);

      // then
      expect(points[2].label).toBe('Response time P95');
    });

    it('should correctly set score to failed color', () => {
      // given
      const data = EvaluationsMock.data.evaluationHistory as Trace[];

      // when
      const chartItems = createChartPoints([data[1]]);

      // then
      expect(chartItems[0].points[0].color).toEqual('#dc172a');
    });

    it('should correctly set score to warning color', () => {
      // given
      const data = EvaluationsMock.data.evaluationHistory as Trace[];

      // when
      const chartItems = createChartPoints([data[12]]);

      // then
      expect(chartItems[0].points[0].color).toEqual('#e6be00');
    });

    it('should correctly set score to success color', () => {
      // given
      const data = EvaluationsMock.data.evaluationHistory as Trace[];

      // when
      const chartItems = createChartPoints([data[0]]);

      // then
      expect(chartItems[0].points[0].color).toEqual('#7dc540');
    });
  });

  describe('createChartXLabels', () => {
    it('should correctly return unique labels', () => {
      // given
      const data = EvaluationsMock.data.evaluationHistory as Trace[];

      // when
      const labels = createChartXLabels([data[0], data[1], data[2], data[12], data[12]]);

      // then
      expect(labels).toEqual({
        0: '2020-11-10 12:12',
        1: '2020-11-10 12:15',
        2: '2020-12-21 13:00',
        3: '2021-04-08 17:47 (1)',
        4: '2021-04-08 17:47 (2)',
      });
    });

    it('should not change labels if they are already unique', () => {
      // given
      const data = EvaluationsMock.data.evaluationHistory as Trace[];

      // when
      const labels = createChartXLabels([data[0], data[1], data[2], data[12]]);

      // then
      expect(labels).toEqual({
        0: '2020-11-10 12:12',
        1: '2020-11-10 12:15',
        2: '2020-12-21 13:00',
        3: '2021-04-08 17:47',
      });
    });
  });

  describe('createChartTooltipLabels', () => {
    it('should correctly return tooltip labels', () => {
      // given
      const data = EvaluationsMock.data.evaluationHistory as Trace[];

      // when
      const labels = createChartTooltipLabels([data[0], data[1], data[2]], (date: string) => `myDate${date}`);

      // then
      expect(labels).toEqual({
        0: 'SLO evaluation of test from myDate2020-11-10T11:12:12.364Z',
        1: 'SLO evaluation of test from myDate2020-11-10T11:15:34.488Z',
        2: 'SLO evaluation of test from myDate2020-12-21T12:00:14.126Z',
      });
    });

    it('should return empty if time is not given', () => {
      // given
      const data = EvaluationsMock.data.evaluationHistory as Trace[];

      // when
      const labels = createChartTooltipLabels([data[12], data[12]], () => 'myDate');

      // then
      expect(labels).toEqual({
        0: '',
        1: '',
      });
    });
  });
});
