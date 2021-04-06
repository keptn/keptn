import {Injectable} from '@angular/core';
import {HttpErrorResponse, HttpEvent, HttpHandler, HttpInterceptor, HttpRequest} from '@angular/common/http';
import {Observable, of, throwError} from 'rxjs';
import {catchError, retryWhen} from 'rxjs/operators';
import {genericRetryStrategy} from './http-generic-retry-strategy';
import {DtToast} from '@dynatrace/barista-components/toast';
import {Router} from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class HttpErrorInterceptor implements HttpInterceptor {
  private isReloading = false;

  constructor(private readonly toast: DtToast,
              private readonly router: Router) {
  }

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    return next.handle(request)
      .pipe(
        retryWhen(genericRetryStrategy()),
        catchError((error: HttpErrorResponse) => {

          if (error.status === 401) {
            if (!this.isReloading) {
              this.isReloading = true;
              this.toast.create('Login required. Redirecting to login.');
              setTimeout(() => window.location.href = '/login', 1000);
            }

            return of(undefined);
          }

          if (error.error instanceof ErrorEvent) {
            // A client-side or network error occurred. Handle it accordingly.
            console.error('An error occurred:', error.error.message);
            this.toast.create(`${error.error.message}`);
          } else {
            // The backend returned an unsuccessful response code.
            // The response body may contain clues as to what went wrong,
            console.error(`${error.status} ${error.message}`);
            this.toast.create(`${error.message}`);
          }

          return throwError(error);
        })
      );
  }
}
