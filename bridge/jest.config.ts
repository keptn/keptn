module.exports = {
  rootDir: './',
  testPathIgnorePatterns: ['<rootDir>/node_modules/', '<rootDir>/cypress'],
  preset: 'jest-preset-angular',
  setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
  collectCoverage: true,
  coverageDirectory: '<rootDir>/coverage',
  moduleNameMapper: {
    '^lodash-es$': 'lodash',
  },
};
