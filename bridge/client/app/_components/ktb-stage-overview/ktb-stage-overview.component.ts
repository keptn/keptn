import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  Input,
  OnInit,
  Output
} from '@angular/core';
import {Project} from '../../_models/project';
import {Stage} from '../../_models/stage';
import {DataService} from '../../_services/data.service';
import {DtFilterFieldDefaultDataSource} from '@dynatrace/barista-components/filter-field';
import {ApiService} from '../../_services/api.service';
import {Service} from '../../_models/service';
import {Root} from '../../_models/root';

@Component({
  selector: 'ktb-stage-overview',
  templateUrl: './ktb-stage-overview.component.html',
  styleUrls: ['./ktb-stage-overview.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbStageOverviewComponent implements OnInit {
  public _project: Project;
  public selectedStage: Stage = null;
  public _dataSource = new DtFilterFieldDefaultDataSource();
  public filter: any[];
  private filteredServices: string[] = [];
  private globalFilter: {[projectName: string]: {services: string[]}};

  @Output() selectedStageChange: EventEmitter<any> = new EventEmitter();
  @Output() filterChange: EventEmitter<string[]> = new EventEmitter<string[]>();

  @Input()
  get project() {
    return this._project;
  }

  set project(project: Project) {
    if (this._project !== project) {
      this._project = project;
      this.setFilter();
    }
  }

  constructor(private dataService: DataService, private apiService: ApiService, private _changeDetectorRef: ChangeDetectorRef) {
  }

  ngOnInit(): void {
  }

  private setFilter(): void {
    this._dataSource.data = {
      autocomplete: [
        {
          name: 'Services',
          autocomplete: this.project.services.map(service => {
            return {
              name: service.serviceName
            };
          })
        }
      ]
    };
    this.globalFilter = this.apiService.environmentFilter;
    const services = this.globalFilter[this.project.projectName]?.services || [];
    this.filteredServices = services.filter(service => this.project.services.some(pService => pService.serviceName === service));
    this.filterChange.emit(this.filteredServices);
    this.filter = [
      ...this.filteredServices.map(service => {
          return [
            // @ts-ignore
            this._dataSource.data.autocomplete[0],
            {name: service}
          ];
        }
      )
    ];

    this._changeDetectorRef.markForCheck();
  }

  public filterChanged(event: any) {
    this.filteredServices = this.getServicesOfFilter(event);
    this.globalFilter[this.project.projectName] = {services: this.filteredServices};
    this.apiService.environmentFilter = this.globalFilter;
    this.filterChange.emit(this.filteredServices);
    this._changeDetectorRef.markForCheck();
  }

  public filterServices(services: Service[]): Service[] {
    return this.filteredServices.length === 0 ? services : services.filter(service => this.filteredServices.includes(service.serviceName));
  }

  public filterRoots(roots: Root[]): Root[] {
    return this.filteredServices.length === 0 ? roots : roots?.filter(root => this.filteredServices.includes(root.getService()));
  }

  private getServicesOfFilter(event: any): string[] {
    return event.filters.reduce((services, filter) => {
      services.push(filter[1].name);
      return services;
    }, []);
  }

  public trackStage(index: number, stage: string[]): string {
    return stage.toString();
  }

  public selectStage($event, stage: Stage, filterType?: string) {
    this.selectedStage = stage;
    $event.stopPropagation();
    this.selectedStageChange.emit({stage, filterType});
  }

}
