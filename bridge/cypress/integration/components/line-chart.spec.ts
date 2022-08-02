import { LineChartComponentPage } from '../../support/pageobjects/LineChartComponentPage';

describe('Test line chart', () => {
  const lineChartComponentPage = new LineChartComponentPage();
  const slis = [
    'Response-time-P95',
    'http_response_time_seconds_main_page_sum',
    'request_throughput',
    'go_routines',
    'A-very-long-metric-name-so-long-it-gets-cut-somewhere-along-the-way',
    'go_routines2',
    'go_routines3',
    'go_routines4',
    'go_routines5',
    'go_routines6',
    'go_routines7',
    'go_routines3a',
  ];

  beforeEach(() => {
    lineChartComponentPage.intercept();
  });

  describe('Test legend and filter', () => {
    beforeEach(() => {
      lineChartComponentPage.visitPageWithHeatmapComponent().selectLineChart();
    });

    it('should initially show 10 score bars', () => {
      lineChartComponentPage.assertScoreBarEnabled(true, 10);
    });

    it('should initially hide all SLIs', () => {
      lineChartComponentPage.assertScoreBarEnabled(true).assertScoreLineEnabled(true).assertMetricCount(14);

      for (const sli of slis) {
        lineChartComponentPage.assertMetricEnabled(sli, false);
      }
    });

    it('should not be possible to hide all metrics', () => {
      lineChartComponentPage
        .toggleMetric('score')
        .assertScoreBarEnabled(false)
        .toggleMetric('score', 1)
        .assertScoreLineEnabled(true);
    });

    it('should show the displayName in favor of the metric name', () => {
      lineChartComponentPage.assertMetricName('Response-time-P95', 'Response time P95');
    });

    it('should disable and enable metric', () => {
      lineChartComponentPage
        .toggleMetric('score')
        .assertScoreBarEnabled(false)
        .toggleMetric('score')
        .assertScoreBarEnabled(true);
    });

    it('should enable and disable metric', () => {
      lineChartComponentPage
        .toggleMetric('Response-time-P95')
        .assertMetricEnabled('Response-time-P95', true)
        .toggleMetric('Response-time-P95')
        .assertMetricEnabled('Response-time-P95', false);
    });
  });

  describe('Test xAxis', () => {
    beforeEach(() => {
      lineChartComponentPage.visitPageWithHeatmapComponent().selectLineChart();
    });

    it('should correctly map duplicate dates to unique dates', () => {
      const dates = [
        { selector: '2022-02-08-12:56-(1)', value: '2022-02-08 12:56 (1)' },
        { selector: '2022-02-08-12:56-(2)', value: '2022-02-08 12:56 (2)' },
      ];
      lineChartComponentPage.assertXAxisLabelCount(10);
      for (const date of dates) {
        lineChartComponentPage.assertXAxisLabel(date.selector, date.value);
      }
    });

    it('should show buildIDs in favor of the evaluation time', () => {
      lineChartComponentPage.assertXAxisLabelCount(10).assertXAxisLabel('myBuildId', 'myBuildId');
    });

    it('should correctly show all dates', () => {
      const dates = [
        { selector: 'myBuildId', value: 'myBuildId' },
        { selector: '2022-02-08-12:56-(1)', value: '2022-02-08 12:56 (1)' },
        { selector: '2022-02-08-12:56-(2)', value: '2022-02-08 12:56 (2)' },
        { selector: '2022-02-08-13:17', value: '2022-02-08 13:17' },
        { selector: '2022-02-08-13:30', value: '2022-02-08 13:30' },
        { selector: '2022-02-08-13:33', value: '2022-02-08 13:33' },
        { selector: '2022-02-08-13:43', value: '2022-02-08 13:43' },
        { selector: '2022-02-08-13:57', value: '2022-02-08 13:57' },
        { selector: '2022-02-08-14:01', value: '2022-02-08 14:01' },
        { selector: '2022-02-09-10:30', value: '2022-02-09 10:30' },
      ];
      lineChartComponentPage.assertXAxisLabelCount(10);
      for (const date of dates) {
        lineChartComponentPage.assertXAxisLabel(date.selector, date.value);
      }
    });

    it('should only show evaluations related to the activated metrics', () => {
      lineChartComponentPage.toggleMetric('Response-time-P95').toggleScores().assertXAxisLabelCount(4);
    });
  });

  describe('Test yAxis', () => {
    beforeEach(() => {
      lineChartComponentPage.visitPageWithHeatmapComponent().selectLineChart();
    });

    it('should always show 25, 50, 75, 100 for the left yAxis', () => {
      const scorePairs = [
        { score1: 0, score2: 0 },
        { score1: 25, score2: 25 },
        { score1: 50, score2: 50 },
        { score1: 75, score2: 75 },
        { score1: 100, score2: 100 },
        { score1: 0, score2: 100 },
      ];
      for (const scorePair of scorePairs) {
        lineChartComponentPage
          .interceptWithScore(scorePair.score1, scorePair.score2)
          .visitPageWithHeatmapComponent()
          .selectLineChart()
          .assertYAxisLeftLabelCount(4)
          .assertYAxisLeftLabels();
      }
      lineChartComponentPage.toggleMetric('score', 0).assertYAxisLeftLabels();
    });

    it('should not change a big right yAxis if score is enabled/disabled', () => {
      const labels = ['2,500', '5,000', '7,500', '10,000'] as const;
      lineChartComponentPage
        .toggleMetric('request_throughput')
        .assertYAxisRightLabels(...labels)
        .toggleScores()
        .assertYAxisRightLabels(...labels);
    });

    it('should not change a small right yAxis if score is enabled/disabled', () => {
      const labels = ['0.5', '1', '1.5', '2'] as const;
      lineChartComponentPage
        .toggleMetric('http_response_time_seconds_main_page_sum')
        .assertYAxisRightLabels(...labels)
        .toggleScores()
        .assertYAxisRightLabels(...labels);
    });
  });

  describe('Test tooltip', () => {
    beforeEach(() => {
      lineChartComponentPage.visitPageWithHeatmapComponent().selectLineChart();
    });

    it('should show tooltip on hover', () => {
      lineChartComponentPage
        .showTooltip('2022-02-08-13:30')
        .assertTooltipHeader('2022-02-08 13:30')
        .assertToolTipValue('score', '100')
        .assertTooltipValueCount(1);
    });

    it('should show tooltip only with enabled SLIs and score on hover', () => {
      lineChartComponentPage
        .toggleMetric('go_routines')
        .toggleMetric('request_throughput')
        .showTooltip('2022-02-09-10:30')
        .assertTooltipHeader('2022-02-09 10:30')
        .assertToolTipValue('score', '33.99')
        .assertToolTipValue('go_routines', '88')
        .assertToolTipValue('request_throughput', '10000')
        .assertTooltipValueCount(3);
    });

    it('should show displayName in favor of metric in tooltip for SLI', () => {
      lineChartComponentPage
        .toggleMetric('Response-time-P95')
        .showTooltip('2022-02-08-13:17')
        .assertToolTipValue('Response time P95', '0');
    });

    it('should show score in tooltip if one of the two scores is enabled', () => {
      lineChartComponentPage
        .toggleMetric('Response-time-P95')
        .showTooltip('2022-02-08-13:17')
        .assertToolTipValue('score', '0')
        .toggleMetric('score', 0)
        .assertToolTipValue('score', '0')
        .toggleMetric('score', 0)
        .toggleMetric('score', 1)
        .assertToolTipValue('score', '0');
    });

    it('should not show score in tooltip if both scores are disabled', () => {
      lineChartComponentPage
        .toggleMetric('Response-time-P95')
        .toggleScores()
        .showTooltip('2022-02-08-13:17')
        .assertTooltipValueCount(1);
    });

    it('should not show metric in tooltip if the evaluation does not have it', () => {
      lineChartComponentPage
        .toggleMetric('Response-time-P95')
        .toggleMetric('go_routines')
        .toggleScores()
        .showTooltip('2022-02-08-13:17')
        .assertToolTipValue('Response time P95', '0')
        .assertTooltipValueCount(1);
    });

    it('should change tooltip if another evaluation is hovered', () => {
      lineChartComponentPage
        .showTooltip('2022-02-08-13:30')
        .assertTooltipHeader('2022-02-08 13:30')
        .assertToolTipValue('score', '100')
        .showTooltip('2022-02-09-10:30')
        .assertTooltipHeader('2022-02-09 10:30')
        .assertToolTipValue('score', '33.99');
    });
  });
});
