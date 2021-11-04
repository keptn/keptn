/// <reference types="cypress" />

class ServicesPage {
  clickOnServicePanelByName(serviceName: string): this {
    cy.get('div.dt-info-group-content').get('h2').contains(serviceName).forceClick();
    return this;
  }

  clickOnServiceInnerPanelByName(serviceName: string): this {
    cy.get('span.ng-star-inserted').contains(serviceName).forceClick();
    return this;
  }

  clickEvaluationBoardButton(): this {
    cy.get('button[uitestid="keptn-event-item-contextButton-evaluation"]').forceClick();
    return this;
  }

  clickViewServiceDetails(): this {
    cy.get('.highcharts-plot-background').should('be.visible');
    cy.contains('View service details').forceClick();
    return this;
  }

  clickViewSequenceDetails(): this {
    cy.contains('View sequence details').forceClick();
    return this;
  }

  clickGoBack(): this {
    cy.contains('Go back').forceClick();
    return this;
  }

  verifyCurrentOpenServiceNameEvaluationPanel(serviceName: string): this {
    cy.get('div.service-title > span').should('have.text', serviceName);
    cy.get('.highcharts-plot-background').should('be.visible');
    return this;
  }
}
export default ServicesPage;
