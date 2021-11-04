/* eslint-disable import/no-extraneous-dependencies */
import 'cypress-file-upload';
/* eslint-enable import/no-extraneous-dependencies */

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace Cypress {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    interface Chainable<Subject> {
      byTestId<E extends Node = HTMLElement>(id: string): Cypress.Chainable<JQuery<E>>;
    }
  }
}
/* eslint-disable @typescript-eslint/explicit-function-return-type */
const byTestId = (testId: string) => cy.get(`[uitestid="${testId}"]`);
/* eslint-enable @typescript-eslint/explicit-function-return-type */
// Commands have to be added by hooking them to Cypress
Cypress.Commands.add('byTestId', byTestId);
