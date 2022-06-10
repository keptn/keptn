import { AxiosError, Method } from 'axios';
// eslint-disable-next-line import/no-extraneous-dependencies
import { jest } from '@jest/globals';
import { printError } from './print-utils';
import { LogDestination, logger } from './logger';

describe('Test print-utils', () => {
  beforeAll(() => {
    logger.configure(LogDestination.STDOUT);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should print error as axiosError', () => {
    const consoleSpy = jest.spyOn(logger, 'error');
    const error: AxiosError = Object.assign(new Error('myError'), {
      config: {
        url: 'https://keptn.sh',
        method: 'GET' as Method,
      },
      isAxiosError: true,
      toJSON: () => {
        return {};
      },
    });
    printError(error);
    expect(consoleSpy).toBeCalledTimes(1);
    expect(consoleSpy).toHaveBeenCalledWith('API', 'Error for GET https://keptn.sh: myError');
  });

  it('should print request error as axiosError', () => {
    const consoleSpy = jest.spyOn(logger, 'error');
    const error: AxiosError = Object.assign(new Error('myError'), {
      request: {
        path: 'https://keptn.sh',
        method: 'GET',
      },
      config: {},
      isAxiosError: true,
      toJSON: () => {
        return {};
      },
    });
    printError(error);
    expect(consoleSpy).toBeCalledTimes(1);
    expect(consoleSpy).toHaveBeenCalledWith('API', 'Error for GET https://keptn.sh: myError');
  });

  it('should not print as axiosError', () => {
    const consoleSpy = jest.spyOn(logger, 'error');
    const msg = 'myError';
    const error = new Error(msg);
    printError(error);
    expect(consoleSpy).toBeCalledTimes(1);
    expect(consoleSpy).toHaveBeenCalledWith('API', `Error: ${msg}`);
  });
});
