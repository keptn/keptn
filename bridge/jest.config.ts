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
    '^yaml$': '<rootDir>/node_modules/yaml/dist/index.js',
    '^uuid$': '<rootDir>/node_modules/uuid/dist/index.js',
  },
};
