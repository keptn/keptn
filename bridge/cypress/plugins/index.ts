// Plugins enable you to tap into, modify, or extend the internal behavior of Cypress
// For more info, visit https://on.cypress.io/plugins-api

export default (on: Cypress.PluginEvents, config: Cypress.PluginConfigOptions) => {
  on('before:browser:launch', (browser, launchOptions) => {
    if (browser.name === 'chrome' && browser.isHeadless) {
      launchOptions.args.push('--window-size=1920,1080');
      launchOptions.args.push('--force-device-scale-factor=1');
    }
  });
};
