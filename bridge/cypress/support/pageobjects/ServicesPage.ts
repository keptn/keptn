/// <reference types="cypress" />

class ServicesPage {
  clickOnServicePanelByName(serviceName: string): this {
    cy.get('div.dt-info-group-content').get('h2').contains(serviceName).click();
    return this;
  }

  clickOnServiceInnerPanelByName(serviceName: string): this {
    cy.get('span.ng-star-inserted').contains(serviceName).click();
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
