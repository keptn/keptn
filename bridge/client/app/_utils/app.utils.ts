import { Observable, of, timer } from 'rxjs';

export class AppUtils {
  public static createTimer(delay: number): Observable<unknown> {
    if (delay === 0) {
      return of(undefined);
    }

    return timer(0, delay);
  }
}
