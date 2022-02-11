/* eslint-disable import/no-extraneous-dependencies */
import 'cypress-file-upload';
/* eslint-enable import/no-extraneous-dependencies */

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace Cypress {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    interface Chainable<Subject> {
      byTestId<E extends Node = HTMLElement>(id: string): Cypress.Chainable<JQuery<E>>;
      clickOutside<E extends Node = HTMLElement>(): Cypress.Chainable<JQuery<E>>;
    }
  }
}
// Commands have to be added by hooking them to Cypress
Cypress.Commands.add('byTestId', (testId: string) => cy.get(`[uitestid="${testId}"]`));
Cypress.Commands.add('clickOutside', () => cy.get('body').click(0, 0));
