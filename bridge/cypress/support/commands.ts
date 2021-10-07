import 'cypress-file-upload';

// eslint-disable-next-line @typescript-eslint/no-namespace
declare namespace Cypress {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  interface Chainable<Subject> {
    byTestId<E extends Node = HTMLElement>(id: string): Cypress.Chainable<JQuery<E>>;
  }
}

// eslint-disable-next-line @typescript-eslint/explicit-function-return-type
const byTestId = (testId: string) => cy.get(`[uitestid="${testId}"]`);

// Commands have to be added by hooking them to Cypress
Cypress.Commands.add('byTestId', byTestId);
