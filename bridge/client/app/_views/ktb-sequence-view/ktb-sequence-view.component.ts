import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnDestroy,
  OnInit,
  ViewEncapsulation
} from '@angular/core';
import {Location} from '@angular/common';
import {ActivatedRoute, Params, Router} from '@angular/router';
import { DtQuickFilterChangeEvent, DtQuickFilterDefaultDataSource, DtQuickFilterDefaultDataSourceConfig } from '@dynatrace/barista-components/quick-filter';
import {isObject} from '@dynatrace/barista-components/core';
import {combineLatest, Observable, Subject, Subscription, timer} from 'rxjs';
import { filter, map, startWith, switchMap, takeUntil, takeWhile } from 'rxjs/operators';
import moment from 'moment';
import {Project} from '../../_models/project';
import {DataService} from '../../_services/data.service';
import {DateUtil} from '../../_utils/date.utils';
import {Sequence} from '../../_models/sequence';

@Component({
  selector: 'ktb-sequence-view',
  templateUrl: './ktb-sequence-view.component.html',
  styleUrls: ['./ktb-sequence-view.component.scss'],
  host: {
    class: 'ktb-sequence-view'
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSequenceViewComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();
  /** configuration for the quick filter */
  private filterFieldData = {
    autocomplete: [
      {
        name: 'Service',
        showInSidebar: true,
        autocomplete: [],
      }, {
        name: 'Stage',
        showInSidebar: true,
        autocomplete: [],
      }, {
        name: 'Sequence',
        showInSidebar: true,
        autocomplete: [
        ],
      }, {
        name: 'Status',
        showInSidebar: true,
        autocomplete: [
          { name: 'Active', value: 'started' },
          { name: 'Failed', value: 'failed' },
          { name: 'Succeeded', value: 'succeeded' },
          { name: 'Waiting', value: 'waiting' }
        ],
      },
    ],
  };
  private _config: DtQuickFilterDefaultDataSourceConfig = {
    // Method to decide if a node should be displayed in the quick filter
    showInSidebar: (node) => isObject(node) && node.showInSidebar,
  };
  private sequenceFilters: {[key: string]: string[]} = {};
  private project?: Project;
  private unfinishedSequences: Sequence[] = [];
  private _tracesTimerInterval = 10_000;
  private _sequenceTimerInterval = 30_000;
  private _tracesTimer: Subscription = Subscription.EMPTY;
  private _rootsTimer: Subscription = Subscription.EMPTY;

  public project$: Observable<Project | undefined>;
  public sequences$: Observable<Sequence[]>;
  public currentSequence?: Sequence;
  public selectedStage?: string;
  public _filterDataSource = new DtQuickFilterDefaultDataSource(
    this.filterFieldData,
    this._config,
  );
  // tslint:disable-next-line:no-any
  public _seqFilters: any[] = [];

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute,
              public dateUtil: DateUtil, private router: Router, private location: Location) {
    const projectName$ = this.route.params
      .pipe(
        map(params => params.projectName)
      );

    this.sequences$ = this.dataService.sequences
      .pipe(
        takeUntil(this.unsubscribe$),
        filter((sequences: Sequence[] | undefined): sequences is Sequence[] => !!sequences?.length)
      );

    this.project$ = projectName$.pipe(
      switchMap(projectName => this.dataService.getProject(projectName))
    );

    this.project$
      .pipe(
        takeUntil(this.unsubscribe$),
        filter((project: Project | undefined): project is Project => !!project && !!project.getServices() && !!project.stages)
      )
      .subscribe(project => {
        if (project.projectName !== this.project?.projectName) {
          this.currentSequence = undefined;
          this.selectedStage = undefined;
          this.updateFilterDataSource(project);
          this.dataService.loadRoots(project);
        }
        this.project = project;
        this._changeDetectorRef.markForCheck();
      });
  }

  ngOnInit() {


    timer(0, this._sequenceTimerInterval)
      .pipe(
        startWith(0),
        switchMap(() => this.project$),
        filter((project: Project | undefined): project is Project => !!project && !!project.getServices()),
        takeUntil(this.unsubscribe$)
      ).subscribe(project => {
      this.dataService.loadSequences(project);
    });

    this._rootsTimer = timer(0, this._tracesTimerInterval)
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(() => {
        // This triggers the subscription for sequences$
        if (this.project) {
          for (const sequence of this.unfinishedSequences) {
            this.dataService.updateSequence(this.project.projectName, sequence.shkeptncontext);
          }
        }
    });

    // init; set parameters
    combineLatest([this.route.params, this.sequences$])
      .pipe(
        takeUntil(this.unsubscribe$),
        takeWhile( ([params]) => !this.currentSequence && params.shkeptncontext)
      )
      .subscribe(([params, sequences]: [Params, Sequence[]]) => {
        if (params.shkeptncontext) {
          const sequence = sequences.find(s => s.shkeptncontext === params.shkeptncontext);
          const stage = params.eventId ? sequence?.traces.find(t => t.id === params.eventId)?.stage : params.stage;
          const eventId = params.eventId;
          if (sequence) {
            this.selectSequence({sequence, stage, eventId});
          } else if (params.shkeptncontext && this.project) {
            this.dataService.loadUntilRoot(this.project, params.shkeptncontext);
          }
        }
    });

    this.sequences$.subscribe(sequences => {
      this.updateFilterSequence(sequences);
      this.refreshFilterDataSource();
      // Set unfinished sequences so that the state updates can be loaded
      this.unfinishedSequences = sequences.filter(sequence => !sequence.isFinished());
    });
  }

  selectSequence(event: {sequence: Sequence, stage: string, eventId: string}): void {
    if (event.eventId) {
      const routeUrl = this.router.createUrlTree(['/project', event.sequence.project, 'sequence',
                                                          event.sequence.shkeptncontext, 'event', event.eventId]);
      this.location.go(routeUrl.toString());
    } else {
      const stage = event.stage || event.sequence.getStages().pop();
      const routeUrl = this.router.createUrlTree(['/project', event.sequence.project, 'sequence', event.sequence.shkeptncontext,
                                                            ...(stage ? ['stage', stage] : [])]);
      this.location.go(routeUrl.toString());
    }

    this.currentSequence = event.sequence;
    this.selectedStage = event.stage || event.sequence.getStages().pop();
    this.loadTraces(this.currentSequence);
  }

  loadTraces(sequence: Sequence): void {
    this._tracesTimer.unsubscribe();
    if (moment().subtract(1, 'day').isBefore(sequence.time)) {
      this._tracesTimer = timer(0, this._tracesTimerInterval)
        .pipe(takeUntil(this.unsubscribe$))
        .subscribe(() => {
          this.dataService.loadTraces(sequence);
        });
    } else {
      this.dataService.loadTraces(sequence);
      this._tracesTimer = Subscription.EMPTY;
    }
  }

  // tslint:disable-next-line:no-any
  filtersChanged(event: DtQuickFilterChangeEvent<any> | {filters: []}) {
    this._seqFilters = event.filters;
    this.sequenceFilters = this._seqFilters.reduce((filters, currentFilter) => {
      if (!filters[currentFilter[0].name]) {
        filters[currentFilter[0].name] = [];
      }
      filters[currentFilter[0].name].push(currentFilter[1].value);
      return filters;
    }, {});
  }

  updateFilterSequence(sequences: Sequence[]) {
    if (sequences) {
      const filterItem = this.filterFieldData.autocomplete.find(f => f.name === 'Sequence');
      if (filterItem) {
        filterItem.autocomplete = sequences.map(s => s.name).filter((v, i, a) => a.indexOf(v) === i).map(seqName => Object.assign({}, {
          name: seqName,
          value: seqName
        }));
      }
    }
  }

  updateFilterDataSource(project: Project) {
    let filterItem  = this.filterFieldData.autocomplete.find(f => f.name === 'Service');
    if (filterItem) {
      filterItem.autocomplete = project.getServices().map(s => Object.assign({}, {name: s.serviceName, value: s.serviceName}));
    }
    filterItem = this.filterFieldData.autocomplete.find(f => f.name === 'Stage');
    if (filterItem) {
      filterItem.autocomplete = project.stages.map(s => Object.assign({}, {name: s.stageName, value: s.stageName}));
    }
    this.updateFilterSequence(project.sequences);
    this.refreshFilterDataSource();

    this.filtersChanged({ filters: [] });
    this._changeDetectorRef.markForCheck();
  }

  private refreshFilterDataSource() {
    this._filterDataSource = new DtQuickFilterDefaultDataSource(
      this.filterFieldData,
      this._config,
    );
  }

  getFilteredSequences(sequences: Sequence[]): Sequence[] {
    return sequences.filter(s => {
      let res = true;
      Object.keys(this.sequenceFilters).forEach((key) => {
        switch (key) {
          case 'Service':
            res = res && this.sequenceFilters[key].includes(s.service);
            break;
          case 'Stage':
            res = res && this.sequenceFilters[key].every(f => s.getStages().includes(f));
            break;
          case 'Sequence':
            res = res && this.sequenceFilters[key].includes(s.name);
            break;
          case 'Status':
            res = res && this.sequenceFilters[key].includes(s.getStatus());
            break;
          default:
            break;
        }
      });
      return res;
    });
  }

  getTracesLastUpdated(sequence: Sequence): Date {
    return this.dataService.getTracesLastUpdated(sequence);
  }

  showReloadButton(sequence: Sequence) {
    return moment().subtract(1, 'day').isAfter(sequence.time);
  }

  selectStage(stageName: string) {
    if (this.currentSequence) {
      const routeUrl = this.router.createUrlTree(['/project', this.currentSequence.project, 'sequence', this.currentSequence.shkeptncontext, 'stage', stageName]);
      this.location.go(routeUrl.toString());

      this.selectedStage = stageName;
      this._changeDetectorRef.markForCheck();
    }
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
    this._tracesTimer.unsubscribe();
    this._rootsTimer.unsubscribe();
  }
}
