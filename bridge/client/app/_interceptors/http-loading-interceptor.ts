import {Injectable} from "@angular/core";
import {HttpEvent, HttpHandler, HttpInterceptor, HttpRequest} from "@angular/common/http";
import {Observable} from "rxjs";
import {delay, finalize, retryWhen, take, tap} from "rxjs/operators";

import {HttpProgressState} from "../_models/http-progress-state";
import {HttpStateService} from "../_services/http-state.service";
import {environment} from "../../environments/environment";


@Injectable({
  providedIn: 'root'
})
export class HttpLoadingInterceptor implements HttpInterceptor {

  constructor(private httpStateService: HttpStateService) { }

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {

    this.httpStateService.state.next({
      url: request.url,
      state: HttpProgressState.start
    });

    return next.handle(request).pipe(finalize(() => {
      this.httpStateService.state.next({
        url: request.url,
        state: HttpProgressState.end
      });
    }));
  }
}
