import { Component, EventEmitter, OnDestroy, OnInit, Output } from '@angular/core';
import { DtSortEvent, DtTableDataSource } from '@dynatrace/barista-components/table';
import { UniformRegistration } from '../../_models/uniform-registration';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute } from '@angular/router';
import { switchMap, takeUntil } from 'rxjs/operators';
import { UniformRegistrationLog } from '../../_models/uniform-registration-log';

@Component({
  selector: 'ktb-keptn-services-list',
  templateUrl: './ktb-keptn-services-list.component.html',
  styleUrls: ['./ktb-keptn-services-list.component.scss']
})
export class KtbKeptnServicesListComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  private selectedUniformRegistrationId$ = new Subject<string>();
  private uniformRegistrationLogsSubject = new BehaviorSubject<UniformRegistrationLog[]>([]);

  public uniformRegistrations: DtTableDataSource<UniformRegistration> = new DtTableDataSource();
  public selectedUniformRegistration?: UniformRegistration;
  public uniformRegistrationLogs$: Observable<UniformRegistrationLog[]> = this.uniformRegistrationLogsSubject.asObservable();
  public isLoadingLogs = false;

  public projectName: string | null = null;

  @Output() selectedUniformRegistrationChanged: EventEmitter<UniformRegistration> = new EventEmitter();

  constructor(private dataService: DataService, private route: ActivatedRoute) {
  }

  ngOnInit(): void {
    this.route.paramMap.pipe(
      takeUntil(this.unsubscribe$)
    ).subscribe(map => {
      this.projectName = map.get('projectName');
    });

    this.selectedUniformRegistrationId$.pipe(
      takeUntil(this.unsubscribe$),
      switchMap(uniformRegistrationId => {
        this.isLoadingLogs = true;
        return this.dataService.getUniformRegistrationLogs(uniformRegistrationId);
      })
    ).subscribe((uniformRegLogs) => {
      uniformRegLogs.sort(this.sortLogs);
      this.isLoadingLogs = false;
      this.uniformRegistrationLogsSubject.next(uniformRegLogs);
    });

    this.dataService.getUniformRegistrations()
      .subscribe((uniformRegistrations) => {
        this.uniformRegistrations.data = uniformRegistrations;
      });
  }

  private sortLogs(a: UniformRegistrationLog, b: UniformRegistrationLog): number {
    let status = 0;
    if (a.time.valueOf() > b.time.valueOf()) {
      status = -1;
    }
    else if (a.time.valueOf() < b.time.valueOf()) {
      status = 1;
    }
    return status;
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

  public setSelectedUniformRegistration(uniformRegistration: UniformRegistration) {
    if (this.selectedUniformRegistration !== uniformRegistration) {
      this.selectedUniformRegistration = uniformRegistration;
      this.selectedUniformRegistrationId$.next(this.selectedUniformRegistration.id);
      this.selectedUniformRegistrationChanged.emit(uniformRegistration);
    }
  }

  public formatSubscriptions(subscriptions: string[]): string {
    return subscriptions.join('<br/>');
  }

  public sortData(sortEvent: DtSortEvent) {
    if (this.uniformRegistrations.data) {
      const isAscending = sortEvent.direction === 'asc';
      const data: UniformRegistration[] = this.uniformRegistrations.data.slice();

      data.sort((a: UniformRegistration, b: UniformRegistration) => {
        switch (sortEvent.active) {
          case 'host':
            return (this.compare(a.metadata.hostname, b.metadata.hostname, isAscending) || this.compare(a.name, b.name, true));
          case 'namespace':
            return this.compare(a.metadata.kubernetesmetadata.namespace, b.metadata.kubernetesmetadata.namespace, isAscending) || this.compare(a.name, b.name, true);
          case 'location':
            return this.compare(a.metadata.location, b.metadata.location, isAscending) || this.compare(a.name, b.name, true);
          case 'name':
          default:
            return this.compare(a.name, b.name, isAscending);
        }
      });

      this.uniformRegistrations.data = data;
    } else {
      this.uniformRegistrations.data = [];
    }
  }

  public toRegistration(row: UniformRegistration): UniformRegistration {
    return row;
  }

  private compare(a: string, b: string, isAsc: boolean): number {
    const result = a.localeCompare(b);
    if (result !== 0 && !isAsc) {
      return -result;
    }
    return result;
  }
}
