import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { Project } from '../../../_models/project';
import { Stage } from '../../../_models/stage';
import { DataService } from '../../../_services/data.service';
import { DtFilterFieldChangeEvent, DtFilterFieldDefaultDataSource } from '@dynatrace/barista-components/filter-field';
import { ApiService } from '../../../_services/api.service';
import { Service } from '../../../_models/service';
import { DtAutoComplete, DtFilter, DtFilterArray } from '../../../_models/dt-filter';
import { DtFilterFieldDefaultDataSourceAutocomplete } from '@dynatrace/barista-components/filter-field/src/filter-field-default-data-source';
import { ServiceFilterType } from '../ktb-stage-details/ktb-stage-details.component';
import { ISelectedStageInfo } from '../ktb-environment-view.component';

@Component({
  selector: 'ktb-stage-overview',
  templateUrl: './ktb-stage-overview.component.html',
  styleUrls: ['./ktb-stage-overview.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbStageOverviewComponent {
  public _dataSource = new DtFilterFieldDefaultDataSource();
  public filter: DtFilterArray[] = [];
  public isTriggerSequenceOpen: boolean;
  private filteredServices: string[] = [];
  private globalFilter: { [projectName: string]: { services: string[] } } = {};
  private _project?: Project;

  @Input() selectedStageInfo?: ISelectedStageInfo;

  @Input() set project(project: Project | undefined) {
    this._project = project;
    if (project) {
      this.setFilter(project, true);
    }
  }
  get project(): Project | undefined {
    return this._project;
  }

  @Output() selectedStageInfoChange: EventEmitter<ISelectedStageInfo> = new EventEmitter();
  @Output() filteredServicesChange: EventEmitter<string[]> = new EventEmitter<string[]>();

  constructor(private dataService: DataService, private apiService: ApiService) {
    this.isTriggerSequenceOpen = this.dataService.isTriggerSequenceOpen;
    this.dataService.isTriggerSequenceOpen = false;
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

  public trackStage(_index: number, stage: string[] | null): string | undefined {
    return stage?.toString();
  }

  public selectStage($event: MouseEvent, project: Project, stage: Stage, filterType?: ServiceFilterType): void {
    $event.stopPropagation();
    this.selectedStageInfoChange.emit({ stage, filterType });
  }
}
