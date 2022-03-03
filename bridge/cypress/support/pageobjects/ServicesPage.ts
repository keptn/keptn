/// <reference types="cypress" />

import { SliResult } from '../../../client/app/_models/sli-result';

class ServicesPage {
  SERVICE_PANEL_TEXT_LOC = 'dt-info-group-title.dt-info-group-title > div > h2';

  visitServicePage(projectName: string): this {
    cy.visit(`/project/${projectName}/service`);
    return this;
  }

  selectService(serviceName: string, version: string): this {
    cy.byTestId(`keptn-service-view-service-${serviceName}`).click();
    cy.get('dt-row').contains(version).click();
    return this;
  }

  clickSliBreakdownHeader(columnName: string): this {
    cy.byTestId('keptn-sli-breakdown')
      .find('dt-header-row')
      .first()
      .find('dt-header-cell')
      .contains(columnName)
      .click();
    return this;
  }

  verifySliBreakdownSorting(columnIndex: number, direction: string, firstElement: string, secondElement: string): this {
    cy.byTestId('keptn-sli-breakdown')
      .find('dt-header-row')
      .first()
      .find('dt-header-cell')
      .eq(columnIndex)
      .find('.dt-sort-header-container')
      .first()
      .should('have.class', 'dt-sort-header-sorted');
    cy.byTestId('keptn-sli-breakdown')
      .find('dt-header-row')
      .first()
      .find('dt-header-cell')
      .eq(columnIndex)
      .find('dt-icon')
      .first()
      .invoke('attr', 'ng-reflect-name')
      .should('equal', `sorter2-${direction}`);

    cy.byTestId('keptn-sli-breakdown')
      .find('dt-row')
      .eq(0)
      .find('dt-cell')
      .eq(columnIndex)
      .should('have.text', firstElement);
    cy.byTestId('keptn-sli-breakdown')
      .find('dt-row')
      .eq(1)
      .find('dt-cell')
      .eq(columnIndex)
      .should('have.text', secondElement);

    return this;
  }

  verifySliBreakdown(result: SliResult, isExpanded: boolean): this {
    cy.byTestId('keptn-sli-breakdown').should('exist');
    if (isExpanded) {
      this.assertSliColumnText(
        result.name,
        'name',
        `${result.name}Absolute change:Relative change:Compared with:`
      ).assertSliColumnText(
        result.name,
        'value',
        `${result.value}${result.calculatedChanges?.absolute > 0 ? '+' : '-'}${result.calculatedChanges?.absolute}${
          result.calculatedChanges?.relative > 0 ? '+' : '-'
        }${result.calculatedChanges?.relative}% ${result.comparedValue}`
      );
    } else {
      this.assertSliColumnText(result.name, 'name', result.name).assertSliColumnText(
        result.name,
        'value',
        `${result.value} (${result.calculatedChanges?.relative > 0 ? '+' : '-'}${result.calculatedChanges?.relative}%) `
      );
    }

    this.assertSliColumnText(result.name, 'weight', result.weight.toString()).assertSliColumnText(
      result.name,
      'score',
      result.score.toString()
    );

    cy.byTestId(`keptn-sli-breakdown-row-${result.name}`)
      .find('dt-cell')
      .eq(2)
      .find('.error')
      .should('have.length', result.result == 'pass' ? 0 : 1);

    return this;
  }

  assertSliColumnText(sliName: string, columnName: string, value: string): this {
    cy.byTestId(`keptn-sli-breakdown-row-${sliName}`)
      .find(`[uitestid="keptn-sli-breakdown-${columnName}-cell"]`)
      .should('have.text', value);
    return this;
  }

  expandSliBreakdown(name: string): this {
    cy.byTestId(`keptn-sli-breakdown-row-${name}`).find('dt-cell').eq(0).find('button').click();
    return this;
  }

  clickOnServicePanelByName(serviceName: string): this {
    cy.wait(500).get('div.dt-info-group-content').get('h2').contains(serviceName).click();
    return this;
  }

  clickOnServiceInnerPanelByName(serviceName: string): this {
    cy.wait(500).get('span.ng-star-inserted').contains(serviceName).click();
    return this;
  }

  clickEvaluationBoardButton(): this {
    cy.get('button[uitestid="keptn-event-item-contextButton-evaluation"]').click();
    return this;
  }

  clickViewServiceDetails(): this {
    cy.get('.highcharts-plot-background').should('be.visible');
    cy.contains('View service details').click();
    return this;
  }

  clickViewSequenceDetails(): this {
    cy.contains('View sequence details').click();
    return this;
  }

  clickGoBack(): this {
    cy.contains('Go back').click();
    return this;
  }

  verifyCurrentOpenServiceNameEvaluationPanel(serviceName: string): this {
    cy.get('div.service-title > span').should('have.text', serviceName);
    cy.get('.highcharts-plot-background').should('be.visible');
    return this;
  }
}
export default ServicesPage;
