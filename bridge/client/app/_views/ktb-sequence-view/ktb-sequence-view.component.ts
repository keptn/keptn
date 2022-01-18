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
import { filter, map, startWith, switchMap, takeUntil, takeWhile } from 'rxjs/operators';
import moment from 'moment';
import { Project } from '../../_models/project';
import { DataService } from '../../_services/data.service';
import { DateUtil } from '../../_utils/date.utils';
import { Sequence } from '../../_models/sequence';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { ISequencesMetadata, SequenceMetadataDeployment } from '../../../../shared/interfaces/sequencesMetadata';

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
        autocomplete: [
          { name: 'Active', value: 'started' },
          { name: 'Failed', value: 'failed' },
          { name: 'Aborted', value: 'aborted' },
          { name: 'Succeeded', value: 'succeeded' },
        ],
      },
    ],
  };
  private _config: DtQuickFilterDefaultDataSourceConfig = {
    // Method to decide if a node should be displayed in the quick filter
    showInSidebar: (node) => isObject(node) && node.showInSidebar,
  };
  private sequenceFilters: { [key: string]: string[] } = {};
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
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public _seqFilters: any[] = [];
  private latestDeployments: SequenceMetadataDeployment[] = [];

  constructor(
    private dataService: DataService,
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

    this.sequencesUpdated$ = this.dataService.sequencesUpdated.pipe(takeUntil(this.unsubscribe$));

    this.project$ = this.projectName$.pipe(switchMap((projectName) => this.dataService.getProject(projectName)));

    this.project$
      .pipe(
        takeUntil(this.unsubscribe$),
        filter(
          (project: Project | undefined): project is Project => !!project && !!project.getServices() && !!project.stages
        )
      )
      .subscribe((project) => {
        if (project.projectName !== this.project?.projectName) {
          this.currentSequence = undefined;
          this.selectedStage = undefined;
          this.filtersChanged({ filters: [] });
        }
        this.project = project;
      });
  }

  ngOnInit(): void {
    AppUtils.createTimer(0, this._sequenceTimerInterval)
      .pipe(
        startWith(0),
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
        takeWhile(([params]) => !this.currentSequence && params.shkeptncontext)
      )
      .subscribe(([params]: [Params, void]) => {
        if (params.shkeptncontext && this.project?.sequences) {
          const sequence = this.project.sequences.find((s) => s.shkeptncontext === params.shkeptncontext);
          const stage = params.eventId ? sequence?.traces.find((t) => t.id === params.eventId)?.stage : params.stage;
          const eventId = params.eventId;
          if (sequence) {
            this.selectSequence({ sequence, stage, eventId });
          } else if (params.shkeptncontext && this.project) {
            this.dataService.loadUntilRoot(this.project, params.shkeptncontext);
          }
        }
      });

    this.sequencesUpdated$.subscribe(() => {
      if (this.project?.sequences) {
        this.updateFilterSequence(this.project.sequences);
        this.refreshFilterDataSource();
        // Set unfinished sequences so that the state updates can be loaded
        this.unfinishedSequences = this.project.sequences.filter((sequence) => !sequence.isFinished());
        // Needed for the updates to work properly
        this.changeDetectorRef_.detectChanges();
      }
    });
  }

  public loadSequenceMetadata(projectName: string): void {
    this.dataService.getSequenceMetadata(projectName).subscribe((metadata) => {
      this.latestDeployments = metadata.deployments;
      this.updateLatestDeployedImage();
      this.updateFilterDataSource(metadata);
    });
  }

  selectSequence(event: { sequence: Sequence; stage?: string; eventId?: string }): void {
    if (event.eventId) {
      const routeUrl = this.router.createUrlTree([
        '/project',
        event.sequence.project,
        'sequence',
        event.sequence.shkeptncontext,
        'event',
        event.eventId,
      ]);
      this.location.go(routeUrl.toString());
    } else {
      const stage = event.stage || event.sequence.getStages().pop();
      const routeUrl = this.router.createUrlTree([
        '/project',
        event.sequence.project,
        'sequence',
        event.sequence.shkeptncontext,
        ...(stage ? ['stage', stage] : []),
      ]);
      this.location.go(routeUrl.toString());
    }

    this.currentSequence = event.sequence;
    this.selectedStage = event.stage || event.sequence.getStages().pop();
    this.updateLatestDeployedImage();
    this.loadTraces(this.currentSequence);
  }

  loadTraces(sequence: Sequence): void {
    this._tracesTimer.unsubscribe();
    if (moment().subtract(1, 'day').isBefore(sequence.time)) {
      this._tracesTimer = AppUtils.createTimer(0, this._tracesTimerInterval)
        .pipe(takeUntil(this.unsubscribe$))
        .subscribe(() => {
          this.dataService.loadTraces(sequence);
        });
    } else {
      this.dataService.loadTraces(sequence);
      this._tracesTimer = Subscription.EMPTY;
    }
  }

  private updateLatestDeployedImage(): void {
    const deployedStage = this.latestDeployments.find((depl) => depl.stage.name === this.selectedStage);
    const deployedService = deployedStage?.stage.services.find((svc) => svc.name === this.currentSequence?.service);
    this.currentLatestDeployedImage = deployedService?.image ?? '';
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  filtersChanged(event: DtQuickFilterChangeEvent<any> | { filters: [] }): void {
    this._seqFilters = event.filters;
    this.sequenceFilters = this._seqFilters.reduce((filters, currentFilter) => {
      if (!filters[currentFilter[0].name]) {
        filters[currentFilter[0].name] = [];
      }
      filters[currentFilter[0].name].push(currentFilter[1].value);
      return filters;
    }, {});
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

  getFilteredSequences(sequences: Sequence[]): Sequence[] {
    return sequences.filter((s) => {
      let res = true;
      Object.keys(this.sequenceFilters).forEach((key) => {
        switch (key) {
          case 'Service':
            res = res && this.sequenceFilters[key].includes(s.service);
            break;
          case 'Stage':
            res = res && this.sequenceFilters[key].every((f) => s.getStages().includes(f));
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

  showReloadButton(sequence: Sequence): boolean {
    return moment().subtract(1, 'day').isAfter(sequence.time);
  }

  selectStage(stageName: string): void {
    if (this.currentSequence) {
      const routeUrl = this.router.createUrlTree([
        '/project',
        this.currentSequence.project,
        'sequence',
        this.currentSequence.shkeptncontext,
        'stage',
        stageName,
      ]);
      this.location.go(routeUrl.toString());

      this.selectedStage = stageName;
      this.updateLatestDeployedImage();
    }
  }

  ngOnDestroy(): void {
    this._tracesTimer.unsubscribe();
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
