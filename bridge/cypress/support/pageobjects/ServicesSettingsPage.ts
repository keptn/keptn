import Chainable = Cypress.Chainable;

export class ServicesSettingsPage {
  inputService(serviceName: string): Chainable<JQuery<HTMLElement>> {
    return cy.get('input[formcontrolname="serviceName"]').type(serviceName);
  }

  createService(): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('createServiceButton').click();
  }
}
