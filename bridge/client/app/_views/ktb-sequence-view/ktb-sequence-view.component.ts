import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnDestroy,
  OnInit,
  ViewEncapsulation
} from '@angular/core';
import {Location} from '@angular/common';
import {ActivatedRoute, Router} from '@angular/router';
import {DtQuickFilterDefaultDataSource, DtQuickFilterDefaultDataSourceConfig} from '@dynatrace/barista-components/experimental/quick-filter';
import {isObject} from '@dynatrace/barista-components/core';

import {Observable, Subject, Subscription, timer} from 'rxjs';
import {filter, take, takeUntil} from 'rxjs/operators';

import * as moment from 'moment';

import {Root} from '../../_models/root';
import {Stage} from '../../_models/stage';
import {Project} from '../../_models/project';

import {DataService} from '../../_services/data.service';
import {DateUtil} from '../../_utils/date.utils';

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
  /** configuration for the quick filter **/
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
          { name: 'Active', value: 'active' },
          { name: 'Failed', value: 'failed' },
          { name: 'Succeeded', value: 'succeeded' },
        ],
      },
    ],
  };
  private _config: DtQuickFilterDefaultDataSourceConfig = {
    // Method to decide if a node should be displayed in the quick filter
    showInSidebar: (node) => isObject(node) && node.showInSidebar,
  };
  private sequenceFilters = {};
  private project: Project;
  private loading = false;

  private _tracesTimerInterval = 10;
  private _tracesTimer: Subscription = Subscription.EMPTY;

  public project$: Observable<Project>;
  public currentSequence: Root;
  public selectedStage: String;

  public _filterDataSource = new DtQuickFilterDefaultDataSource(
    this.filterFieldData,
    this._config,
  );
  public _seqFilters = [];

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute, public dateUtil: DateUtil, private router: Router, private location: Location) { }

  ngOnInit() {
    this.currentSequence = null;
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        this.project$ = this.dataService.getProject(params.projectName);
        this.project$
          .pipe(
            filter(project => !!project && !!project.getServices() && !!project.stages),
            take(1)
          )
          .subscribe(project => {
            this.currentSequence = null;
            this.selectedStage = null;
            this.project = project;
            this.updateFilterDataSource(project);

            this._changeDetectorRef.markForCheck();
          });

        this.dataService.roots
          .pipe(
            takeUntil(this.unsubscribe$),
            filter(roots => !!roots)
          )
          .subscribe(roots => {
            if (!this.currentSequence && roots && params.shkeptncontext) {
              const root = roots.find(sequence => sequence.shkeptncontext === params.shkeptncontext);
              if (root) {
                this.selectSequence({root, stage: params.stage});
              } else if (!this.loading) {
                this.loading = true;
                this.dataService.loadUntilRoot(this.project, params.shkeptncontext);
              }
            }
            if (roots) {
              this.updateFilterSequence(roots);
              this._filterDataSource.data = this.filterFieldData;
            }
            this._changeDetectorRef.markForCheck();
          });
      });
  }

  selectSequence(event: {root: Root, stage: string}): void {
    if (event.stage) {
      const routeUrl = this.router.createUrlTree(['/project', event.root.getProject(), 'sequence', event.root.shkeptncontext, 'stage', event.stage]);
      this.location.go(routeUrl.toString());
    } else {
      const routeUrl = this.router.createUrlTree(['/project', event.root.getProject(), 'sequence', event.root.shkeptncontext]);
      this.location.go(routeUrl.toString());
    }

    this.currentSequence = event.root;
    this.selectedStage = event.stage || event.root.getStages().pop();
    this.loadTraces(this.currentSequence);
  }

  loadTraces(root: Root): void {
    this._tracesTimer.unsubscribe();
    if(moment().subtract(1, 'day').isBefore(root.time)) {
      this._tracesTimer = timer(0, this._tracesTimerInterval*1000)
        .subscribe(() => {
          this.dataService.loadTraces(root);
        });
    } else {
      this.dataService.loadTraces(root);
      this._tracesTimer = Subscription.EMPTY;
    }
  }

  filtersChanged(event) {
    this._seqFilters = event.filters;
    this.sequenceFilters = this._seqFilters.reduce((filters, filter) => {
      if(!filters[filter[0].name])
        filters[filter[0].name] = [];
      filters[filter[0].name].push(filter[1].value);
      return filters;
    }, {});
  }

  updateFilterSequence(sequences: Root[]) {
    if (sequences) {
      this.filterFieldData.autocomplete.find(f => f.name == 'Sequence').autocomplete = sequences.map(s => s.getShortType()).filter((v, i, a) => a.indexOf(v) === i).map(seqName => Object.assign({}, {
        name: seqName,
        value: seqName
      }));
    }
  }

  updateFilterDataSource(project: Project) {
    this.filterFieldData.autocomplete.find(f => f.name == 'Service').autocomplete = project.services.map(s => Object.assign({}, { name: s.serviceName, value: s.serviceName }));
    this.filterFieldData.autocomplete.find(f => f.name == 'Stage').autocomplete = project.stages.map(s => Object.assign({}, { name: s.stageName, value: s.stageName }));
    this.updateFilterSequence(project.sequences);
    this._filterDataSource.data = this.filterFieldData;

    this.filtersChanged({ filters: [] });
    this._changeDetectorRef.markForCheck();
  }

  getFilteredSequences(sequences: Root[]) {
    if(sequences)
      return sequences.filter(s => {
        let res = true;
        Object.keys(this.sequenceFilters||{}).forEach((key) => {
          switch(key) {
            case "Service":
              res = res && this.sequenceFilters[key].includes(s.getService());
              break;
            case "Stage":
              res = res && this.sequenceFilters[key].every(f => s.getStages().includes(f));
              break;
            case "Sequence":
              res = res && this.sequenceFilters[key].includes(s.getShortType());
              break;
            case "Status":
              res = res && this.sequenceFilters[key].includes(s.getStatus());
              break;
          }
        });
        return res;
      });
  }

  getTracesLastUpdated(root: Root): Date {
    return this.dataService.getTracesLastUpdated(root);
  }

  showReloadButton(root: Root) {
    return moment().subtract(1, 'day').isAfter(root.time);
  }

  selectStage(event: {stageName: string, triggerByEvent: boolean}) {
    if (!event.triggerByEvent) {
      const routeUrl = this.router.createUrlTree(['/project', this.currentSequence.getProject(), 'sequence', this.currentSequence.shkeptncontext, 'stage', event.stageName]);
      this.location.go(routeUrl.toString());
    }

    this.selectedStage = event.stageName;
    this._changeDetectorRef.markForCheck();
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
