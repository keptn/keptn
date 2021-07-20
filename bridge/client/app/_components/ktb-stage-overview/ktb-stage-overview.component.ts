import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  Input,
  OnDestroy,
  OnInit,
  Output
} from '@angular/core';
import {Project} from '../../_models/project';
import {Stage} from '../../_models/stage';
import {DataService} from '../../_services/data.service';
import { DtFilterFieldChangeEvent, DtFilterFieldDefaultDataSource } from '@dynatrace/barista-components/filter-field';
import {ApiService} from '../../_services/api.service';
import {Service} from '../../_models/service';
import {Root} from '../../_models/root';
import {filter, takeUntil, tap} from 'rxjs/operators';
import {Subject, Subscription, timer} from 'rxjs';

@Component({
  selector: 'ktb-stage-overview[project]',
  templateUrl: './ktb-stage-overview.component.html',
  styleUrls: ['./ktb-stage-overview.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbStageOverviewComponent implements OnInit, OnDestroy {
  public _project?: Project;
  public selectedStage?: Stage;
  public _dataSource = new DtFilterFieldDefaultDataSource();
  // tslint:disable-next-line:no-any
  public filter: any[] = [];
  private filteredServices: string[] = [];
  private globalFilter: {[projectName: string]: {services: string[]}} = {};
  private unfinishedRoots: Root[] = [];
  private _rootsTimerInterval = 10;
  private _rootsTimer: Subscription = Subscription.EMPTY;
  private readonly unsubscribe$ = new Subject<void>();

  @Output() selectedStageChange: EventEmitter<{ stage: Stage, filterType?: string }> = new EventEmitter();
  @Output() filterChange: EventEmitter<string[]> = new EventEmitter<string[]>();

  @Input()
  get project(): Project | undefined {
    return this._project;
  }

  set project(project: Project | undefined) {
    if (this._project !== project) {
      this._project = project;
      this.setFilter();
    }
  }

  constructor(private dataService: DataService, private apiService: ApiService, private _changeDetectorRef: ChangeDetectorRef) {
  }

  ngOnInit(): void {
    this._rootsTimer = timer(0, this._rootsTimerInterval * 1000)
      .subscribe(() => {
        if (this.project) {
          this.dataService.loadRoots(this.project);
          if (this.unfinishedRoots) {
            this.unfinishedRoots.forEach(root => {
              this.dataService.loadRootTraces(root);
            });
          }
        }
      });

    this.dataService.roots
      .pipe(
        takeUntil(this.unsubscribe$),
        filter((roots: Root[] | undefined): roots is Root[] => !!roots),
        tap(roots => {
            // Set unfinished roots so that the traces for updates can be loaded
            // Also ignore currently selected root, as this is getting already polled
            this.unfinishedRoots = roots.filter(root => root && !root.isFinished());
          }
        )).subscribe();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this._rootsTimer.unsubscribe();
  }

  private setFilter(): void {
    this._dataSource.data =  {
      autocomplete: [
        {
          name: 'Services',
          autocomplete: this.project?.getServices().map(service => {
            return {
              name: service.serviceName
            };
          }) ?? []
        }
      ]
    };
    this.globalFilter = this.apiService.environmentFilter;
    if (this.project) {
      const services = this.globalFilter[this.project.projectName]?.services || [];
      // tslint:disable-next-line:no-non-null-assertion
      this.filteredServices = services.filter(service => this.project!.getServices().some(pService => pService.serviceName === service));
    }
    else {
      this.filteredServices = [];
    }
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

  // tslint:disable-next-line:no-any
  public filterChanged(event: DtFilterFieldChangeEvent<any>) {
    this.filteredServices = this.getServicesOfFilter(event);
    if (this.project) {
      this.globalFilter[this.project.projectName] = {services: this.filteredServices};
    }
    this.apiService.environmentFilter = this.globalFilter;
    this.filterChange.emit(this.filteredServices);
    this._changeDetectorRef.markForCheck();
  }

  public filterServices(services: Service[]): Service[] {
    return this.filteredServices.length === 0 ? services : services.filter(service => this.filteredServices.includes(service.serviceName));
  }

  public filterRoots(roots: Root[]): Root[] {
    return this.filteredServices.length === 0
          ? roots
          : roots?.filter(root => root.service ? this.filteredServices.includes(root.service) : false);
  }

  // tslint:disable-next-line:no-any
  private getServicesOfFilter(event: DtFilterFieldChangeEvent<any>): string[] {
    // tslint:disable-next-line:no-any
    return event.filters.reduce((services: string[], currentFilter: any) => {
      services.push(currentFilter[1].name);
      return services;
    }, []);
  }

  public trackStage(index: number, stage: string[] | null): string | undefined {
    return stage?.toString();
  }

  public selectStage($event: MouseEvent, stage: Stage, filterType?: string) {
    this.selectedStage = stage;
    $event.stopPropagation();
    this.selectedStageChange.emit({stage, filterType});
  }

}
