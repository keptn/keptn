/* eslint-disable import/no-extraneous-dependencies */
import 'cypress-file-upload';
/* eslint-enable import/no-extraneous-dependencies */

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace Cypress {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    interface Chainable<Subject> {
      byTestId<E extends Node = HTMLElement>(id: string): Cypress.Chainable<JQuery<E>>;
      parentsUntilTestId<E extends Node = HTMLElement>(id: string): Cypress.Chainable<JQuery<E>>;
      toggleSwitch<E extends Node = HTMLElement>(status: boolean): Cypress.Chainable<JQuery<E>>;
      dtSelect<E extends Node = HTMLElement>(element: string): Cypress.Chainable<JQuery<E>>;
      clickOutside<E extends Node = HTMLElement>(): Cypress.Chainable<JQuery<E>>;
    }
  }
}
// Commands have to be added by hooking them to Cypress
Cypress.Commands.add('byTestId', (testId: string) => cy.get(`[uitestid="${testId}"]`));
Cypress.Commands.add('clickOutside', () => cy.get('body').click(0, 0));
Cypress.Commands.add('parentsUntilTestId', { prevSubject: 'element' }, (subject: JQuery<HTMLElement>, testId: string) =>
  cy.wrap(subject).parentsUntil(`[uitestid="${testId}"]`).parent()
);
Cypress.Commands.add('toggleSwitch', { prevSubject: 'element' }, (subject: JQuery<HTMLElement>, status: boolean) => {
  const isChecked = subject.find('input').attr('aria-checked') as 'true' | 'false' | undefined;

  if ((status && isChecked !== 'true') || (!status && isChecked === 'true')) {
    cy.wrap(subject).click();
  }
});
Cypress.Commands.add('dtSelect', { prevSubject: 'element' }, (subject: JQuery<HTMLElement>, element: string) => {
  subject.trigger('click');
  cy.get('.dt-select-content dt-option').contains(element).click();
});
