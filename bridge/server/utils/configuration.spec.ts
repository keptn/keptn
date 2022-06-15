import { getConfiguration } from './configuration';
import { LogDestination } from './logger';

describe('Configuration', () => {
  afterEach(() => {
    const envToClear = ['LOGGING_COMPONENTS'];
    envToClear.forEach((name) => delete process.env[name]);
  });

  it('should use default values', () => {
    const result = getConfiguration();
    expect(result.logging.destination).toBe(LogDestination.STDOUT);
    expect(result.logging.enabledComponents).toStrictEqual({});
  });

  it('should set values using options object', () => {
    let result = getConfiguration({
      logging: {
        enabledComponents: 'a=true,b=false,c=true',
      },
    });
    expect(result.logging.destination).toBe(LogDestination.STDOUT);
    expect(result.logging.enabledComponents).toStrictEqual({
      a: true,
      b: false,
      c: true,
    });

    result = getConfiguration({
      logging: {
        enabledComponents: 'a=false',
        destination: LogDestination.FILE,
      },
    });
    expect(result.logging.destination).toBe(LogDestination.FILE);
    expect(result.logging.enabledComponents).toStrictEqual({
      a: false,
    });
  });

  it('should set values using env var', () => {
    process.env.LOGGING_COMPONENTS = 'a=true,b=false,c=true';
    const result = getConfiguration();
    expect(result.logging.destination).toBe(LogDestination.STDOUT);
    expect(result.logging.enabledComponents).toStrictEqual({
      a: true,
      b: false,
      c: true,
    });
  });

  it('option object should win over env var', () => {
    process.env.LOGGING_COMPONENTS = 'a=false,b=true,c=false';
    const result = getConfiguration({
      logging: {
        enabledComponents: 'a=true,b=false,c=true',
      },
    });
    expect(result.logging.destination).toBe(LogDestination.STDOUT);
    expect(result.logging.enabledComponents).toStrictEqual({
      a: true,
      b: false,
      c: true,
    });
  });

  it('should correctly eval booleans', () => {
    const result = getConfiguration({
      logging: {
        enabledComponents: 'a=tRue,b=FaLsE,c=0,d=1,e=enabled',
      },
    });
    expect(result.logging.enabledComponents).toStrictEqual({
      a: true,
      b: false,
      c: false,
      d: true,
      e: true,
    });
  });
});
