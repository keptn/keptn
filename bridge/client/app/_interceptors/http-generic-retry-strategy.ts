import { Observable, throwError, timer } from 'rxjs';
import { mergeMap } from 'rxjs/operators';
import { HttpErrorResponse } from '@angular/common/http';

export interface RetryParams {
  maxAttempts: number;
  scalingDuration: number;
  shouldRetry: ({ status }: { status: number }) => boolean;
}

/**
 * Retry for the following status codes
 *  - 500: Internal Server Error
 *  - 502: Bad Gateway
 *  - 503: Service Unavailable
 *  - 504: Gateway Timeout
 */
const retryFor = [500, 502, 503, 504];

const defaultParams: RetryParams = {
  maxAttempts: 3,
  scalingDuration: 1000,
  shouldRetry: ({ status }) => retryFor.includes(status),
};

export const genericRetryStrategy =
  (params?: RetryParams) =>
  (attempts: Observable<HttpErrorResponse>): Observable<number> =>
    attempts.pipe(
      mergeMap((error, i) => {
        const { maxAttempts, scalingDuration, shouldRetry } = { ...defaultParams, ...params };
        const retryAttempt = i + 1;

        if (retryAttempt > maxAttempts || !shouldRetry(error)) {
          return throwError(error);
        }
        return timer(retryAttempt * scalingDuration);
      })
    );
