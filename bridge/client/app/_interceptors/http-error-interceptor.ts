import { Inject, Injectable, OnDestroy } from '@angular/core';
import { HttpErrorResponse, HttpEvent, HttpHandler, HttpInterceptor, HttpRequest } from '@angular/common/http';
import { EMPTY, NEVER, Observable, Subject, throwError } from 'rxjs';
import { catchError, retryWhen, takeUntil } from 'rxjs/operators';
import { genericRetryStrategy, RetryParams } from './http-generic-retry-strategy';
import { Location } from '@angular/common';
import { RETRY_ON_HTTP_ERROR } from '../_utils/app.utils';
import { NotificationsService } from '../_services/notifications.service';
import { NotificationType } from '../_models/notification';
import { AuthType } from '../../../shared/models/auth-type';
import { DataService } from '../_services/data.service';

@Injectable({
  providedIn: 'root',
})
export class HttpErrorInterceptor implements HttpInterceptor, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();

  private isReloading = false;
  private isAuthorizedErrorShown = false;
  private keptnInfo$ = this.dataService.keptnInfo;

  constructor(
    private readonly location: Location,
    private readonly notificationService: NotificationsService,
    private readonly dataService: DataService,
    @Inject(RETRY_ON_HTTP_ERROR) private hasRetry: boolean
  ) {}

  intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    const params: RetryParams | undefined = this.hasRetry
      ? undefined
      : { maxAttempts: 3, scalingDuration: 0, shouldRetry: (): boolean => false };
    return next.handle(request).pipe(
      retryWhen(genericRetryStrategy(params)),
      catchError((response: HttpErrorResponse) => {
        return this.handleError(response);
      })
    );
  }

  private handleError(response: HttpErrorResponse): Observable<HttpEvent<unknown>> {
    if (response.status === 401) {
      return this.handleUnauthorizedError(response);
    }

    if (response.status === 403) {
      this.keptnInfo$.pipe(takeUntil(this.unsubscribe$)).subscribe((keptnInfo) => {
        this.notificationService.addNotification(
          NotificationType.ERROR,
          `${
            keptnInfo?.bridgeInfo.user ? keptnInfo?.bridgeInfo.user : 'User'
          } does not have the permissions to perform this action.`
        );
      });
      return throwError(() => response);
    }

    if (response.status === 409 && response.url?.endsWith('/api/secrets/v1/secret')) {
      // Special case for already existing secrets - for unit test has to be before instanceof ErrorEvent
      return throwError(() => response);
    }

    if (typeof response.error === 'string') {
      // The backend returned an unsuccessful response code.
      // The response body may contain clues as to what went wrong,
      this.notificationService.addNotification(NotificationType.ERROR, response.error);
      return throwError(() => response);
    }

    // A client-side or network error occurred. Handle it accordingly.
    this.notificationService.addNotification(NotificationType.ERROR, response.message);
    return throwError(() => response);
  }

  private handleUnauthorizedError(error: HttpErrorResponse): Observable<HttpEvent<unknown>> {
    const authType = error.headers.get('keptn-auth-type') as AuthType | null;

    if (this.isAuthorizedErrorShown) {
      return EMPTY;
    }

    if (this.isReloading) {
      return NEVER;
    }

    if (authType === AuthType.OAUTH) {
      this.isReloading = true;
      this.notificationService.addNotification(NotificationType.INFO, 'Login required. Redirecting to login.');
      // Wait for few moments to let user see the toast message and navigate to external login route
      setTimeout(() => (window.location.href = this.location.prepareExternalUrl('/oauth/login')), 1000);
      return NEVER;
    }

    if (authType === AuthType.BASIC) {
      this.isAuthorizedErrorShown = true;
      this.notificationService.addNotification(
        NotificationType.ERROR,
        'Login credentials invalid. Please check your provided username and password.'
      );
      return throwError(() => error);
    }

    let errorInfo;
    if (error.error === 'incorrect api key auth') {
      errorInfo = 'Could not authorize API token. Please check the configured API token.';
    } else {
      errorInfo = 'Could not authorize.';
    }
    this.isAuthorizedErrorShown = true;
    this.notificationService.addNotification(NotificationType.ERROR, errorInfo);

    return throwError(() => error);
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
