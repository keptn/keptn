export default {
  rootDir: './',
  preset: 'ts-jest',
  testEnvironment: 'node',
  extensionsToTreatAsEsm: ['.ts'],
  transform: {
    '^.+\\.ts$': [
      'ts-jest/legacy',
      {
        useESM: true,
        tsconfig: './tsconfig.spec.json',
      },
    ],
  },
  collectCoverage: true,
  coverageDirectory: '<rootDir>/coverage',
  testPathIgnorePatterns: ['<rootDir>/dist', '<rootDir>/node_modules'],
};
