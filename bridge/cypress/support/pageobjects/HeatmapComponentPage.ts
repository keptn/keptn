import { interceptHeatmapComponent } from '../intercept';

export type ResultState = 'pass' | 'warning' | 'fail' | 'info';

export class HeatmapComponentPage {
  private readonly tilePrefix = 'ktb-heatmap-tile-';

  public intercept(): this {
    interceptHeatmapComponent();
    return this;
  }

  interceptWithManyEvaluations(): this {
    cy.intercept('GET', 'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
      statusCode: 200,
      fixture: 'get.sockshop.service.carts.evaluations.heatmap.manyscores.mock.json',
    });
    return this;
  }

  interceptWith10Metrics(): this {
    cy.intercept('GET', 'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
      statusCode: 200,
      fixture: 'get.sockshop.service.carts.evaluations.heatmap.10metrics.mock.json',
    });
    return this;
  }

  interceptWithTwoHighlights(): this {
    cy.intercept('GET', 'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
      statusCode: 200,
      fixture: 'get.sockshop.service.carts.evaluations.heatmap.twohighlights.mock.json',
    });
    return this;
  }

  visitPageWithHeatmapComponent(): this {
    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging');
    return this;
  }

  clickExpandButton(): this {
    cy.get('ktb-heatmap button.show-more-button').click();
    return this;
  }

  clickScore(scoreTestId: string): this {
    cy.get(
      `ktb-heatmap .data-point-container g[uitestid="Score"] rect[uitestid="${this.tilePrefix + scoreTestId}"]`
    ).click();
    return this;
  }

  clickMetric(metricName: string, metricTestId: string): this {
    cy.get(
      `ktb-heatmap .data-point-container g[uitestid="${metricName}"] rect[uitestid="${this.tilePrefix + metricTestId}"]`
    ).click();
    return this;
  }

  clickLegendCircle(circle: ResultState): this {
    // Wait, because element changes after being clicked
    cy.wait(1000);
    cy.contains(`ktb-heatmap .legend-container .legend-item text`, circle).click();
    return this;
  }

  mouseOverOnTile(tileId: string): this {
    cy.byTestId(this.tilePrefix + tileId)
      .first()
      .trigger('mouseover');
    return this;
  }

  mouseLeaveOnTile(tileId: string): this {
    cy.byTestId(this.tilePrefix + tileId)
      .first()
      .trigger('mouseleave');
    return this;
  }

  assertComponentExists(): this {
    cy.get('ktb-heatmap').should('exist');
    return this;
  }

  assertNumberOfRows(length: number): this {
    cy.get('ktb-heatmap svg .y-axis-container .tick').should('have.length', length);
    return this;
  }

  assertTileColor(tileId: string, color: ResultState, disabled = false): this {
    cy.byTestId(this.tilePrefix + tileId)
      .should('have.class', color)
      .and('have.class', 'data-point')
      .and(disabled ? 'have.class' : 'not.have.class', 'disabled');
    return this;
  }

  private assertHighlight(name: '.highlight-primary' | '.highlight-secondary', amount: number): this {
    if (amount <= 0) {
      cy.get(`ktb-heatmap ${name}`).should('not.exist');
      return this;
    }
    cy.get(`ktb-heatmap ${name}`).should('have.length', amount);
    return this;
  }

  assertPrimaryHighlight(amountOfExistence: number): this {
    return this.assertHighlight('.highlight-primary', amountOfExistence);
  }

  assertSecondaryHighlight(amountOfExistence: number): this {
    return this.assertHighlight('.highlight-secondary', amountOfExistence);
  }

  assertMetricIsTruncated(longName: string, truncatedName: string): this {
    cy.contains('ktb-heatmap title', longName)
      .parent()
      .within(() => {
        cy.contains('text', truncatedName);
      });
    return this;
  }

  assertTooltipIsVisible(visible: boolean): this {
    const should = visible ? 'not.have.class' : 'have.class';
    cy.get('ktb-heatmap ktb-heatmap-tooltip').should(should, 'hidden');
    return this;
  }

  assertTooltipIs(type: 'score' | 'metric'): this {
    const selector = 'ktb-heatmap ktb-heatmap-tooltip dt-key-value-list-item dt-key-value-list-key';
    cy.contains(selector, 'Value').should('exist');
    if (type === 'score') {
      cy.contains(selector, 'Total passed SLIs').should('exist');
      cy.contains(selector, 'Total warning SLIs').should('exist');
      cy.contains(selector, 'Total failed SLIs').should('exist');
    } else if (type === 'metric') {
      cy.contains(selector, 'Score').should('exist');
      cy.contains(selector, 'Key SLI').should('exist');
    }
    return this;
  }

  assertXAxisLabelExistsOnce(labelText: string): this {
    cy.contains('ktb-heatmap .x-axis-container g.tick text', labelText).should('have.length', 1);
    return this;
  }

  assertXAxisTickLength(length: number): this {
    cy.get('ktb-heatmap .x-axis-container g.tick').should('have.length', length);
    return this;
  }

  assertXAxisTickLabels(labels: string[]): this {
    const sorter = (a: string, b: string): number => a.localeCompare(b);
    const sortedLabels = [...labels].sort(sorter);
    // eslint-disable-next-line promise/catch-or-return
    cy.get('ktb-heatmap .x-axis-container g.tick')
      .then(($els) =>
        Cypress.$.makeArray($els)
          .map((el) => el.textContent ?? '')
          .sort(sorter)
      )
      .should('deep.equal', sortedLabels);
    return this;
  }

  assertExpandExists(exists: boolean): this {
    cy.get('ktb-heatmap button.show-more-button').should(!exists ? 'not.exist' : 'exist');
    return this;
  }
}

export const range = (start: number, stop: number, step: number): number[] =>
  Array.from({ length: (stop - start) / step + 1 }, (_, i) => start + i * step);
