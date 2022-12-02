import { AxiosError } from 'axios';
import { ComponentLogger } from './logger';
import { IncomingHttpHeaders } from 'http';

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

export function filterHeaders(headers: IncomingHttpHeaders): IncomingHttpHeaders {
  const filteredHeaders = { ...headers };
  const denyList = ['cookies', 'user-agent'];
  for (const item of denyList) {
    // it is safe to delete because all properties are marked as optional ?
    delete filteredHeaders[item];
  }
  // remove additional sec headers
  Object.keys(headers).forEach((key) => {
    if (key.startsWith('sec')) {
      delete filteredHeaders[key];
    }
  });
  return filteredHeaders;
}
