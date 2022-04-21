import { Component, OnDestroy } from '@angular/core';
import { environment } from '../../../environments/environment';
import { ActivatedRoute } from '@angular/router';
import { map, takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';

enum ServerErrors {
  INTERNAL = 500,
  INSUFFICIENT_PERMISSION = 403,
}

@Component({
  selector: 'ktb-error-view',
  templateUrl: './ktb-error-view.component.html',
  styleUrls: ['../ktb-logout-view/ktb-logout-view.component.scss'],
})
export class KtbErrorViewComponent implements OnDestroy {
  public logoUrl = environment.config.logoInvertedUrl;
  public error?: ServerErrors;
  public ServerErrors = ServerErrors;
  private unsubscribe$ = new Subject<void>();

  constructor(private route: ActivatedRoute) {
    this.route.queryParamMap
      .pipe(
        map((params) => params.get('status')),
        map((status) => {
          if (status) {
            const parsedStatus = +status;
            return isNaN(parsedStatus) ? null : parsedStatus;
          }
          return null;
        }),
        takeUntil(this.unsubscribe$)
      )
      .subscribe((status) => {
        if (status === null) {
          this.error = ServerErrors.INTERNAL;
        } else {
          this.error = status;
        }
      });
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
