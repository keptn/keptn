import { interceptHeatmapComponent } from '../intercept';

export class HeatmapComponent {
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

  visitPageWithHeatmapComponent(): this {
    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging');
    return this;
  }

  clickExpandButton(): this {
    cy.get('ktb-heatmap button.show-more-button').click();
    return this;
  }

  clickScore(scoreTestId: string): this {
    cy.get(`ktb-heatmap .data-point-container g[uitestid="Score"] rect[uitestid="${scoreTestId}"]`).click();
    return this;
  }

  mouseOverOnTile(tileId: string): this {
    cy.byTestId(tileId).first().trigger('mouseover');
    return this;
  }

  mouseLeaveOnTile(tileId: string): this {
    cy.byTestId(tileId).first().trigger('mouseleave');
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

  assertTileColor(tileId: string, color: 'pass' | 'warning' | 'fail' | 'info'): this {
    cy.byTestId(tileId).should('have.class', color).and('have.class', 'data-point');
    return this;
  }

  private assertHighlight(name: '.highlight-primary' | '.highlight-secondary'): this {
    cy.get(`ktb-heatmap ${name}`).should('have.length', 1);
    return this;
  }

  assertPrimaryHighlightExists(): this {
    return this.assertHighlight('.highlight-primary');
  }

  assertSecondaryHighlightExists(): this {
    return this.assertHighlight('.highlight-secondary');
  }

  assertSecondaryHighlightDoesNotExist(): this {
    cy.get(`ktb-heatmap .highlight-secondary`).should('not.exist');
    return this;
  }

  assertMetricIsTruncated(longName: string, truncatedName: string): this {
    cy.contains('ktb-heatmap title', longName)
      .parent()
      .within(() => {
        cy.contains('text', truncatedName);
      });
    return this;
  }

  assertTooltipIsVisible(): this {
    cy.get('ktb-heatmap ktb-heatmap-tooltip').should('not.have.class', 'hidden');
    return this;
  }

  assertTooltipIsHidden(): this {
    cy.get('ktb-heatmap ktb-heatmap-tooltip').should('have.class', 'hidden');
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
