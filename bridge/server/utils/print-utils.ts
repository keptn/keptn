import { AxiosError } from 'axios';

export function printError(err: AxiosError | Error): void {
  if (isAxiosError(err)) {
    const method = (err.request || err.config).method;
    const url = err.request?.path ?? err.config.url;
    console.error(`Error for ${method} ${url}: ${err.message}`);
  } else {
    console.error(err);
  }
}

function isAxiosError(err: Error | AxiosError): err is AxiosError {
  return err.hasOwnProperty('isAxiosError');
}
