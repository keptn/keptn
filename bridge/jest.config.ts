module.exports = {
  rootDir: './',
  testPathIgnorePatterns: [
    '<rootDir>/node_modules/',
    '<rootDir>/server',
    '<rootDir>/cypress',
    '<rootDir>/client/environments',
  ],
  preset: 'jest-preset-angular',
  setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
  globalSetup: 'jest-preset-angular/global-setup',
  collectCoverage: true,
  coverageDirectory: '<rootDir>/coverage',
  moduleNameMapper: {
    '^lodash-es$': 'lodash',
    d3: '<rootDir>/node_modules/d3/dist/d3.min.js',
  },
};
