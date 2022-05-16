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
      /**
       * Toggles a dt-switch and checks a dt-checkbox
       * @param status
       */
      dtCheck<E extends Node = HTMLElement>(status: boolean): Cypress.Chainable<JQuery<E>>;
      dtSelect<E extends Node = HTMLElement>(element: string): Cypress.Chainable<JQuery<E>>;
      dtQuickFilterCheck<E extends Node = HTMLElement>(
        filterName: string,
        itemName: string,
        status: boolean
      ): Cypress.Chainable<JQuery<E>>;
      clearDtFilter<E extends Node = HTMLElement>(): Cypress.Chainable<JQuery<E>>;
      clickOutside<E extends Node = HTMLElement>(): Cypress.Chainable<JQuery<E>>;
    }
  }
}
// Commands have to be added by hooking them to Cypress
Cypress.Commands.add(
  'byTestId',
  { prevSubject: ['optional', 'element'] },
  (sbj: JQuery<HTMLElement> | void, testId: string) => {
    const selector = `[uitestid="${testId}"]`;
    if (sbj) {
      return cy.wrap(sbj).find(selector);
    }
    return cy.get(selector);
  }
);
Cypress.Commands.add('clickOutside', () => cy.get('body').click(0, 0));
Cypress.Commands.add('parentsUntilTestId', { prevSubject: 'element' }, (subject: JQuery<HTMLElement>, testId: string) =>
  cy.wrap(subject).parentsUntil(`[uitestid="${testId}"]`).parent()
);
Cypress.Commands.add('dtCheck', { prevSubject: 'element' }, (subject: JQuery<HTMLElement>, status: boolean) => {
  const isChecked = subject.find('input').attr('aria-checked') as 'true' | 'false' | undefined;

  if ((status && isChecked !== 'true') || (!status && isChecked === 'true')) {
    cy.wrap(subject).click();
  }
});

Cypress.Commands.add('dtSelect', { prevSubject: 'element' }, (subject: JQuery<HTMLElement>, element: string) => {
  cy.wrap(subject).click();
  cy.get('.dt-select-content dt-option').contains(element).click();
});

Cypress.Commands.add(
  'dtQuickFilterCheck',
  { prevSubject: 'element' },
  (subject: JQuery<HTMLElement>, filterName: string, itemName: string, status: boolean) => {
    cy.wrap(subject)
      .find('.dt-quick-filter-group .dt-quick-filter-group-headline')
      .contains(filterName)
      .parent()
      .find('dt-checkbox')
      .contains(itemName)
      .dtCheck(status);
  }
);

Cypress.Commands.add('clearDtFilter', { prevSubject: 'element' }, (subject: JQuery<HTMLElement>) => {
  subject.find('.dt-filter-field-clear-all-button').trigger('click');
  cy.wrap(subject).find('.dt-filter-field-input ').type('{esc}');
});
