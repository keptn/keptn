import { Observable, of, timer } from 'rxjs';
import { InjectionToken } from '@angular/core';

export class AppUtils {
  public static createTimer(delay: number, dueTime: number): Observable<unknown> {
    if (dueTime === 0) {
      return of(undefined);
    }

    return timer(delay, dueTime);
  }

  public static copyObject<T>(data: T): T {
    return JSON.parse(JSON.stringify(data));
  }

  public static round(value: number, places: number): number {
    return +(Math.round(Number(`${value}e+${places}`)) + `e-${places}`);
  }

  public static formatNumber(value: number): number {
    let n = Math.abs(value);
    if (n < 1) {
      n = Math.floor(n * 1000) / 1000;
    } else if (n < 100) {
      n = Math.floor(n * 100) / 100;
    } else if (n < 1000) {
      n = Math.floor(n * 10) / 10;
    } else {
      n = Math.floor(n);
    }

    return n;
  }
}

export const POLLING_INTERVAL_MILLIS = new InjectionToken<number>('Polling interval in millis');
export const RETRY_ON_HTTP_ERROR = new InjectionToken<boolean>('If retry is turned on for interceptor');
