import { Component, Inject, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { DataService } from '../../_services/data.service';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';

@Component({
  selector: 'ktb-service-settings-list',
  templateUrl: './ktb-service-settings-list.component.html',
})
export class KtbServiceSettingsListComponent implements OnDestroy {
  public projectName?: string;
  public isLoading = false;
  public dataSource: DtTableDataSource<string> = new DtTableDataSource<string>();
  private _timer: Subscription = Subscription.EMPTY;

  constructor(
    private router: ActivatedRoute,
    private dataService: DataService,
    @Inject(POLLING_INTERVAL_MILLIS) private initialDelayMillis: number
  ) {
    this.router.paramMap
      .pipe(
        map((params) => params.get('projectName')),
        filter((projectName): projectName is string => !!projectName)
      )
      .subscribe((projectName) => {
        this.projectName = projectName;
        this.isLoading = true;

        this._timer.unsubscribe();

        this._timer = AppUtils.createTimer(0, initialDelayMillis).subscribe(() => {
          if (this.projectName) {
            this.dataService.getServiceNames(this.projectName).subscribe((services) => {
              this.dataSource = new DtTableDataSource<string>(services);
              this.isLoading = false;
            });
          }
        });
      });
  }

  public ngOnDestroy(): void {
    this._timer.unsubscribe();
  }
}
