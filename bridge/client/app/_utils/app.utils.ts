import { Observable, of, timer } from 'rxjs';
import { InjectionToken } from '@angular/core';

export class AppUtils {
  public static createTimer(delay: number, dueTime: number): Observable<unknown> {
    if (delay === 0) {
      return of(undefined);
    }

    return timer(delay, dueTime);
  }
}

export const POLLING_INTERVAL_MILLIS = new InjectionToken<number>('Polling interval in millis');
export const RETRY_ON_HTTP_ERROR = new InjectionToken<boolean>('If retry is turned on for interceptor');
