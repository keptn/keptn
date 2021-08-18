import { Observable, of, timer } from 'rxjs';
import { InjectionToken } from '@angular/core';

export class AppUtils {
  public static createTimer(delay: number, dueTime: number = 0): Observable<unknown> {
    if (delay === 0) {
      return of(undefined);
    }

    return timer(dueTime, delay);
  }
}

export const INITIAL_DELAY_MILLIS = new InjectionToken<number>('Initial delay in millis');
export const RETRY_ON_HTTP_ERROR = new InjectionToken<boolean>('If retry is turned on for interceptor');
