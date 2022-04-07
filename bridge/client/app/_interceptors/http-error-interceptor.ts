import { Inject, Injectable } from '@angular/core';
import { HttpErrorResponse, HttpEvent, HttpHandler, HttpInterceptor, HttpRequest } from '@angular/common/http';
import { Observable, throwError, of } from 'rxjs';
import { catchError, retryWhen } from 'rxjs/operators';
import { genericRetryStrategy, RetryParams } from './http-generic-retry-strategy';
import { Location } from '@angular/common';
import { RETRY_ON_HTTP_ERROR } from '../_utils/app.utils';
import { NotificationsService } from '../_services/notifications.service';
import { NotificationType } from '../_models/notification';

@Injectable({
  providedIn: 'root',
})
export class HttpErrorInterceptor implements HttpInterceptor {
  private isReloading = false;
  private isAuthorizedErrorShown = false;

  constructor(
    private readonly location: Location,
    private readonly notificationService: NotificationsService,
    @Inject(RETRY_ON_HTTP_ERROR) private hasRetry: boolean
  ) {}

  intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    const params: RetryParams | undefined = this.hasRetry
      ? undefined
      : { maxAttempts: 3, scalingDuration: 0, shouldRetry: (): boolean => false };
    return next.handle(request).pipe(
      retryWhen(genericRetryStrategy(params)),
      catchError((error) => {
        if (error.status === 401) {
          this._handleUnauthorizedError(error);
          return of(error);
        } else if (error.status === 403) {
          this.notificationService.addNotification(
            NotificationType.ERROR,
            'You do not have the permissions to perform this action.'
          );
          return of(error);
        } else if (error.status === 409 && error.url?.endsWith('/api/secrets/v1/secret')) {
          // Special case for already existing secrets - for unit test has to be before instanceof ErrorEvent
          return throwError(error);
        } else if (error.error instanceof ErrorEvent) {
          // A client-side or network error occurred. Handle it accordingly.
          this.notificationService.addNotification(NotificationType.ERROR, error.error.message);
        } else {
          // The backend returned an unsuccessful response code.
          // The response body may contain clues as to what went wrong,
          this.notificationService.addNotification(NotificationType.ERROR, error.message);
        }

        return throwError(error);
      })
    );
  }

  private _handleUnauthorizedError(error: HttpErrorResponse): void {
    const authType = error.headers.get('keptn-auth-type');
    if (authType === 'OAUTH') {
      if (!this.isReloading) {
        if (window.location.href.indexOf('/oauth/login') === -1) {
          this.isReloading = true;
          this.notificationService.addNotification(NotificationType.INFO, 'Login required. Redirecting to login.');
          // Wait for few moments to let user see the toast message and navigate to external login route
          setTimeout(() => (window.location.href = this.location.prepareExternalUrl('/oauth/login')), 1000);
        }
      }
    } else if (authType === 'BASIC') {
      if (!this.isAuthorizedErrorShown) {
        this.isAuthorizedErrorShown = true;
        this.notificationService.addNotification(
          NotificationType.ERROR,
          'Login credentials invalid. Please check your provided username and password.'
        );
      }
    } else if (!this.isAuthorizedErrorShown) {
      let errorInfo;
      if (error.error === 'incorrect api key auth') {
        errorInfo = 'Could not authorize API token. Please check the configured API token.';
      } else {
        errorInfo = 'Could not authorize.';
      }
      this.isAuthorizedErrorShown = true;
      this.notificationService.addNotification(NotificationType.ERROR, ' ' + errorInfo);
    }
  }
}
