// eslint-disable-next-line import/no-extraneous-dependencies
import 'cypress-xpath';
import './commands';
import './intercept';

Cypress.on('window:before:load', (window) => {
  cy.spy(window.console, 'error');
  cy.spy(window.console, 'warn');
});

afterEach(() => {
  // eslint-disable-next-line promise/always-return,promise/catch-or-return
  cy.window().then((window) => {
    const errorCount = window.errorCount ?? 0;
    window.errorCount = 0;
    expect(window.console.error).to.have.callCount(errorCount);
    expect(window.console.warn).to.have.callCount(0);
  });
});
