import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subject } from 'rxjs';
import { filter, map, switchMap, takeUntil } from 'rxjs/operators';
import { DataService } from '../../_services/data.service';
import { DtTableDataSource } from '@dynatrace/barista-components/table';

@Component({
  selector: 'ktb-service-settings-list',
  templateUrl: './ktb-service-settings-list.component.html',
})
export class KtbServiceSettingsListComponent implements OnDestroy {
  public projectName?: string;
  public dataSource: DtTableDataSource<string> = new DtTableDataSource<string>();
  private unsubscribe$: Subject<void> = new Subject<void>();

  constructor(private router: ActivatedRoute, private dataService: DataService) {
    const projectName$ = this.router.paramMap.pipe(
      map(params => params.get('projectName')),
      filter((projectName): projectName is string => !!projectName),
    );

    projectName$.pipe(
      switchMap(projectName => {
        this.projectName = projectName;
        return this.dataService.getProject(projectName);
      }),
      takeUntil(this.unsubscribe$),
    ).subscribe(project => {
      const services: string[] = project?.getServices()?.map(service => service.serviceName) ?? [];
      this.dataSource = new DtTableDataSource<string>(services);
    });
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
