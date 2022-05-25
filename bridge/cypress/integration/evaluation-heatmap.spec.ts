import { HeatmapComponent } from '../support/pageobjects/HeatmapComponent';

describe('evaluation-heatmap', () => {
  const heatmap = new HeatmapComponent();
  beforeEach(() => {
    heatmap.intercept().visitPageWithHeatmapComponent();
  });
  it('should display ktb-heatmap if the feature flag is enabled', () => {
    heatmap.assertComponentExists();
  });
  it('should be expandable and collapsable', () => {
    heatmap
      .assertNumberOfRows(10)
      .clickExpandButton()
      .assertNumberOfRows(13)
      .clickExpandButton()
      .assertNumberOfRows(10);
  });
  it('should set correct color classes', () => {
    heatmap
      .assertTileColor('ktb-heatmap-tile-182d10b8-b68d-49d4-86cd-5521352d7a42', 'pass')
      .assertTileColor('ktb-heatmap-tile-182d10b8-b68d-49d4-86cd-5521352d7a42', 'warning')
      .assertTileColor('ktb-heatmap-tile-52b4b2c7-fa49-41f3-9b5c-b9aea2370bb4', 'fail');
  });
  it('should have a primary and secondary highlight', () => {
    heatmap.assertPrimaryHighlightExists().assertSecondaryHighlightExists();
  });
  it('should truncate long metric names', () => {
    const longName = 'A very long metric name so long it gets cut somewhere along the way';
    const shortName = 'A very long metric name ...';
    heatmap.assertMetricIsTruncated(longName, shortName);
  });
  it('should show and hide tooltips', () => {
    const tileId = 'ktb-heatmap-tile-182d10b8-b68d-49d4-86cd-5521352d7a42';
    heatmap
      .assertTooltipIsHidden()
      .mouseOverOnTile(tileId)
      .assertTooltipIsVisible()
      .mouseLeaveOnTile(tileId)
      .assertTooltipIsHidden();
  });
  it('should enumerate duplicate dates', () => {
    heatmap.assertXAxisLabelExistsOnce('2022-02-08 12:56 (1)').assertXAxisLabelExistsOnce('2022-02-08 12:56 (2)');
  });
  it('should show all dates', () => {
    heatmap.assertXAxisTickLength(10);
  });
  it('should not show secondary highlight if clicked evaluation has no other to compare', () => {
    heatmap
      .clickScore('ktb-heatmap-tile-25ab0f26-e6d8-48d5-a08f-08c8a136a688')
      .assertPrimaryHighlightExists()
      .assertSecondaryHighlightDoesNotExist();
  });
});
