import { getConfiguration } from './configuration';
import { LogDestination } from './logger';

describe('Configuration', () => {
  it('Default values', () => {
    const result = getConfiguration();
    expect(result.logging.destination).toBe(LogDestination.StdOut);
    expect(result.logging.enabledComponents).toStrictEqual({});
  });

  it('set values using options object', () => {

  });

  it('set values using env var', () => {

  });

});
