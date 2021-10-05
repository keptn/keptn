// ***********************************************************
// This example support/index.js is processed and
// loaded automatically before your test files.
//
// This is a great place to put global configuration and
// behavior that modifies Cypress.
//
// You can change the location of this file or turn off
// automatically serving support files with the
// 'supportFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/configuration
// ***********************************************************

// When a command from ./commands is ready to use, import with `import './commands'` syntax
// import './commands';

// eslint-disable-next-line import/no-extraneous-dependencies
import 'cypress-xpath';

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
