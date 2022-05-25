import { interceptHeatmapComponent } from '../intercept';

export class HeatmapComponent {
  public intercept(): this {
    interceptHeatmapComponent();
    return this;
  }

  visitPageWithHeatmapComponent(): this {
    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging');
    return this;
  }

  clickExpandButton(): this {
    cy.get('ktb-heatmap button').click();
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

  assertTileColor(tileId: string, color: 'pass' | 'warning' | 'fail'): this {
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
}
