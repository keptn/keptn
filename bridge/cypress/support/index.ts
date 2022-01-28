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
    expect(window.console.error).to.have.callCount(0);
    expect(window.console.warn).to.have.callCount(0);
  });
});
