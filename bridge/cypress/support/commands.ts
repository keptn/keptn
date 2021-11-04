/* eslint-disable import/no-extraneous-dependencies */
import 'cypress-file-upload';
import 'cypress-wait-until';
/* eslint-enable import/no-extraneous-dependencies */

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
const forceClick = (selector: HTMLElement) =>
  cy
    .waitUntil(
      () =>
        cy
          .wrap(Cypress.$(selector))
          .as('clickSubject')
          .wait(10) // for some reason this is needed, otherwise next line returns `true` even if click() fails due to detached element in the next step
          .then(($el) => Cypress.dom.isAttached($el)),
      { timeout: 1000, interval: 100 }
    )
    .get('@clickSubject')
    .click();
/* eslint-enable @typescript-eslint/explicit-function-return-type */
// Commands have to be added by hooking them to Cypress
Cypress.Commands.add('byTestId', byTestId);
Cypress.Commands.add('forceClick', { prevSubject: 'element' }, forceClick);
