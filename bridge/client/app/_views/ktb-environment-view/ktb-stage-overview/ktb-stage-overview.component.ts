import { AfterContentInit, Component, EventEmitter, OnDestroy, Output } from '@angular/core';
import { Project } from '../../../_models/project';
import { Stage } from '../../../_models/stage';
import { DataService } from '../../../_services/data.service';
import { DtFilterFieldChangeEvent, DtFilterFieldDefaultDataSource } from '@dynatrace/barista-components/filter-field';
import { ApiService } from '../../../_services/api.service';
import { Service } from '../../../_models/service';
import { DtAutoComplete, DtFilter, DtFilterArray } from '../../../_models/dt-filter';
import { distinctUntilChanged, filter, map, switchMap, takeUntil, tap, withLatestFrom } from 'rxjs/operators';
import { ActivatedRoute, Router } from '@angular/router';
import { combineLatest, Subject } from 'rxjs';
import { DtFilterFieldDefaultDataSourceAutocomplete } from '@dynatrace/barista-components/filter-field/src/filter-field-default-data-source';
import { ServiceFilterType } from '../ktb-stage-details/ktb-stage-details.component';

@Component({
  selector: 'ktb-stage-overview',
  templateUrl: './ktb-stage-overview.component.html',
  styleUrls: ['./ktb-stage-overview.component.scss'],
})
export class KtbStageOverviewComponent implements AfterContentInit, OnDestroy {
  public _dataSource = new DtFilterFieldDefaultDataSource();
  public filter: DtFilterArray[] = [];
  public isTriggerSequenceOpen: boolean;
  private filteredServices: string[] = [];
  private globalFilter: { [projectName: string]: { services: string[] } } = {};
  private unsubscribe$: Subject<void> = new Subject<void>();

  public project$ = this.route.params.pipe(
    map((params) => params.projectName),
    filter((projectName): projectName is string => !!projectName),
    distinctUntilChanged(),
    switchMap((projectName) => this.dataService.getProject(projectName)),
    filter((project): project is Project => !!project),
    tap((project) => {
      this.setFilter(project, true);
    })
  );

  public readonly selectedStageName$ = this.route.paramMap.pipe(
    map((params) => params.get('stageName')),
    withLatestFrom(this.project$),
    filter(([stageName, project]) => Boolean(stageName && project)),
    map(([stageName]) => stageName)
  );

  private readonly paramFilterType$ = this.route.queryParamMap.pipe(map((params) => params.get('filterType')));

  @Output() selectedStageChange: EventEmitter<{ stage: Stage; filterType: ServiceFilterType }> = new EventEmitter();
  @Output() filteredServicesChange: EventEmitter<string[]> = new EventEmitter<string[]>();

  constructor(
    private dataService: DataService,
    private apiService: ApiService,
    private route: ActivatedRoute,
    private router: Router
  ) {
    this.isTriggerSequenceOpen = this.dataService.isTriggerSequenceOpen;
    this.dataService.isTriggerSequenceOpen = false;
  }

  ngAfterContentInit(): void {
    combineLatest([this.selectedStageName$, this.paramFilterType$, this.project$])
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(([stageName, filterType, project]) => {
        const stage = project.stages.find((s) => s.stageName === stageName);
        if (stage) {
          this.selectedStageChange.emit({ stage: stage, filterType: (filterType as ServiceFilterType) ?? undefined });
        }
      });
  }

  private setFilter(project: Project | undefined, projectChanged: boolean): void {
    this._dataSource.data = {
      autocomplete: [
        {
          name: 'Services',
          autocomplete:
            project?.getServices().map((service) => ({
              name: service.serviceName,
            })) ?? [],
        } as DtAutoComplete,
      ],
    };
    this.globalFilter = this.apiService.environmentFilter;

    const createFilter = (p: Project): string[] => {
      // services can be deleted or added; adjust filter
      const services = this.globalFilter[p.projectName]?.services || [];
      return services.filter((service) => p.getServices().some((pService) => pService.serviceName === service));
    };
    const newFilter: string[] = project ? createFilter(project) : [];

    if (projectChanged || newFilter.length !== this.filteredServices.length) {
      this.filteredServices = newFilter;
      this.filteredServicesChange.emit(this.filteredServices);
      this.filter = [
        ...this.filteredServices.map(
          (service) =>
            [
              (this._dataSource.data as DtFilterFieldDefaultDataSourceAutocomplete).autocomplete[0],
              { name: service },
            ] as DtFilterArray
        ),
      ];
    }
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public filterChanged(project: Project, event: DtFilterFieldChangeEvent<any>): void {
    // can't set another type because of "is not assignable to..."
    this.filteredServices = this.getServicesOfFilter(event);
    this.globalFilter[project.projectName] = { services: this.filteredServices };

    this.apiService.environmentFilter = this.globalFilter;
    this.filteredServicesChange.emit(this.filteredServices);
  }

  public filterServices(services: Service[]): Service[] {
    return this.filteredServices.length === 0
      ? services
      : services.filter((service) => this.filteredServices.includes(service.serviceName));
  }

  private getServicesOfFilter(event: DtFilterFieldChangeEvent<DtFilter>): string[] {
    const services: string[] = [];
    for (const currentFilter of event.filters) {
      services.push((currentFilter as DtFilterArray)[1].name);
    }
    return services;
  }

  public changeIsTriggerSequence(state: boolean): void {
    this.isTriggerSequenceOpen = state;
  }

  public trackStage(_index: number, stage: string[] | null): string | undefined {
    return stage?.toString();
  }

  public selectStage($event: MouseEvent, project: Project, stage: Stage, filterType?: ServiceFilterType): void {
    this.router.navigate(['/project', project.projectName, 'environment', 'stage', stage.stageName], {
      queryParams: { filterType: filterType },
    });
    $event.stopPropagation();
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
