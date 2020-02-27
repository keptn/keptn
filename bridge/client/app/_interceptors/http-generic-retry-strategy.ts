import { timer, throwError, Observable } from 'rxjs';
import { mergeMap } from 'rxjs/operators';

export interface RetryParams {
  maxAttempts?: number;
  scalingDuration?: number;
  shouldRetry?: ({ status: number }) => boolean;
}

const defaultParams: RetryParams = {
  maxAttempts: 3,
  scalingDuration: 1000,
  shouldRetry: ({ status }) => status >= 400
}

export const genericRetryStrategy = (params: RetryParams = {}) => (attempts: Observable<any>) => attempts.pipe(
  mergeMap((error, i) => {
    const { maxAttempts, scalingDuration, shouldRetry } = { ...defaultParams, ...params };
    const retryAttempt = i + 1;
    if (retryAttempt > maxAttempts || !shouldRetry(error))
      return throwError(error);
    return timer(retryAttempt * scalingDuration);
  })
);
