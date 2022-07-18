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
    return data === undefined ? undefined : JSON.parse(JSON.stringify(data));
  }

  public static truncateNumber(value: number, decimals: number): number {
    return Math.trunc(value * Math.pow(10, decimals)) / Math.pow(10, decimals);
  }

  public static round(value: number, places: number): number {
    return +(Math.round(Number(`${value}e+${places}`)) + `e-${places}`);
  }

  public static formatNumber(value: number): number {
    const abs = Math.abs(value);
    let n = value;
    if (abs < 1) {
      n = Math.trunc(n * 1000) / 1000;
    } else if (abs < 100) {
      n = Math.trunc(n * 100) / 100;
    } else if (abs < 1000) {
      n = Math.trunc(n * 10) / 10;
    } else {
      n = Math.trunc(n);
    }

    return n;
  }

  public static isValidJson(value: string): boolean {
    try {
      JSON.parse(value);
      return true;
    } catch (e) {
      return false;
    }
  }

  public static isValidUrl(value: string): boolean {
    try {
      new URL(value);
    } catch (_) {
      return false;
    }
    return true;
  }

  public static splitURLPort(url: string): { host: string; port: string } {
    const index = url.lastIndexOf(':');
    let host = url;
    let port = '';

    if (index !== -1) {
      const portSubstr = url.substring(index + 1);
      if (!isNaN(+portSubstr)) {
        host = url.substring(0, index);
        port = portSubstr;
      }
    }
    return {
      host,
      port,
    };
  }
}

export const POLLING_INTERVAL_MILLIS = new InjectionToken<number>('Polling interval in millis');
export const RETRY_ON_HTTP_ERROR = new InjectionToken<boolean>('If retry is turned on for interceptor');
