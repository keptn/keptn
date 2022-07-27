import { ChangeDetectorRef, Component, HostBinding, Inject, OnDestroy } from '@angular/core';
import { Location } from '@angular/common';
import { ActivatedRoute, ParamMap, Router } from '@angular/router';
import {
  DtQuickFilterChangeEvent,
  DtQuickFilterDefaultDataSource,
  DtQuickFilterDefaultDataSourceConfig,
} from '@dynatrace/barista-components/quick-filter';
import { isObject } from '@dynatrace/barista-components/core';
import { combineLatest, Observable, of, Subject, Subscription } from 'rxjs';
import { distinctUntilChanged, filter, map, switchMap, take, takeUntil, tap } from 'rxjs/operators';
import moment from 'moment';
import { DataService } from '../../_services/data.service';
import { DateUtil } from '../../_utils/date.utils';
import { Sequence } from '../../_models/sequence';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { ISequencesFilter } from '../../../../shared/interfaces/sequencesFilter';
import { ApiService } from '../../_services/api.service';
import { isEqual } from 'lodash-es';
import { FilterName, FilterType, ISequenceViewState, SequencesState } from './ktb-sequence-view.utils';

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
})
export class KtbSequenceViewComponent implements OnDestroy {
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
  private unfinishedSequences: Sequence[] = [];
  private readonly _tracesTimerInterval: number = 10_000;
  private readonly _sequenceTimerInterval: number = 30_000;
  private _tracesTimer: Subscription = Subscription.EMPTY;
  private sequences: Sequence[] = [];
  public currentSequence?: Sequence;
  public selectedStage?: string;
  public _filterDataSource = new DtQuickFilterDefaultDataSource(this.filterFieldData, this._config);
  public _seqFilters: FilterType[] = [];
  public metadata: ISequencesFilter = {
    stages: [],
    services: [],
  };
  public filteredSequences?: Sequence[];
  public loading = false;
  public state$: Observable<ISequenceViewState>;

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

    const eventId$ = this.route.paramMap.pipe(map((params) => params.get('eventId')));
    const projectName$ = this.route.paramMap.pipe(
      map((params) => params.get('projectName')),
      filter((projectName): projectName is string => !!projectName),
      distinctUntilChanged(),
      takeUntil(this.unsubscribe$)
    );

    const state$ = combineLatest([projectName$, eventId$]).pipe(
      switchMap(([projectName, eventId]) =>
        this.dataService
          .getSequences(projectName)
          .pipe(map((sequenceInfo) => ({ projectName, sequenceInfo, eventId: eventId ?? undefined })))
      )
    );

    this.state$ = state$.pipe(
      tap((state) => {
        this.sequences = state.sequenceInfo?.sequences ?? [];
        if (!state.sequenceInfo) {
          return;
        }
        this.loading = false;
        this.updateFilterSequence(state.sequenceInfo.sequences);
        this.updateSequencesData(state.sequenceInfo.sequences, state.projectName);
        // Needed for the updates to work properly
        this.changeDetectorRef_.detectChanges();
      })
    );

    // not inside tap because it is additionally subscribed somewhere else and then tab would be called way too often
    projectName$.subscribe((projectName) => {
      this.currentSequence = undefined;
      this.selectedStage = undefined;
      this.loadSequenceMetadata(projectName);
    });

    // route params: select sequence (and event/stage)
    const loadedState$ = state$.pipe(filter((state) => state.sequenceInfo?.state === SequencesState.UPDATE));
    combineLatest([this.route.paramMap, loadedState$])
      .pipe(takeUntil(this.unsubscribe$), take(1))
      .subscribe(([params, state]) => {
        this.setParams(params, state);
      });

    // route params: set params after sequence is loaded that wasn't initially loaded
    const loadUntilRootTriggered$ = state$.pipe(
      filter((state) => state.sequenceInfo?.state === SequencesState.LOAD_UNTIL_ROOT)
    );
    combineLatest([this.route.paramMap, loadUntilRootTriggered$])
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(([params, state]) => {
        this.setParams(params, state);
      });

    // set filter through query params
    combineLatest([this.route.queryParams, projectName$])
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(([queryParams, projectName]) => {
        if (Object.keys(queryParams).length === 0) {
          this.loadSequenceFilters(projectName);
          return;
        }
        const sequenceFilters = Object.keys(queryParams).reduce((params: Record<string, string[]>, param) => {
          params[param] = Array.isArray(queryParams[param]) ? queryParams[param] : [queryParams[param]];
          return params;
        }, {});
        this.setSequenceFilters(sequenceFilters, projectName);
      });

    // fetch new sequences
    AppUtils.createTimer(0, this._sequenceTimerInterval)
      .pipe(
        switchMap(() => projectName$),
        takeUntil(this.unsubscribe$)
      )
      .subscribe((projectName) => {
        this.dataService.loadSequences(projectName);
      });

    // update unfinished sequences
    AppUtils.createTimer(0, this._tracesTimerInterval)
      .pipe(
        switchMap(() => projectName$),
        takeUntil(this.unsubscribe$)
      )
      .subscribe((projectName) => {
        // This triggers the subscription for sequences$
        for (const sequence of this.unfinishedSequences) {
          this.dataService.updateSequence(projectName, sequence.shkeptncontext);
        }
      });
  }

  private setParams(params: ParamMap, state: ISequenceViewState): void {
    const keptnContext = params.get('shkeptncontext');
    if (!keptnContext) {
      return;
    }
    const sequence = state.sequenceInfo?.sequences.find((s) => s.shkeptncontext === keptnContext);
    const eventId = params.get('eventId') ?? undefined;
    const stage = eventId ? undefined : params.get('stage') ?? undefined;

    if (sequence) {
      if (eventId && !sequence.traces.length) {
        this.loadTraces(sequence, eventId);
        return;
      }

      // while traces are loading we can already show the sequence and the timeline
      this.selectSequence({ sequence, eventId, stage });
      this.loadTraces(sequence, eventId, stage);
      return;
    }

    this.dataService.loadUntilRoot(state.projectName, keptnContext);
  }

  public loadSequenceMetadata(projectName: string): void {
    this.dataService.getSequenceFilter(projectName).subscribe((metadata) => {
      this.metadata = metadata;
      this.updateFilterDataSource(metadata, this.sequences);
    });
  }

  public selectSequence(event: { sequence: Sequence; stage?: string; eventId?: string }): void {
    const sequenceFilters = this.apiService.getSequenceFilters(event.sequence.project);
    let stage = event.stage || event.sequence.getStages().pop();
    const additionalCommands = [];
    if (event.eventId) {
      stage = event.sequence.findTrace((t) => t.id === event.eventId)?.stage;
      additionalCommands.push('event', event.eventId);
    } else if (stage) {
      additionalCommands.push('stage', stage);
    }
    const routeUrl = this.router.createUrlTree(
      ['/project', event.sequence.project, 'sequence', event.sequence.shkeptncontext, ...additionalCommands],
      { queryParams: sequenceFilters }
    );
    this.location.go(routeUrl.toString());
    this.currentSequence = event.sequence;
    this.selectedStage = stage;
  }

  public updateSequencesData(sequences: Sequence[], projectName: string): void {
    // Update filteredSequences based on current filters
    this.filteredSequences = this.getFilteredSequences(sequences, this.apiService.getSequenceFilters(projectName));
    // Set unfinished sequences so that the state updates can be loaded
    this.unfinishedSequences = sequences.filter((sequence: Sequence) => !sequence.isFinished());
  }

  public loadTraces(sequence: Sequence, eventId?: string, stage?: string): void {
    this._tracesTimer.unsubscribe();
    let setTraces$;
    if (moment().subtract(1, 'day').isBefore(sequence.time)) {
      setTraces$ = AppUtils.createTimer(0, this._tracesTimerInterval);
    } else {
      setTraces$ = of(null);
    }
    this._tracesTimer = setTraces$.pipe(takeUntil(this.unsubscribe$)).subscribe(() => {
      this.setTraces(sequence, eventId, stage);
    });
  }

  private setTraces(sequence: Sequence, eventId?: string, stage?: string): void {
    this.dataService.getTracesOfSequence(sequence).subscribe((traces) => {
      sequence.traces = traces;
      this.selectSequence({ sequence, stage, eventId });
    });
  }

  public filtersClicked(
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    event: DtQuickFilterChangeEvent<any> | { filters: any[] },
    sequences: Sequence[],
    projectName: string
  ): void {
    this._seqFilters = event.filters as FilterType[];
    const sequenceFilters: Record<string, string[]> = this._seqFilters.reduce(
      (filters: Record<string, string[]>, currentFilter: FilterType) => {
        if (!filters[currentFilter[0].name]) {
          // Stage | Service | Sequence | Status
          filters[currentFilter[0].name] = [];
        }
        filters[currentFilter[0].name].push(currentFilter[1].value);
        return filters;
      },
      {}
    );
    this.saveSequenceFilters(sequenceFilters, projectName);
    this.updateSequencesData(sequences, projectName);
  }

  private updateFilterSequence(sequences: Sequence[]): void {
    const filterItem = this.filterFieldData.autocomplete.find((f) => f.name === 'Sequence');
    if (!filterItem) {
      return;
    }
    const newFilter = sequences
      .map((s) => s.name)
      .filter((v, i, a) => a.indexOf(v) === i)
      .map((seqName) => ({
        name: seqName,
        value: seqName,
      }));
    if (isEqual(newFilter, filterItem.autocomplete)) {
      return;
    }

    // only update the filter and refresh it if there are changes
    filterItem.autocomplete = newFilter;
    this.refreshFilterDataSource();
  }

  private mapServiceFilters(metadata: ISequencesFilter, sequences: Sequence[]): void {
    const filterItem = this.filterFieldData.autocomplete.find((f) => f.name === 'Service');
    if (filterItem) {
      // Take basis from metadatadata ...
      const serviceFilters: { name: string; value: string }[] = [];
      for (const svc of metadata.services) {
        serviceFilters.push({ name: svc, value: svc });
      }

      // ... and enhance with sequence services (if deleted service has a sequence)
      sequences
        .map((s) => s.service)
        .filter((v, i, a) => a.indexOf(v) === i)
        .forEach((serviceName) => {
          if (!serviceFilters.some((fltr) => fltr.name === serviceName)) {
            serviceFilters.push({ name: serviceName, value: serviceName });
          }
        });

      filterItem.autocomplete = serviceFilters;

      // Remove service from active filters if not in list of services anymore
      this._seqFilters = this._seqFilters.filter(
        (fltr) => fltr[0].name !== 'Service' || serviceFilters.some((svc) => svc.name === fltr[1].name)
      );
    }
  }

  private updateFilterDataSource(metadata: ISequencesFilter, sequences: Sequence[]): void {
    this.mapServiceFilters(metadata, sequences);

    const filterItem = this.filterFieldData.autocomplete.find((f) => f.name === 'Stage');
    if (filterItem) {
      filterItem.autocomplete = metadata.stages.map((s) => {
        return { name: s, value: s };
      });
    }
    this.refreshFilterDataSource();
  }

  private refreshFilterDataSource(): void {
    this._filterDataSource = new DtQuickFilterDefaultDataSource(this.filterFieldData, this._config);
  }

  private getFilteredSequences(sequences: Sequence[], filters: Record<string, string[]>): Sequence[] {
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

  public showReloadButton(sequence: Sequence): boolean {
    return moment().subtract(1, 'day').isAfter(sequence.time);
  }

  public selectStage(stageName: string): void {
    if (!this.currentSequence) {
      return;
    }
    const sequenceFilters = this.apiService.getSequenceFilters(this.currentSequence.project);
    const routeUrl = this.router.createUrlTree(
      ['/project', this.currentSequence.project, 'sequence', this.currentSequence.shkeptncontext, 'stage', stageName],
      { queryParams: sequenceFilters }
    );
    this.location.go(routeUrl.toString());
    this.selectedStage = stageName;
  }

  public navigateToTriggerSequence(projectName: string): void {
    this.dataService.setIsTriggerSequenceOpen(true);
    this.router.navigate(['/project/' + projectName]);
  }

  public saveSequenceFilters(sequenceFilters: { [p: string]: string[] }, projectName: string): void {
    this.apiService.setSequenceFilters(sequenceFilters, projectName);
    const routeUrl = this.router.createUrlTree([], {
      relativeTo: this.route,
      queryParams: sequenceFilters,
    });
    this.location.go(routeUrl.toString());
  }

  public loadSequenceFilters(projectName: string): void {
    const queryParams = this.apiService.getSequenceFilters(projectName);
    if (!Object.keys(queryParams).length) {
      return;
    }
    this.router.navigate(['project', projectName, 'sequence'], {
      queryParams,
      replaceUrl: true,
    });
  }

  public setSequenceFilters(sequenceFilters: Record<string, string[]>, projectName: string): void {
    this.apiService.setSequenceFilters(sequenceFilters, projectName);
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

    this.updateSequencesData(this.sequences, projectName);
  }

  public loadOldSequences(projectName: string): void {
    this.loading = true;
    this.dataService.loadOldSequences(projectName);
  }

  public ngOnDestroy(): void {
    this._tracesTimer.unsubscribe();
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
