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

beforeEach(() => {
  errorLogs = [];
  warningLogs = [];

  cy.intercept(/api\/(.)*/, (request) => {
    request.continue((response) => {
      if (response.statusCode < 200 || response.statusCode > 399) {
        errorLogs.push([`Request to "${request.url}" failed`, `payload: ${JSON.stringify(request.body)}`]);
      }
      expect(response.statusCode).to.be.gte(200).and.to.be.lte(399);
    });
  });
});

afterEach(() => {
  /*eslint-disable promise/catch-or-return, promise/always-return*/
  cy.task('logError', [...errorLogs])
    .then((expectedErrors) => {
      expect(errorLogs.length).to.eq(expectedErrors);
    })
    .task('logWarning', [...warningLogs])
    .then((expectedWarnings) => {
      expect(warningLogs.length).to.eq(expectedWarnings);
    });
  /*eslint-enable promise/catch-or-return, promise/always-return*/
});
