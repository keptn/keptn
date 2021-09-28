import { ChangeDetectionStrategy, ChangeDetectorRef, Component, EventEmitter, OnDestroy, Output } from '@angular/core';
import { Project } from '../../_models/project';
import { Stage } from '../../_models/stage';
import { DataService } from '../../_services/data.service';
import { DtFilterFieldChangeEvent, DtFilterFieldDefaultDataSource } from '@dynatrace/barista-components/filter-field';
import { ApiService } from '../../_services/api.service';
import { Service } from '../../_models/service';
import { DtAutoComplete, DtFilter, DtFilterArray } from '../../_models/dt-filter';
import { filter, map, switchMap, takeUntil } from 'rxjs/operators';
import { ActivatedRoute } from '@angular/router';
import { Subject } from 'rxjs';

@Component({
  selector: 'ktb-stage-overview',
  templateUrl: './ktb-stage-overview.component.html',
  styleUrls: ['./ktb-stage-overview.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbStageOverviewComponent implements OnDestroy {
  public project?: Project;
  public selectedStage?: Stage;
  public _dataSource = new DtFilterFieldDefaultDataSource();
  public filter: DtFilterArray[] = [];
  private filteredServices: string[] = [];
  private globalFilter: { [projectName: string]: { services: string[] } } = {};
  private unsubscribe$: Subject<void> = new Subject<void>();

  @Output() selectedStageChange: EventEmitter<{ stage: Stage; filterType?: string }> = new EventEmitter();
  @Output() filterChange: EventEmitter<string[]> = new EventEmitter<string[]>();

  constructor(
    private dataService: DataService,
    private apiService: ApiService,
    private _changeDetectorRef: ChangeDetectorRef,
    private route: ActivatedRoute
  ) {
    const project$ = this.route.params.pipe(
      map((params) => params.projectName),
      filter((projectName) => !!projectName),
      switchMap((projectName) => this.dataService.getProject(projectName)),
      takeUntil(this.unsubscribe$)
    );
    project$.subscribe((project) => {
      this.project = project;
      this.setFilter();
    });
  }

  private setFilter(): void {
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
    if (this.project) {
      const services = this.globalFilter[this.project.projectName]?.services || [];
      this.filteredServices = services.filter((service) =>
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        this.project!.getServices().some((pService) => pService.serviceName === service)
      );
    } else {
      this.filteredServices = [];
    }
    this.filterChange.emit(this.filteredServices);
    this.filter = [
      ...this.filteredServices.map((service) => [
          // @ts-ignore
          this._dataSource.data.autocomplete[0],
          { name: service },
        ] as DtFilterArray),
    ];

    this._changeDetectorRef.markForCheck();
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public filterChanged(event: DtFilterFieldChangeEvent<any>): void {
    // can't set another type because of "is not assignable to..."
    this.filteredServices = this.getServicesOfFilter(event);
    if (this.project) {
      this.globalFilter[this.project.projectName] = { services: this.filteredServices };
    }
    this.apiService.environmentFilter = this.globalFilter;
    this.filterChange.emit(this.filteredServices);
    this._changeDetectorRef.markForCheck();
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

  public selectStage($event: MouseEvent, stage: Stage, filterType?: string): void {
    this.selectedStage = stage;
    $event.stopPropagation();
    this.selectedStageChange.emit({ stage, filterType });
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
