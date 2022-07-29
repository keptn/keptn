import { Component, TemplateRef } from '@angular/core';
import { DtSortEvent, DtTableDataSource } from '@dynatrace/barista-components/table';
import { combineLatestWith, merge, Observable, shareReplay, Subject, switchMap } from 'rxjs';
import { DataService } from '../../../_services/data.service';
import { ActivatedRoute, Router } from '@angular/router';
import { distinctUntilChanged, filter, finalize, map, tap } from 'rxjs/operators';
import { UniformRegistrationLog } from '../../../../../shared/interfaces/uniform-registration-log';
import { UniformRegistration } from '../../../_models/uniform-registration';
import { Location } from '@angular/common';

export type Params = { projectName: string; integrationId?: string };

const sortConfig: Record<string, (u: UniformRegistration) => string> = {
  host: (u) => u.metadata.hostname,
  namespace: (u) => u.metadata.kubernetesmetadata.namespace,
  location: (u) => u.metadata.location,
};

const empty = '--EMPTY--';

@Component({
  selector: 'ktb-keptn-services-list',
  templateUrl: './ktb-integration-view.component.html',
  styleUrls: ['./ktb-integration-view.component.scss'],
})
export class KtbIntegrationViewComponent {
  private selectUniformRegistrationId$ = new Subject<string>();
  public isLoadingUniformRegistrations = true;
  public isLoadingLogs = false;
  public lastSeen?: Date;
  private uniformRegistrations: DtTableDataSource<UniformRegistration> = new DtTableDataSource();

  public params$: Observable<Params> = this.route.paramMap.pipe(
    filter((paramMap) => !!paramMap.get('projectName')),
    map((paramMap) => ({
      projectName: paramMap.get('projectName') ?? empty,
      integrationId: paramMap.get('integrationId') ?? undefined,
    }))
  );

  private registrations$ = this.dataService
    .getUniformRegistrations()
    .pipe(finalize(() => (this.isLoadingUniformRegistrations = false)));

  public uniformRegistrations$ = this.registrations$.pipe(
    map((registrations) => {
      this.uniformRegistrations.data = registrations;
      return this.uniformRegistrations;
    })
  );

  public selectedUniformRegistrationId$ = merge(
    this.params$.pipe(map((params) => params.integrationId ?? empty)),
    this.selectUniformRegistrationId$.asObservable()
  ).pipe(distinctUntilChanged(), shareReplay(1));

  public selectedUniformRegistration$ = this.selectedUniformRegistrationId$.pipe(
    combineLatestWith(this.registrations$),
    map(([regId, registrations]) => (regId ? registrations.find((r) => r.id === regId) : undefined))
  );

  public uniformRegistrationLogs$ = this.selectedUniformRegistrationId$.pipe(
    switchMap((uniformRegistrationId) => this.loadLogs(uniformRegistrationId)),
    map((logs) => sortLogs(logs))
  );

  constructor(
    private dataService: DataService,
    private route: ActivatedRoute,
    private router: Router,
    private location: Location
  ) {}

  public setSelectedUniformRegistration(uniformRegistration: UniformRegistration, projectName: string): void {
    this.isLoadingLogs = true;
    const routeUrl = this.router.createUrlTree([
      '/project',
      projectName,
      'settings',
      'uniform',
      'integrations',
      uniformRegistration.id,
    ]);
    this.location.go(routeUrl.toString());
    this.lastSeen = this.dataService.getUniformDate(uniformRegistration.id);
    uniformRegistration.unreadEventsCount = 0;
    const noUnreadLogs = this.uniformRegistrations.data.every((r) => r.unreadEventsCount === 0);
    if (noUnreadLogs) {
      this.dataService.setHasUnreadUniformRegistrationLogs(false);
    }
    this.selectUniformRegistrationId$.next(uniformRegistration.id);
  }

  public sortData(sortEvent: DtSortEvent): void {
    this.uniformRegistrations.data = sortRegistrations(
      this.uniformRegistrations.data,
      sortEvent.active,
      sortEvent.direction === 'asc'
    );
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

  public getSubscriptions(uniformRegistration: UniformRegistration, projectName: string): string[] {
    return uniformRegistration.subscriptions
      .filter((s) => s.hasProject(projectName, true))
      .map((s) => s.formattedEvent);
  }

  public toUniformRegistration(item: unknown): UniformRegistration {
    return <UniformRegistration>item;
  }

  private loadLogs(uniformRegistrationId: string): Observable<UniformRegistrationLog[]> {
    return this.dataService.getUniformRegistrationLogs(uniformRegistrationId).pipe(
      tap((logs) => this.setUniformDate(uniformRegistrationId, logs)),
      finalize(() => (this.isLoadingLogs = false))
    );
  }

  private setUniformDate(uniformRegistrationId: string, logs: UniformRegistrationLog[]): void {
    this.dataService.setUniformDate(uniformRegistrationId, logs[0]?.time);
  }
}

export function sortRegistrations(
  registrations: UniformRegistration[],
  column: string,
  ascending: boolean
): UniformRegistration[] {
  return [...registrations].sort((a: UniformRegistration, b: UniformRegistration) => {
    const sortBy = sortConfig[column];
    const sortResult = sortBy ? compare(sortBy(a), sortBy(b), ascending) : 0;
    return sortResult || compare(a.name, b.name, ascending);
  });
}

export function sortLogs(logs: UniformRegistrationLog[]): UniformRegistrationLog[] {
  return logs.sort((a, b) => {
    if (a.time.valueOf() > b.time.valueOf()) {
      return -1;
    } else if (a.time.valueOf() < b.time.valueOf()) {
      return 1;
    }
    return 0;
  });
}

export function compare(a: string, b: string, isAsc: boolean): number {
  const result = a.localeCompare(b);
  const factor = isAsc ? 1 : -1;
  return result * factor;
}
