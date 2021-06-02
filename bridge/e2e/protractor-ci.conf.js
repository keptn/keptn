const config = require('./protractor.conf').config;

config.capabilities = {
  browserName: 'chrome',
  'goog:chromeOptions': {
    args: ['--headless', '--disable-gpu', '--window-size=1920,1080', '--no-sandbox', "--disable-dev-shm-usage"]
  }
};
exports.config = config;
