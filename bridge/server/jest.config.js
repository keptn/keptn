export default {
  rootDir: './',
  preset: 'ts-jest',
  testEnvironment: 'node',
  extensionsToTreatAsEsm: ['.ts'],
  globals: {
    'ts-jest': {
      useESM: true,
      tsconfig: './tsconfig.spec.json',
    },
  },
  collectCoverage: true,
  coverageDirectory: '<rootDir>/coverage',
  setupFiles: ['<rootDir>/.jest/setEnvVars.ts'],
  testPathIgnorePatterns: ['<rootDir>/dist', '<rootDir>/node_modules'],
};
