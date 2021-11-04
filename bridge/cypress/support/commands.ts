import 'cypress-file-upload';

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace Cypress {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    interface Chainable<Subject> {
      byTestId<E extends Node = HTMLElement>(id: string): Cypress.Chainable<JQuery<E>>;
      forceClick<E extends Node = HTMLElement>(): Cypress.Chainable<JQuery<E>>;
    }
  }
}
/* eslint-disable @typescript-eslint/explicit-function-return-type */
const byTestId = (testId: string) => cy.get(`[uitestid="${testId}"]`);
const forceClick = (selector: HTMLElement) => cy.wrap(Cypress.$(selector)).should('be.visible').click(); // cy.wrap automatically retries when the .should assertion passes
/* eslint-enable @typescript-eslint/explicit-function-return-type */
// Commands have to be added by hooking them to Cypress
Cypress.Commands.add('byTestId', byTestId);
Cypress.Commands.add('forceClick', { prevSubject: 'element' }, forceClick);
