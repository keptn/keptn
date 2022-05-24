import { ChangeDetectorRef, Component, HostBinding, Inject, OnDestroy, OnInit, ViewEncapsulation } from '@angular/core';
import { Location } from '@angular/common';
import { ActivatedRoute, Params, Router } from '@angular/router';
import {
  DtQuickFilterChangeEvent,
  DtQuickFilterDefaultDataSource,
  DtQuickFilterDefaultDataSourceConfig,
} from '@dynatrace/barista-components/quick-filter';
import { isObject } from '@dynatrace/barista-components/core';
import { combineLatest, Observable, Subject, Subscription } from 'rxjs';
import { filter, map, switchMap, takeUntil, takeWhile } from 'rxjs/operators';
import moment from 'moment';
import { Project } from '../../_models/project';
import { DataService } from '../../_services/data.service';
import { DateUtil } from '../../_utils/date.utils';
import { Sequence } from '../../_models/sequence';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { ISequencesMetadata } from '../../../../shared/interfaces/sequencesMetadata';
import { ApiService } from '../../_services/api.service';

enum FilterName {
  SERVICE = 'Service',
  STAGE = 'Stage',
  SEQUENCE = 'Sequence',
  STATUS = 'Status',
}

export type FilterType = [
  {
    name: FilterName;
    autocomplete: { name: string; value: string }[];
    showInSidebar: boolean;
  },
  ...{ name: string; value: string }[]
];

const SEQUENCE_STATUS = {
  started: 'Active',
  waiting: 'Waiting',
  failed: 'Failed',
  aborted: 'Aborted',
  succeeded: 'Succeeded',
} as Record<string, string>;

@Component({
  selector: 'ktb-sequence-view',
  templateUrl: './ktb-sequence-view.component.html',
  styleUrls: ['./ktb-sequence-view.component.scss'],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
})
export class KtbSequenceViewComponent implements OnInit, OnDestroy {
  @HostBinding('class') cls = 'ktb-sequence-view';
  private readonly unsubscribe$ = new Subject<void>();
  /** configuration for the quick filter */
  private filterFieldData = {
    autocomplete: [
      {
        name: 'Service',
        showInSidebar: true,
        autocomplete: [],
      },
      {
        name: 'Stage',
        showInSidebar: true,
        autocomplete: [],
      },
      {
        name: 'Sequence',
        showInSidebar: true,
        autocomplete: [],
      },
      {
        name: 'Status',
        showInSidebar: true,
        autocomplete: Object.entries(SEQUENCE_STATUS).map(([value, name]) => ({
          name: name.toString(),
          value: value.toString().toLowerCase(),
        })),
      },
    ],
  };
  private _config: DtQuickFilterDefaultDataSourceConfig = {
    // Method to decide if a node should be displayed in the quick filter
    showInSidebar: (node) => isObject(node) && node.showInSidebar,
  };
  private project?: Project;
  private unfinishedSequences: Sequence[] = [];
  private _tracesTimerInterval = 10_000;
  private _sequenceTimerInterval = 30_000;
  private _tracesTimer: Subscription = Subscription.EMPTY;
  private sequencesUpdated$: Observable<void>;

  public project$: Observable<Project | undefined>;
  private projectName$: Observable<string>;
  public currentSequence?: Sequence;
  public currentLatestDeployedImage?: string;
  public selectedStage?: string;
  public _filterDataSource = new DtQuickFilterDefaultDataSource(this.filterFieldData, this._config);
  public _seqFilters: FilterType[] = [];
  public metadata: ISequencesMetadata = {
    deployments: [],
    filter: {
      stages: [],
      services: [],
    },
  };

  public filteredSequences: Sequence[] = [];
  public loading = false;

  constructor(
    private dataService: DataService,
    private apiService: ApiService,
    private route: ActivatedRoute,
    public dateUtil: DateUtil,
    private router: Router,
    private location: Location,
    private changeDetectorRef_: ChangeDetectorRef,
    @Inject(POLLING_INTERVAL_MILLIS) private initialDelayMillis: number
  ) {
    if (this.initialDelayMillis === 0) {
      this._sequenceTimerInterval = 0;
      this._tracesTimerInterval = 0;
    }

    this.projectName$ = this.route.params.pipe(map((params) => params.projectName));

    this.sequencesUpdated$ = this.dataService.sequencesUpdated;

    this.project$ = this.projectName$.pipe(switchMap((projectName) => this.dataService.getProject(projectName)));

    this.project$
      .pipe(
        takeUntil(this.unsubscribe$),
        filter(
          (project: Project | undefined): project is Project => !!project && !!project.getServices() && !!project.stages
        )
      )
      .subscribe((project: Project) => {
        const currentProjectName = this.project?.projectName;
        if (project.projectName !== currentProjectName) {
          this.currentSequence = undefined;
          this.selectedStage = undefined;
        }
        this.project = project;
      });
  }

  ngOnInit(): void {
    let initParametersHandled = false;
    AppUtils.createTimer(0, this._sequenceTimerInterval)
      .pipe(
        switchMap(() => this.project$),
        filter((project: Project | undefined): project is Project => !!project && !!project.getServices()),
        takeUntil(this.unsubscribe$)
      )
      .subscribe((project) => {
        this.dataService.loadSequences(project);
        this.updateFilterSequence(project.sequences);
        this.changeDetectorRef_.detectChanges();
      });

    AppUtils.createTimer(0, this._tracesTimerInterval)
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(() => {
        // This triggers the subscription for sequences$
        if (this.project) {
          for (const sequence of this.unfinishedSequences) {
            this.dataService.updateSequence(this.project.projectName, sequence.shkeptncontext);
          }
        }
      });

    AppUtils.createTimer(0, this._sequenceTimerInterval)
      .pipe(
        switchMap(() => this.projectName$),
        takeUntil(this.unsubscribe$)
      )
      .subscribe((projectName) => {
        this.loadSequenceMetadata(projectName);
      });

    // init; set parameters
    combineLatest([this.route.params, this.sequencesUpdated$])
      .pipe(
        takeUntil(this.unsubscribe$),
        takeWhile(([params]) => !this.currentSequence && params.shkeptncontext && !initParametersHandled)
        // initParametersHandled to prevent executing it again after traces are loaded if a eventId is provided
      )
      .subscribe(([params]: [Params, void]) => {
        if (params.shkeptncontext && this.project?.sequences) {
          const sequence = this.project.sequences.find((s) => s.shkeptncontext === params.shkeptncontext);
          const stage = params.eventId ? undefined : params.stage;
          const eventId = params.eventId;

          if (sequence) {
            if (params.eventId && !sequence.traces.length) {
              // at this moment, no traces for the sequence are loaded, wait till next sequencesUpdated$
              initParametersHandled = true;
              this.loadTraces(sequence, params.eventId);
            } else {
              this.selectSequence({ sequence, stage, eventId });
            }
          } else if (params.shkeptncontext && this.project) {
            // is running twice because project is changed on start before the first call finishes
            this.dataService.loadUntilRoot(this.project, params.shkeptncontext);
          }
        }
      });

    this.sequencesUpdated$.pipe(takeUntil(this.unsubscribe$)).subscribe(() => {
      this.loading = false;
      this.updateSequenceView();
      this.updateSequencesData();
    });

    this.route.queryParams.pipe(takeUntil(this.unsubscribe$)).subscribe((queryParams) => {
      if (Object.keys(queryParams).length === 0) {
        this.loadSequenceFilters();
      } else {
        const sequenceFilters = Object.keys(queryParams).reduce((params: { [p: string]: string[] }, param) => {
          params[param] = Array.isArray(queryParams[param]) ? queryParams[param] : [queryParams[param]];
          return params;
        }, {});
        this.setSequenceFilters(sequenceFilters);
      }
    });
  }

  public loadSequenceMetadata(projectName: string): void {
    this.dataService.getSequenceMetadata(projectName).subscribe((metadata) => {
      this.metadata = metadata;
      this.updateLatestDeployedImage();
      this.updateFilterDataSource(metadata);
    });
  }

  public selectSequence(event: { sequence: Sequence; stage?: string; eventId?: string }, loadTraces = true): void {
    if (event.eventId) {
      const routeUrl = this.router.createUrlTree(
        ['/project', event.sequence.project, 'sequence', event.sequence.shkeptncontext, 'event', event.eventId],
        { queryParamsHandling: 'preserve' }
      );
      this.location.go(routeUrl.toString());
    } else {
      const stage = event.stage || event.sequence.getStages().pop();
      const routeUrl = this.router.createUrlTree(
        [
          '/project',
          event.sequence.project,
          'sequence',
          event.sequence.shkeptncontext,
          ...(stage ? ['stage', stage] : []),
        ],
        { queryParamsHandling: 'preserve' }
      );
      this.location.go(routeUrl.toString());
    }

    this.currentSequence = event.sequence;
    this.selectedStage = event.stage || event.sequence.getStages().pop();
    this.updateLatestDeployedImage();
    if (loadTraces) {
      this.loadTraces(this.currentSequence);
    }
  }

  public updateSequenceView(): void {
    const sequences = this.project?.sequences;
    if (sequences !== undefined) {
      this.updateFilterSequence(sequences);
      this.refreshFilterDataSource();
      // Needed for the updates to work properly
      this.changeDetectorRef_.detectChanges();
    }
  }

  public updateSequencesData(): void {
    if (this.project) {
      const sequences = this.project.sequences;
      if (sequences !== undefined) {
        // Update filteredSequences based on current filters
        this.filteredSequences = this.getFilteredSequences(
          sequences,
          this.apiService.getSequenceFilters(this.project.projectName)
        );
        // Set unfinished sequences so that the state updates can be loaded
        this.unfinishedSequences = sequences.filter((sequence: Sequence) => !sequence.isFinished());
      }
    }
  }

  public loadTraces(sequence: Sequence, eventId?: string): void {
    this._tracesTimer.unsubscribe();
    if (moment().subtract(1, 'day').isBefore(sequence.time)) {
      this._tracesTimer = AppUtils.createTimer(0, this._tracesTimerInterval)
        .pipe(takeUntil(this.unsubscribe$))
        .subscribe(() => {
          this.setTraces(sequence, eventId);
        });
    } else {
      this.setTraces(sequence, eventId);
      this._tracesTimer = Subscription.EMPTY;
    }
  }

  private setTraces(sequence: Sequence, eventId?: string): void {
    this.dataService.getTracesOfSequence(sequence).subscribe((traces) => {
      sequence.traces = traces;
      if (eventId) {
        const stage = sequence.findTrace((t) => t.id === eventId)?.stage;
        this.selectSequence({ sequence, stage, eventId }, false);
      }
    });
  }

  private updateLatestDeployedImage(): void {
    const deployedStage = this.metadata.deployments.find((depl) => depl.stage.name === this.selectedStage);
    const deployedService = deployedStage?.stage.services.find((svc) => svc.name === this.currentSequence?.service);
    this.currentLatestDeployedImage = deployedService?.image ?? '';
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  filtersClicked(event: DtQuickFilterChangeEvent<any> | { filters: any[] }): void {
    this._seqFilters = event.filters as FilterType[];
    const sequenceFilters: Record<string, string[]> = this._seqFilters.reduce(
      (
        filters: Record<string, string[]>,
        currentFilter: [
          { name: string; autocomplete: { name: string; value: string }[] },
          ...{ name: string; value: string }[]
        ]
      ) => {
        if (!filters[currentFilter[0].name]) {
          // Stage | Service | Sequence | Status
          filters[currentFilter[0].name] = [];
        }
        filters[currentFilter[0].name].push(currentFilter[1].value);
        return filters;
      },
      {}
    );
    this.saveSequenceFilters(sequenceFilters);
    this.updateSequencesData();
  }

  updateFilterSequence(sequences?: Sequence[]): void {
    if (sequences) {
      const filterItem = this.filterFieldData.autocomplete.find((f) => f.name === 'Sequence');
      if (filterItem) {
        filterItem.autocomplete = sequences
          .map((s) => s.name)
          .filter((v, i, a) => a.indexOf(v) === i)
          .map((seqName) =>
            Object.assign(
              {},
              {
                name: seqName,
                value: seqName,
              }
            )
          );
      }
    }
  }

  private mapServiceFilters(metadata: ISequencesMetadata): void {
    const filterItem = this.filterFieldData.autocomplete.find((f) => f.name === 'Service');
    if (filterItem) {
      // Take basis from metadatadata ...
      const serviceFilters: { name: string; value: string }[] = [];
      for (const svc of metadata.filter.services) {
        serviceFilters.push({ name: svc, value: svc });
      }

      // ... and enhance with sequence services (if deleted service has a sequence)
      if (this.project?.sequences) {
        this.project.sequences
          .map((s) => s.service)
          .filter((v, i, a) => a.indexOf(v) === i)
          .forEach((serviceName) => {
            if (serviceFilters.find((fltr) => fltr.name === serviceName) === undefined) {
              serviceFilters.push({ name: serviceName, value: serviceName });
            }
          });
      }

      filterItem.autocomplete = serviceFilters;

      // Remove service from active filters if not in list of services anymore
      this._seqFilters = this._seqFilters.filter(
        (fltr) => fltr[0].name !== 'Service' || serviceFilters.some((svc) => svc.name === fltr[1].name)
      );
    }
  }

  updateFilterDataSource(metadata: ISequencesMetadata): void {
    this.mapServiceFilters(metadata);

    const filterItem = this.filterFieldData.autocomplete.find((f) => f.name === 'Stage');
    if (filterItem) {
      filterItem.autocomplete = metadata.filter.stages.map((s) => {
        return { name: s, value: s };
      });
    }
    this.refreshFilterDataSource();
  }

  private refreshFilterDataSource(): void {
    this._filterDataSource = new DtQuickFilterDefaultDataSource(this.filterFieldData, this._config);
  }

  getFilteredSequences(sequences: Sequence[], filters: Record<string, string[]>): Sequence[] {
    const filterSequence = (s: Sequence): boolean => {
      const mapFilter = (key: string): boolean => {
        switch (key) {
          case 'Service':
            return filters[key].includes(s.service);
          case 'Stage':
            return filters[key].every((f) => s.getStages().includes(f));
          case 'Sequence':
            return filters[key].includes(s.name);
          case 'Status':
            return filters[key].includes(s.getStatus());
          default:
            return true;
        }
      };
      const reduceFilter = (prior: boolean, current: boolean): boolean => prior && current;
      return Object.keys(filters).map(mapFilter).reduce(reduceFilter, true);
    };

    return sequences.filter(filterSequence);
  }

  public getTracesLastUpdated(sequence: Sequence): Date | undefined {
    return this.dataService.getTracesLastUpdated(sequence);
  }

  showReloadButton(sequence: Sequence): boolean {
    return moment().subtract(1, 'day').isAfter(sequence.time);
  }

  selectStage(stageName: string): void {
    if (this.currentSequence) {
      const routeUrl = this.router.createUrlTree(
        ['/project', this.currentSequence.project, 'sequence', this.currentSequence.shkeptncontext, 'stage', stageName],
        { queryParamsHandling: 'preserve' }
      );
      this.location.go(routeUrl.toString());

      this.selectedStage = stageName;
      this.updateLatestDeployedImage();
    }
  }

  public navigateToTriggerSequence(): void {
    this.dataService.isTriggerSequenceOpen = true;
    this.router.navigate(['/project/' + this.project?.projectName]);
  }

  public saveSequenceFilters(sequenceFilters: { [p: string]: string[] }): void {
    if (this.project) {
      this.apiService.setSequenceFilters(sequenceFilters, this.project.projectName);
    }
    const routeUrl = this.router.createUrlTree([], {
      relativeTo: this.route,
      queryParams: sequenceFilters,
    });
    this.location.go(routeUrl.toString());
  }

  public loadSequenceFilters(): void {
    if (this.project) {
      this.router.navigate([], {
        relativeTo: this.route,
        queryParams: this.apiService.getSequenceFilters(this.project.projectName),
        replaceUrl: true,
      });
    }
  }

  public setSequenceFilters(sequenceFilters: { [p: string]: string[] }): void {
    if (this.project) {
      this.apiService.setSequenceFilters(sequenceFilters, this.project.projectName);
    }
    this._seqFilters = Object.keys(sequenceFilters).reduce((_seqFilters: FilterType[], filterName: string) => {
      sequenceFilters[filterName].forEach((value: string) => {
        const name = filterName === 'Status' ? SEQUENCE_STATUS[value] : value;

        if (!name) {
          return;
        }

        _seqFilters.push([
          { name: filterName as FilterName, autocomplete: [], showInSidebar: false },
          { name, value },
        ]);
      });
      return _seqFilters;
    }, []);

    this.updateSequencesData();
  }

  loadOldSequences(): void {
    if (this.project) {
      this.loading = true;
      this.dataService.loadOldSequences(this.project);
    }
  }

  ngOnDestroy(): void {
    this._tracesTimer.unsubscribe();
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
