import { HeatmapComponentPage, range } from '../support/pageobjects/HeatmapComponentPage';

describe('evaluation-heatmap', () => {
  const heatmap = new HeatmapComponentPage();
  beforeEach(() => {
    heatmap.intercept().visitPageWithHeatmapComponent();
  });
  it('should display ktb-heatmap if the feature flag is enabled', () => {
    heatmap.assertComponentExists();
  });
  it('should be expandable and collapsable', () => {
    heatmap
      .assertNumberOfRows(10)
      .assertExpandExists(true)
      .clickExpandButton()
      .assertNumberOfRows(13)
      .assertExpandExists(true)
      .clickExpandButton()
      .assertNumberOfRows(10)
      .assertExpandExists(true);
  });
  it('should set correct color classes', () => {
    heatmap
      .assertTileColor('182d10b8-b68d-49d4-86cd-5521352d7a42', 'pass')
      .assertTileColor('182d10b8-b68d-49d4-86cd-5521352d7a42', 'warning')
      .assertTileColor('52b4b2c7-fa49-41f3-9b5c-b9aea2370bb4', 'fail')
      .clickExpandButton()
      .assertTileColor('e074893e-a7f9-4fa8-9e7e-898937a3d2b6', 'info');
  });
  it('should have a primary and one secondary highlight', () => {
    heatmap.assertPrimaryHighlight(1).assertSecondaryHighlight(1);
  });
  it('should have a primary and two secondary highlights', () => {
    heatmap.interceptWithTwoHighlights().assertPrimaryHighlight(1).assertSecondaryHighlight(2);
  });
  it('should truncate long metric names', () => {
    const longName = 'A very long metric name so long it gets cut somewhere along the way';
    const shortName = 'A very long metric name ...';
    heatmap.assertMetricIsTruncated(longName, shortName);
  });
  it('should show and hide tooltips', () => {
    const tileId = '182d10b8-b68d-49d4-86cd-5521352d7a42';
    heatmap
      .assertTooltipIsVisible(false)
      .mouseOverOnTile(tileId)
      .assertTooltipIsVisible(true)
      .mouseLeaveOnTile(tileId)
      .assertTooltipIsVisible(false);
  });
  it('should enumerate duplicate dates', () => {
    heatmap.assertXAxisLabelExistsOnce('2022-02-08 12:56 (1)').assertXAxisLabelExistsOnce('2022-02-08 12:56 (2)');
  });
  it('should show all dates', () => {
    heatmap.assertXAxisTickLength(10);
  });
  it('should not show secondary highlight if clicked evaluation has no other to compare', () => {
    heatmap.clickScore('25ab0f26-e6d8-48d5-a08f-08c8a136a688').assertPrimaryHighlight(1).assertSecondaryHighlight(0);
  });
  it('should reduce X elements on many evaluations', () => {
    const labels = range(1, 44, 2).map((value) => `2022-02-01 03:46 (${value})`);
    heatmap.interceptWithManyEvaluations().assertXAxisTickLength(22).assertXAxisTickLabels(labels);
  });
  it('should not show expand button if indicator results with score are less than 10', () => {
    heatmap.interceptWith10Metrics().assertExpandExists(false);
  });
  it('should set tiles disabled/enabled via legend', () => {
    heatmap
      .clickLegendCircle('pass')
      .assertTileColor('182d10b8-b68d-49d4-86cd-5521352d7a42', 'pass', true)
      .clickLegendCircle('pass')
      .assertTileColor('182d10b8-b68d-49d4-86cd-5521352d7a42', 'pass', false)
      .clickLegendCircle('warning')
      .assertTileColor('182d10b8-b68d-49d4-86cd-5521352d7a42', 'warning', true)
      .clickLegendCircle('warning')
      .assertTileColor('182d10b8-b68d-49d4-86cd-5521352d7a42', 'warning', false)
      .clickLegendCircle('info')
      .clickExpandButton()
      .assertTileColor('e074893e-a7f9-4fa8-9e7e-898937a3d2b6', 'info', true)
      .clickLegendCircle('info')
      .assertTileColor('e074893e-a7f9-4fa8-9e7e-898937a3d2b6', 'info', false)
      .clickLegendCircle('fail')
      .assertTileColor('52b4b2c7-fa49-41f3-9b5c-b9aea2370bb4', 'fail', true)
      .clickLegendCircle('fail')
      .assertTileColor('52b4b2c7-fa49-41f3-9b5c-b9aea2370bb4', 'fail', false);
  });
  it('should not show tooltip if tile is disabled', () => {
    heatmap
      .clickLegendCircle('pass')
      .clickScore('182d10b8-b68d-49d4-86cd-5521352d7a42')
      .assertTooltipIsVisible(false)
      .clickLegendCircle('warning')
      .clickMetric('go_routines3a', '182d10b8-b68d-49d4-86cd-5521352d7a42')
      .assertTooltipIsVisible(false);
  });
  it('should show score/metric tooltip accordingly', () => {
    heatmap
      .clickScore('182d10b8-b68d-49d4-86cd-5521352d7a42')
      .assertTooltipIsVisible(true)
      .assertTooltipIs('score')
      .clickMetric('go_routines3a', '182d10b8-b68d-49d4-86cd-5521352d7a42')
      .assertTooltipIsVisible(true)
      .assertTooltipIs('metric');
  });
});
