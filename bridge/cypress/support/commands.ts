// ***********************************************
// This example namespace declaration will help
// with Intellisense and code completion in your
// IDE or Text Editor.
// ***********************************************
// tslint:disable-next-line:no-namespace
declare namespace Cypress {
  // tslint:disable-next-line:no-any
  interface Chainable<Subject = any> {
    byTestId<E extends Node = HTMLElement>(id: string):
      Cypress.Chainable<JQuery<E>>;
  }
}

// Actual function
const byTestId = (testId: string) =>
  cy.get(`[uitestid="${testId}"]`);

// Hooking into Cypress
Cypress.Commands.add('byTestId', byTestId);
//
// function customCommand(param: any): void {
//   console.warn(param);
// }
//
// NOTE: You can use it like so:
// Cypress.Commands.add('customCommand', customCommand);
//
// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add("login", (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add("drag", { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add("dismiss", { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite("visit", (originalFn, url, options) => { ... })
