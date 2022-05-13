import { Component, EventEmitter, OnDestroy, OnInit, AfterContentInit, Output } from '@angular/core';
import { Project } from '../../_models/project';
import { Stage } from '../../_models/stage';
import { DataService } from '../../_services/data.service';
import { DtFilterFieldChangeEvent, DtFilterFieldDefaultDataSource } from '@dynatrace/barista-components/filter-field';
import { ApiService } from '../../_services/api.service';
import { Service } from '../../_models/service';
import { DtAutoComplete, DtFilter, DtFilterArray } from '../../_models/dt-filter';
import { filter, map, switchMap, takeUntil, tap } from 'rxjs/operators';
import { Router, ActivatedRoute } from '@angular/router';
import { Subject } from 'rxjs';
import { DtFilterFieldDefaultDataSourceAutocomplete } from '@dynatrace/barista-components/filter-field/src/filter-field-default-data-source';
import { ServiceFilterType } from '../ktb-stage-details/ktb-stage-details.component';

@Component({
  selector: 'ktb-stage-overview',
  templateUrl: './ktb-stage-overview.component.html',
  styleUrls: ['./ktb-stage-overview.component.scss'],
})
export class KtbStageOverviewComponent implements OnDestroy, OnInit, AfterContentInit {
  public project?: Project;
  public selectedStage?: Stage;
  public _dataSource = new DtFilterFieldDefaultDataSource();
  public filter: DtFilterArray[] = [];
  public isTriggerSequenceOpen = false;
  private filteredServices: string[] = [];
  private globalFilter: { [projectName: string]: { services: string[] } } = {};
  private unsubscribe$: Subject<void> = new Subject<void>();

  @Output() selectedStageChange: EventEmitter<{ stage: Stage; filterType: ServiceFilterType }> = new EventEmitter();
  @Output() filteredServicesChange: EventEmitter<string[]> = new EventEmitter<string[]>();

  constructor(
    private dataService: DataService,
    private apiService: ApiService,
    private router: Router,
    private route: ActivatedRoute
  ) {}

  public ngOnInit(): void {
    // needs to be in init because of emitter
    const project$ = this.route.params.pipe(
      map((params) => params.projectName),
      filter((projectName) => !!projectName),
      tap(() => {
        this.isTriggerSequenceOpen = this.dataService.isTriggerSequenceOpen;
        this.dataService.isTriggerSequenceOpen = false;
      }),
      switchMap((projectName) => this.dataService.getProject(projectName)),
      takeUntil(this.unsubscribe$)
    );

    project$.subscribe((project) => {
      const differentProject = project?.projectName !== this.project?.projectName;
      this.project = project;
      this.setFilter(differentProject);
    });
  }

  public ngAfterContentInit(): void {
    this.route.queryParams.subscribe((params) => {
      if (params['stage'] && this.project) {
        let stage = this.project.getStage(params['stage']);
        if (stage) {
          this.selectedStage = stage;
          this.selectedStageChange.emit({ stage: stage, filterType: undefined });
        }
      }
    });
  }

  private setFilter(projectChanged: boolean): void {
    this._dataSource.data = {
      autocomplete: [
        {
          name: 'Services',
          autocomplete:
            this.project?.getServices().map((service) => ({
              name: service.serviceName,
            })) ?? [],
        } as DtAutoComplete,
      ],
    };
    this.globalFilter = this.apiService.environmentFilter;
    let newFilter: string[];
    if (this.project) {
      // services can be deleted or added; adjust filter
      const services = this.globalFilter[this.project.projectName]?.services || [];
      newFilter = services.filter((service) =>
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        this.project!.getServices().some((pService) => pService.serviceName === service)
      );
    } else {
      newFilter = [];
    }
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
  public filterChanged(event: DtFilterFieldChangeEvent<any>): void {
    // can't set another type because of "is not assignable to..."
    this.filteredServices = this.getServicesOfFilter(event);
    if (this.project) {
      this.globalFilter[this.project.projectName] = { services: this.filteredServices };
    }
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

  public trackStage(index: number, stage: string[] | null): string | undefined {
    return stage?.toString();
  }

  public selectStage($event: MouseEvent, stage: Stage): void {
    this.router.navigate([], {
      queryParams: {
        stage: stage.stageName,
      },
      queryParamsHandling: 'merge',
    });
    $event.stopPropagation();
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
