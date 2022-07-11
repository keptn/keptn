import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, Router, UrlTree } from '@angular/router';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Trace } from '../_models/trace';
import { DataService } from '../_services/data.service';
import { ServerErrors } from '../_models/server-error';

@Injectable({
  providedIn: 'root',
})
export class TraceDeepLinkGuard implements CanActivate {
  constructor(private dataService: DataService, private router: Router) {}

  canActivate(route: ActivatedRouteSnapshot): Observable<UrlTree> | UrlTree {
    const keptnContext = route.paramMap.get('keptnContext');
    const eventSelector = route.paramMap.get('eventSelector');
    if (!keptnContext) {
      return this.traceErrorRoute();
    }

    return this.dataService
      .getTracesByContext(keptnContext)
      .pipe(map((traces) => this.navigateToTrace(traces, keptnContext, eventSelector)));
  }

  private navigateToTrace(
    traces: Trace[] | undefined,
    keptnContext: string | null,
    eventSelector: string | null
  ): UrlTree {
    if (!traces?.length) {
      return this.traceErrorRoute(keptnContext);
    }
    if (eventSelector) {
      let trace = this.findTraceForStage(traces, eventSelector);
      if (trace) {
        return this.router.createUrlTree([
          '/project',
          trace.project,
          'sequence',
          trace.shkeptncontext,
          'stage',
          trace.stage,
        ]);
      }
      trace = this.findTraceForEvent(traces, eventSelector);
      if (trace) {
        return this.router.createUrlTree([
          '/project',
          trace.project,
          'sequence',
          trace.shkeptncontext,
          'event',
          trace.id,
        ]);
      }
      return this.traceErrorRoute(keptnContext);
    } else {
      const trace = this.findTraceForKeptnContext(traces);
      if (trace) {
        return this.router.createUrlTree(['/project', trace.project, 'sequence', trace.shkeptncontext]);
      }
      return this.traceErrorRoute(keptnContext);
    }
  }

  private traceErrorRoute(keptnContext?: string | null): UrlTree {
    return this.router.createUrlTree(['error'], {
      queryParams: {
        status: ServerErrors.TRACE,
        ...(keptnContext && { keptnContext }),
      },
    });
  }

  private findTraceForKeptnContext(traces: Trace[]): Trace | undefined {
    return traces.find((t: Trace) => !!t.project && !!t.service);
  }

  private findTraceForStage(traces: Trace[], eventselector: string | null): Trace | undefined {
    return traces.find((t: Trace) => t.data.stage === eventselector && !!t.project && !!t.service);
  }

  private findTraceForEvent(traces: Trace[], eventselector: string | null): Trace | undefined {
    return [...traces].reverse().find((t: Trace) => t.type === eventselector && !!t.project && !!t.service);
  }
}
