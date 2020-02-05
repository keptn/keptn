import { Injectable } from '@angular/core';
import {
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
  HttpResponse
} from "@angular/common/http";
import {Observable, of} from "rxjs";
import {map} from "rxjs/operators";

const evaluationLabelsMockData = {
  "testid": "12345",
  "buildnr": "build17",
  "runby": "JohnDoe"
};

@Injectable({
  providedIn: 'root'
})
export class HttpMockInterceptor implements HttpInterceptor {

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    return next.handle(request).pipe(
      map((event: HttpEvent<any>) => {
        if (event instanceof HttpResponse) {
          if(request.url.indexOf('/api/traces/') != -1) {
            if(event.body) {
              let traces = event.body;
              traces
                .map((trace) => {
                  if(trace.type == 'sh.keptn.events.evaluation-done') {
                    trace.data.labels = evaluationLabelsMockData;
                  }
                  return trace;
                });
            }
          }
        }
        return event;
      }));
  }
}
