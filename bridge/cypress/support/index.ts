// eslint-disable-next-line import/no-extraneous-dependencies
import 'cypress-xpath';
import './commands';
import './intercept';

let errorLogs: (string | undefined | Error)[][] = [];
let warningLogs: string[] = [];
Cypress.on('window:before:load', (window) => {
  cy.stub(window.console, `error`).callsFake((...args: (string | undefined | Error)[]) => {
    errorLogs.push(
      args.map((arg) =>
        // instanceof does not work here for Error
        arg?.hasOwnProperty('stack') ? (arg as Error).stack : arg
      )
    );
    console.error(...args); // also write it to the console
  });
  cy.stub(window.console, `warn`).callsFake((...args: (string | undefined)[]) => {
    warningLogs.push(args.join(', '));
    console.warn(...args); // also write it to the console
  });
});

afterEach(() => {
  /*eslint-disable promise/catch-or-return, promise/always-return*/
  if (errorLogs.length) {
    cy.task('log', errorLogs).then(() => {
      expect(errorLogs.length).to.eq(0); // this should be inside then, else log will be aborted
      errorLogs = [];
    });
  }
  if (warningLogs.length) {
    cy.task('log', warningLogs).then(() => {
      expect(errorLogs.length).to.eq(0); // this should be inside then, else log will be aborted
      warningLogs = [];
    });
  }
  /*eslint-enable promise/catch-or-return, promise/always-return*/
});
