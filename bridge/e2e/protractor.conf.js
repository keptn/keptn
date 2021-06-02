// @ts-check
// Protractor configuration file, see link for more information
// https://github.com/angular/protractor/blob/master/lib/config.ts

const { SpecReporter } = require('jasmine-spec-reporter');

let chromeExtensions = [];

// add basic auth to chrome
if (process.env.KEPTN_BRIDGE_USER && process.env.KEPTN_BRIDGE_PASSWORD) {
  const { Authenticator } = require('authenticator-browser-extension');

  chromeExtensions.push(Authenticator.for(process.env.KEPTN_BRIDGE_USER || null, process.env.KEPTN_BRIDGE_PASSWORD || null).asBase64());
}

/**
 * @type { import("protractor").Config }
 */
exports.config = {
  allScriptsTimeout: 11000,
  specs: [
    './src/**/*.e2e-spec.ts'
  ],
  capabilities: {
    browserName: 'chrome',

    chromeOptions: {
      extensions: chromeExtensions
    }
  },
  directConnect: true,
  baseUrl: 'http://localhost:4200/',
  framework: 'jasmine',
  jasmineNodeOpts: {
    showColors: true,
    defaultTimeoutInterval: 30000,
    print: function() {}
  },
  onPrepare() {
    require('ts-node').register({
      project: require('path').join(__dirname, './tsconfig.json')
    });
    jasmine.getEnv().addReporter(new SpecReporter({ spec: { displayStacktrace: true } }));
  }
};
