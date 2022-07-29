import { LineChartComponentPage } from '../../support/pageobjects/LineChartComponentPage';

describe('Test line chart', () => {
  const lineChartComponentPage = new LineChartComponentPage();

  beforeEach(() => {
    lineChartComponentPage.intercept().visitPageWithHeatmapComponent().selectLineChart();
    cy.wait(100000);
  });
  //

  describe('Test legend and filter', () => {
    it('should not be possible to hide all metrics', () => {});

    it('should show the displayName in favor of the metric name', () => {});

    it('should enable and disable metric', () => {});

    it('should initially hide all SLIs', () => {});

    it('should show SLI if enabled', () => {});
  });

  describe('Test xAxis', () => {
    it('should correctly map duplicate dates to unique dates', () => {});

    it('should show buildIDs in favor of the evaluation time', () => {});

    it('should always show 25, 50, 75, 100 for the xAxis', () => {
      // test with all/one score(s) 0, 25, 50, 75, 100
    });
  });

  describe('Test yAxis', () => {
    it('should expand the yAxis if an SLI with a bigger range is enabled', () => {});

    it('should collapse the yAxis if an SLI with a big range compared to others is disabled', () => {});
  });

  describe('Test tooltip', () => {
    // how to properly test tooltip?
    // one solution would be: hover on metric X (metric X may start after/before metric Y)
    // other solution: take position of the xAxis-label and go some pixels up
    it('should show tooltip on hover', () => {});

    it('should show tooltip only with enabled SLIs and score on hover', () => {});

    it('should show score in tooltip if one of the two scores is enabled', () => {});

    it('should not show score in tooltip if both scores are disabled', () => {});

    it('should change tooltip if another evaluation is hovered', () => {});
  });
});
