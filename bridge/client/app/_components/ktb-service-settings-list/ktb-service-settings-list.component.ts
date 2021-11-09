import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subject } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { DataService } from '../../_services/data.service';
import { DtTableDataSource } from '@dynatrace/barista-components/table';

@Component({
  selector: 'ktb-service-settings-list',
  templateUrl: './ktb-service-settings-list.component.html',
})
export class KtbServiceSettingsListComponent implements OnDestroy {
  public projectName?: string;
  public isLoading = false;
  public dataSource: DtTableDataSource<string> = new DtTableDataSource<string>();
  private unsubscribe$: Subject<void> = new Subject<void>();

  constructor(private router: ActivatedRoute, private dataService: DataService) {
    this.router.paramMap
      .pipe(
        map((params) => params.get('projectName')),
        filter((projectName): projectName is string => !!projectName)
      )
      .subscribe((projectName) => {
        this.projectName = projectName;
        this.isLoading = true;
        this.dataService.getServiceNames(this.projectName).subscribe((services) => {
          this.dataSource = new DtTableDataSource<string>(services);
          this.isLoading = false;
        });
      });
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
