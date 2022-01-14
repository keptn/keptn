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
  setupFilesAfterEnv: ['<rootDir>/.jest/setupServer.ts'],
  globalTeardown: '<rootDir>/.jest/shutdownServer.ts',
  testPathIgnorePatterns: ['<rootDir>/dist', '<rootDir>/node_modules'],
};
