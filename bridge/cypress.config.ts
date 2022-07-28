// eslint-disable-next-line import/no-extraneous-dependencies
import { defineConfig } from 'cypress';

let expectedErrorCount = 0;
let expectedWarningCount = 0;

export default defineConfig({
  videosFolder: 'dist/cypress/videos',
  screenshotsFolder: 'dist/cypress/screenshots',
  fixturesFolder: 'cypress/fixtures',
  video: false,
  screenshotOnRunFailure: true,
  viewportHeight: 1080,
  viewportWidth: 1920,
  chromeWebSecurity: false,
  retries: {
    // Configure retry attempts for `cypress run`
    runMode: 2,
    // Configure retry attempts for `cypress open`
    openMode: 0,
  },
  e2e: {
    baseUrl: 'http://localhost:5000', // workaround until https://github.com/cypress-io/cypress/issues/21555 is fixed
    specPattern: 'cypress/integration/**/**.spec.{js,jsx,ts,tsx}',
    setupNodeEvents(on) {
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
            console.log('Received errors:', logs);
          }
          expectedErrorCount = 0;
          return expected;
        },
        logWarning(logs: string[]): number {
          const expected = expectedWarningCount;
          if (logs.length !== expected) {
            console.log('Received warnings:', logs);
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
    },
  },
});
