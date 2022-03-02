import Chainable = Cypress.Chainable;

export class ServicesSettingsPage {
  inputService(serviceName: string): Chainable<JQuery<HTMLElement>> {
    return cy.get('input[formcontrolname="serviceName"]').type(serviceName);
  }

  createService(serviceName?: string): Chainable<JQuery<HTMLElement>> {
    if (serviceName) {
      this.inputService(serviceName);
    }
    return cy.byTestId('createServiceButton').click();
  }
}
