import { timer, throwError, Observable } from 'rxjs';
import { mergeMap } from 'rxjs/operators';

export interface RetryParams {
  maxAttempts?: number;
  scalingDuration?: number;
  shouldRetry?: ({ status: number }) => boolean;
}

/**
 * Retry for either client or server errors.
 *
 * Server purposefully send 401 OR 403 based on login flow enabled for this Keptn instance.
 *
 *  - 401 : Unauthorized. Login required
 *  - 403 : Forbidden. Login was successful, however, it's not acceptable for this instance.
 */
const avoidRetryFor = [401, 403];

const defaultParams: RetryParams = {
  maxAttempts: 3,
  scalingDuration: 1000,
  shouldRetry: ({status}) => (status >= 400 && !avoidRetryFor.includes(status))
};

export const genericRetryStrategy = (params: RetryParams = {}) => (attempts: Observable<any>) => attempts.pipe(
  mergeMap((error, i) => {
    const { maxAttempts, scalingDuration, shouldRetry } = { ...defaultParams, ...params };
    const retryAttempt = i + 1;
    if (retryAttempt > maxAttempts || !shouldRetry(error))
      return throwError(error);
    return timer(retryAttempt * scalingDuration);
  })
);
