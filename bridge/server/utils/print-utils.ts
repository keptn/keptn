import { AxiosError } from 'axios';
import { ComponentLogger } from './logger';

const log = new ComponentLogger('API');

export function printError(err: AxiosError | Error): void {
  if (isAxiosError(err)) {
    const method = (err.request || err.config).method;
    const url = err.request?.path ?? err.config.url;
    log.error(`Error for ${method} ${url}: ${err.message}`);
  } else {
    const msg = err instanceof Error ? `${err.name}: ${err.message}` : `${err}`;
    log.error(msg);
  }
}

function isAxiosError(err: Error | AxiosError): err is AxiosError {
  return err.hasOwnProperty('isAxiosError');
}
