import { Component, EventEmitter, OnDestroy, OnInit, Output, TemplateRef } from '@angular/core';
import { DtSortEvent, DtTableDataSource } from '@dynatrace/barista-components/table';
import { BehaviorSubject, combineLatest, Observable, Subject } from 'rxjs';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, Router } from '@angular/router';
import { filter, map, switchMap, takeUntil } from 'rxjs/operators';
import { UniformRegistrationLog } from '../../../../server/interfaces/uniform-registration-log';
import { UniformRegistration } from '../../_models/uniform-registration';
import { Location } from '@angular/common';
import { UniformSubscription } from '../../_models/uniform-subscription';

@Component({
  selector: 'ktb-keptn-services-list',
  templateUrl: './ktb-keptn-services-list.component.html',
  styleUrls: ['./ktb-keptn-services-list.component.scss'],
})
export class KtbKeptnServicesListComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  private selectedUniformRegistrationId$ = new Subject<string>();
  private uniformRegistrationLogsSubject = new BehaviorSubject<UniformRegistrationLog[]>([]);
  public UniformRegistrationClass = UniformRegistration;
  public isLoadingUniformRegistrations = true;
  public uniformRegistrations: DtTableDataSource<UniformRegistration> = new DtTableDataSource();
  public selectedUniformRegistration?: UniformRegistration;
  public uniformRegistrationLogs$: Observable<UniformRegistrationLog[]> =
    this.uniformRegistrationLogsSubject.asObservable();
  public isLoadingLogs = false;
  public projectName?: string;
  public lastSeen?: Date;

  @Output() selectedUniformRegistrationChanged: EventEmitter<UniformRegistration> = new EventEmitter();

  constructor(
    private dataService: DataService,
    private route: ActivatedRoute,
    private router: Router,
    private location: Location
  ) {
    this.route.paramMap
      .pipe(
        map((paramMap) => paramMap.get('projectName')),
        takeUntil(this.unsubscribe$),
        filter((projectName: string | null): projectName is string => !!projectName)
      )
      .subscribe((projectName) => {
        this.projectName = projectName;
      });
  }

  ngOnInit(): void {
    this.selectedUniformRegistrationId$
      .pipe(
        takeUntil(this.unsubscribe$),
        switchMap((uniformRegistrationId) => {
          this.isLoadingLogs = true;
          const routeUrl = this.router.createUrlTree([
            '/project',
            this.projectName,
            'uniform',
            'services',
            uniformRegistrationId,
          ]);
          this.location.go(routeUrl.toString());
          return this.dataService.getUniformRegistrationLogs(uniformRegistrationId);
        })
      )
      .subscribe((uniformRegLogs) => {
        uniformRegLogs.sort(this.sortLogs);
        this.isLoadingLogs = false;
        if (this.selectedUniformRegistration) {
          this.dataService.setUniformDate(this.selectedUniformRegistration.id, uniformRegLogs[0]?.time);
        }
        this.uniformRegistrationLogsSubject.next(uniformRegLogs);
      });

    const registrations$ = this.dataService.getUniformRegistrations();
    const integrationId$ = this.route.paramMap.pipe(map((paramMap) => paramMap.get('integrationId')));

    combineLatest([registrations$, integrationId$])
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(([uniformRegistrations, integrationId]) => {
        this.isLoadingUniformRegistrations = false;
        this.uniformRegistrations.data = uniformRegistrations;
        if (integrationId) {
          const selectedUniformRegistration = uniformRegistrations.find(
            (uniformRegistration) => uniformRegistration.id === integrationId
          );
          if (selectedUniformRegistration) {
            this.setSelectedUniformRegistration(selectedUniformRegistration);
          }
        }
      });
  }

  private sortLogs(a: UniformRegistrationLog, b: UniformRegistrationLog): number {
    let status = 0;
    if (a.time.valueOf() > b.time.valueOf()) {
      status = -1;
    } else if (a.time.valueOf() < b.time.valueOf()) {
      status = 1;
    }
    return status;
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

  public setSelectedUniformRegistration(uniformRegistration: UniformRegistration) {
    if (this.selectedUniformRegistration !== uniformRegistration) {
      this.lastSeen = this.dataService.getUniformDate(uniformRegistration.id);
      if (this.selectedUniformRegistration) {
        this.selectedUniformRegistration.unreadEventsCount = 0;
        if (!this.uniformRegistrations.data.some((registration) => registration.unreadEventsCount !== 0)) {
          this.dataService.setHasUnreadUniformRegistrationLogs(false);
        }
      }
      this.selectedUniformRegistration = uniformRegistration;
      this.selectedUniformRegistrationId$.next(this.selectedUniformRegistration.id);
      this.selectedUniformRegistrationChanged.emit(uniformRegistration);
    }
  }

  public sortData(sortEvent: DtSortEvent) {
    if (this.uniformRegistrations.data) {
      const isAscending = sortEvent.direction === 'asc';
      const data: UniformRegistration[] = this.uniformRegistrations.data.slice();

      data.sort((a: UniformRegistration, b: UniformRegistration) => {
        switch (sortEvent.active) {
          case 'host':
            return (
              this.compare(a.metadata.hostname, b.metadata.hostname, isAscending) || this.compare(a.name, b.name, true)
            );
          case 'namespace':
            return (
              this.compare(
                a.metadata.kubernetesmetadata.namespace,
                b.metadata.kubernetesmetadata.namespace,
                isAscending
              ) || this.compare(a.name, b.name, true)
            );
          case 'location':
            return (
              this.compare(a.metadata.location, b.metadata.location, isAscending) || this.compare(a.name, b.name, true)
            );
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

  public getOverlay(
    registration: UniformRegistration,
    projectName: string,
    template: TemplateRef<unknown>
  ): TemplateRef<unknown> {
    // The overlay must be conditional but in general directives are not meant to be dynamic.
    // That's why we ignore the fact that undefined is not assignable to TemplateRef
    return (registration.hasSubscriptions(projectName) ? undefined : template) as TemplateRef<unknown>;
  }

  private compare(a: string, b: string, isAsc: boolean): number {
    const result = a.localeCompare(b);
    if (result !== 0 && !isAsc) {
      return -result;
    }
    return result;
  }

  public formatSubscriptions(uniformRegistration: UniformRegistration, projectName: string): string | undefined {
    const subscriptions = uniformRegistration.subscriptions.reduce(
      (accSubscriptions: string[], subscription: UniformSubscription) => {
        if (subscription.hasProject(projectName, true)) {
          accSubscriptions.push(subscription.formattedEvent);
        }
        return accSubscriptions;
      },
      []
    );
    return subscriptions.length !== 0 ? subscriptions.join('<br/>') : undefined;
  }
}
