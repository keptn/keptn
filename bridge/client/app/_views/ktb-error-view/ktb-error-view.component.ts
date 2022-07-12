import { Component, Input } from '@angular/core';
import { environment } from '../../../environments/environment';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { map } from 'rxjs/operators';
import { Observable, of } from 'rxjs';
import { ServerErrors } from '../../_models/server-error';

@Component({
  selector: 'ktb-error-view',
  templateUrl: './ktb-error-view.component.html',
  styleUrls: ['../ktb-logout-view/ktb-logout-view.component.scss'],
})
export class KtbErrorViewComponent {
  public logoUrl = environment.config.logoInvertedUrl;
  public error$: Observable<ServerErrors>;
  public ServerErrors = ServerErrors;
  public queryParams$: Observable<ParamMap>;

  @Input() set error(error: ServerErrors) {
    this.error$ = of(error);
  }

  constructor(private route: ActivatedRoute) {
    this.error$ = this.route.queryParamMap.pipe(
      map((params) => params.get('status')),
      map((status) => {
        if (status) {
          const parsedStatus = +status;
          return isNaN(parsedStatus) ? null : parsedStatus;
        }
        return null;
      }),
      map((status) => {
        return status ?? ServerErrors.INTERNAL;
      })
    );

    this.queryParams$ = this.route.queryParamMap;
  }
}
