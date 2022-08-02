import { interceptHeatmapComponent, interceptHeatmapComponentWithScores } from '../intercept';

// index for selecting the right score item
enum ScoreSelectionIndex {
  BAR = 0,
  LINE = 1,
}

export class LineChartComponentPage {
  public intercept(): this {
    interceptHeatmapComponent();
    return this;
  }

  public interceptWithScore(score1: number, score2: number): this {
    interceptHeatmapComponentWithScores(score1, score2);
    return this;
  }

  public visitPageWithHeatmapComponent(): this {
    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging').wait(
      '@heatmapEvaluations'
    );
    return this;
  }

  public selectLineChart(): this {
    cy.byTestId('keptn-evaluation-details-contextButtons-chart').click();
    return this;
  }

  public toggleMetric(metric: string, selectorIndex = 0 /*if there are two metrics like 'score'*/): this {
    cy.byTestId(`chart-legend-item-${metric}`).eq(selectorIndex).click();
    return this;
  }

  public toggleScores(): this {
    return this.toggleMetric('score', 0).toggleMetric('score', 1);
  }

  public assertMetricCount(count: number): this {
    cy.get('.ktb-chart-legend-item').should('have.length', count);
    return this;
  }

  public assertScoreBarEnabled(status: boolean, barCount?: number): this {
    const metric = 'score';
    this.assertChartLegendItemEnabled(metric, status, ScoreSelectionIndex.BAR).assertBarMetricExists(metric, status);
    if (barCount !== undefined) {
      this.assertBarMetricCount(metric, barCount);
    }
    return this;
  }

  private assertBarMetricExists(metric: string, status: boolean): this {
    cy.byTestId(`bar-${metric}`)
      .find('rect')
      .should(status ? 'exist' : 'not.exist');
    return this;
  }

  private assertBarMetricCount(metric: string, count: number): this {
    cy.byTestId(`bar-${metric}`).find('rect').should('have.length', count);
    return this;
  }

  public assertScoreLineEnabled(status: boolean): this {
    const metric = 'score';
    return this.assertChartLegendItemEnabled(metric, status, ScoreSelectionIndex.LINE).assertLineMetricExists(
      metric,
      status
    );
  }

  private assertChartLegendItemEnabled(metric: string, status: boolean, selectorIndex = 0): this {
    cy.byTestId(`chart-legend-item-${metric}`)
      .eq(selectorIndex)
      .should(status ? 'not.have.class' : 'have.class', 'invisible');
    return this;
  }

  private assertLineMetricExists(metric: string, status: boolean): this {
    cy.byTestId(`line-${metric}`).should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertMetricEnabled(metric: string, status: boolean): this {
    return this.assertChartLegendItemEnabled(metric, status).assertLineMetricExists(metric, status);
  }

  public assertMetricName(metric: string, name: string): this {
    cy.byTestId(`chart-legend-item-${metric}`).should('have.text', name);
    return this;
  }

  public assertXAxisLabelCount(count: number): this {
    cy.byTestId('axis-x').find('.tick').should('have.length', count);
    return this;
  }

  public assertXAxisLabel(selector: string, value: string): this {
    cy.byTestId(`axis-x-item-${selector}`).should('have.text', value);
    return this;
  }

  public assertYAxisLeftLabelCount(count: number): this {
    cy.byTestId('axis-y-left').find('.tick').should('have.length', count);
    return this;
  }

  public assertYAxisLeftLabel(value: number): this {
    cy.byTestId('axis-y-left').find('.tick text').contains(value.toString());
    return this;
  }

  public assertYAxisRightLabel(value: string): this {
    cy.byTestId('axis-y-right').find('.tick text').contains(value);
    return this;
  }

  public assertYAxisLeftLabels(): this {
    for (const label of [25, 50, 75, 100]) {
      this.assertYAxisLeftLabel(label);
    }
    return this;
  }

  public assertYAxisRightLabels(...labels: [string, string, string, string]): this {
    for (const label of labels) {
      this.assertYAxisRightLabel(label);
    }
    return this;
  }

  public showTooltip(normalizedEvaluationDate: string): this {
    cy.byTestId(`area-${normalizedEvaluationDate}`).trigger('mouseenter').trigger('mousemove');
    return this;
  }

  public assertTooltipHeader(date: string): this {
    cy.get('.tooltip.tooltip-container>h4').should('have.text', `SLO evaluation of test from ${date}`);
    return this;
  }

  public assertToolTipValue(name: string, value: string): this {
    cy.get('.tooltip.tooltip-container .dt-key-value-list-item')
      .find('.dt-key-value-list-item-key')
      .contains(name)
      .parentsUntil('.dt-key-value-list-item')
      .parent()
      .find('.dt-key-value-list-item-value')
      .should('have.text', value);
    return this;
  }

  public assertTooltipValueCount(count: number): this {
    cy.get('.tooltip.tooltip-container .dt-key-value-list-item').should('have.length', count);
    return this;
  }
}
