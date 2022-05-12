import { AxiosError, Method } from 'axios';
// eslint-disable-next-line import/no-extraneous-dependencies
import { jest } from '@jest/globals';
import { printError } from './print-utils';

describe('Test print-utils', () => {
  it('should print error as axiosError', () => {
    const consoleSpy = jest.spyOn(console, 'error');
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
    expect(consoleSpy).toHaveBeenCalledWith('Error for GET https://keptn.sh: myError');
  });

  it('should print request error as axiosError', () => {
    const consoleSpy = jest.spyOn(console, 'error');
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
    expect(consoleSpy).toHaveBeenCalledWith('Error for GET https://keptn.sh: myError');
  });

  it('should not print as axiosError', () => {
    const consoleSpy = jest.spyOn(console, 'error');
    const error = new Error('myError');
    printError(error);
    expect(consoleSpy).toHaveBeenCalledWith(error);
  });
});
