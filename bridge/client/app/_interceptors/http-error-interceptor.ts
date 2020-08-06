import { Injectable } from '@angular/core';
import {HttpErrorResponse, HttpEvent, HttpHandler, HttpInterceptor, HttpRequest} from "@angular/common/http";
import {Observable, throwError} from "rxjs";
import {catchError, retryWhen} from "rxjs/operators";
import {genericRetryStrategy} from "./http-generic-retry-strategy";
import {DtToast} from "@dynatrace/barista-components/toast";

@Injectable({
  providedIn: 'root'
})
export class HttpErrorInterceptor implements HttpInterceptor {

  constructor(private readonly toast: DtToast){
  }

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    return next.handle(request)
      .pipe(
        retryWhen(genericRetryStrategy()),
        catchError((error: HttpErrorResponse) => {
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
