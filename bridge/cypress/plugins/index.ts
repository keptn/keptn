// Plugins enable you to tap into, modify, or extend the internal behavior of Cypress
// For more info, visit https://on.cypress.io/plugins-api

let expectedErrorCount = 0;
let expectedWarningCount = 0;
// eslint-disable-next-line @typescript-eslint/no-unused-vars
export default (on: Cypress.PluginEvents, config: Cypress.PluginConfigOptions): void => {
  on('before:browser:launch', (browser, launchOptions) => {
    expectedErrorCount = 0;
    expectedWarningCount = 0;

    if (browser.name === 'chrome' && browser.isHeadless) {
      launchOptions.args.push('--window-size=1920,1080');
      launchOptions.args.push('--force-device-scale-factor=1');
    }
  });
  on('task', {
    logError(logs: (string | undefined | Error)[][]): number {
      const expected = expectedErrorCount;
      if (logs.length !== expected) {
        console.log(logs);
      }
      expectedErrorCount = 0;
      return expected;
    },
    logWarning(logs: string[]): number {
      const expected = expectedWarningCount;
      if (logs.length !== expected) {
        console.log(logs);
      }
      expectedWarningCount = 0;
      return expected;
    },
    setExpectedErrorCount(count: number): null {
      expectedErrorCount = count;
      return null;
    },
    setExpectedWarningCount(count: number): null {
      expectedWarningCount = count;
      return null;
    },
  });
};
