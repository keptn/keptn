export default {
  preset: 'ts-jest',
  testEnvironment: 'node',
  extensionsToTreatAsEsm: ['.ts'],
  globals: {
    'ts-jest': {
      useESM: true,
    },
  },
  setupFiles: ['./.jest/setEnvVars.js'],
  testPathIgnorePatterns: ['./dist', './node_modules'],
};
