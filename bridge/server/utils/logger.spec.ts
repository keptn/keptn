// eslint-disable-next-line import/no-extraneous-dependencies
import { jest } from '@jest/globals';
import { ComponentLogger, Level, LogDestination, logger as logimpl, LoggerImpl } from './logger';

describe('Logger', () => {
  let logger: LoggerImpl;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let spyLog: any;

  beforeEach(() => {
    logger = new LoggerImpl();
    spyLog = jest.spyOn(console, 'log');
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should use default values', () => {
    // do not log without conf
    logger.error('component', 'message');
    expect(spyLog).toBeCalledTimes(0);

    logger.configure();
    logger.error('component', 'message');
    expect(spyLog).toBeCalledTimes(1);
  });

  it('should select the correct log destination', () => {
    logger.configure(LogDestination.STDOUT);
    logger.error('component', 'message');
    expect(spyLog).toBeCalledTimes(1);

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const fsSpy = jest.spyOn(LoggerImpl.prototype as any, 'fileLog');
    fsSpy.mockImplementation(() => {});
    logger.configure(LogDestination.FILE);
    logger.error('component', 'message');
    expect(fsSpy).toBeCalledTimes(1);
  });

  it('should format correctly the messages', () => {
    logger.configure(LogDestination.STDOUT);
    logger.error('Test1', 'error msg');
    logger.warning('Test2', 'warning msg');
    logger.info('Test3', 'info msg');
    logger.debug('Test4', 'debug msg'); //debug disabled

    expect(spyLog).toBeCalledTimes(3);

    const msg1 = spyLog.mock.calls[0][0];
    const msg2 = spyLog.mock.calls[1][0];
    const msg3 = spyLog.mock.calls[2][0];
    const timestamp = '\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}Z';
    expect(msg1).toMatch(new RegExp(`^\\[Keptn] ${timestamp} error   \\[Test1] error msg`));
    expect(msg2).toMatch(new RegExp(`^\\[Keptn] ${timestamp} warning \\[Test2] warning msg`));
    expect(msg3).toMatch(new RegExp(`^\\[Keptn] ${timestamp} info    \\[Test3] info msg`));
  });

  it('should not log debug level by default', () => {
    logger.configure(LogDestination.STDOUT);
    logger.debug('Test', 'debug msg'); //debug disabled
    expect(spyLog).toBeCalledTimes(0);

    logger.configure(LogDestination.STDOUT, { Test: true });
    logger.debug('Test', 'debug msg');
    expect(spyLog).toBeCalledTimes(1);
  });

  describe('ComponentLogger', () => {
    const logimplSpy = jest.spyOn(logimpl, 'log');

    it('should uses Logger implementation', () => {
      const name = 'Component';
      const msg = 'myError';
      const l = new ComponentLogger(name);
      l.error(msg);
      expect(logimplSpy).toHaveBeenCalledWith(Level.ERROR, name, msg);
      l.warning(msg);
      expect(logimplSpy).toHaveBeenCalledWith(Level.WARNING, name, msg);
      l.info(msg);
      expect(logimplSpy).toHaveBeenCalledWith(Level.INFO, name, msg);
      l.debug(msg);
      expect(logimplSpy).toHaveBeenCalledWith(Level.DEBUG, name, msg);
    });

    it('should pretty print values', () => {
      const obj = {
        key: 'value',
        array: [
          {},
          {
            k: 'v',
            k1: 'v1',
          },
        ],
        k: 'v',
      };
      const l = new ComponentLogger('test');
      expect(l.prettyPrint(obj)).toBe('key=value, array=[object Object],[object Object], k=v');
    });
  });
});
