import { Component, Inject, Input, NgZone, OnDestroy, ViewEncapsulation } from '@angular/core';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { DateUtil } from '../../_utils/date.utils';
import { DataService } from '../../_services/data.service';
import { Sequence } from '../../_models/sequence';
import { Subscription } from 'rxjs';
import { Project } from '../../_models/project';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { Router } from '@angular/router';

@Component({
  selector: 'ktb-sequence-state-list',
  templateUrl: './ktb-sequence-state-list.component.html',
  styleUrls: ['./ktb-sequence-state-list.component.scss'],
  host: {
    class: 'ktb-sequence-state-list',
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
})
export class KtbSequenceStateListComponent implements OnDestroy {
  private _project?: Project;
  private _sequenceStates: Sequence[] = [];
  private _timer: Subscription = Subscription.EMPTY;
  public dataSource: DtTableDataSource<Sequence> = new DtTableDataSource();
  public SequenceClass = Sequence;
  public PAGE_SIZE = 5;

  @Input()
  get project(): Project | undefined {
    return this._project;
  }

  set project(value: Project | undefined) {
    if (this._project !== value) {
      this._project = value;
      this._timer.unsubscribe();
      this._timer = AppUtils.createTimer(0, this.initialDelayMillis).subscribe(() => {
        this.loadLatestSequences();
      });
    }
  }

  get sequenceStates(): Sequence[] {
    return this._sequenceStates;
  }

  set sequenceStates(value: Sequence[]) {
    if (this._sequenceStates !== value) {
      this._sequenceStates = value;
      this.updateDataSource();
    }
  }

  constructor(
    public dataService: DataService,
    public dateUtil: DateUtil,
    private ngZone: NgZone,
    @Inject(POLLING_INTERVAL_MILLIS) private initialDelayMillis: number,
    private router: Router
  ) {}

  loadLatestSequences(): void {
    if (this.project) {
      this.dataService.loadLatestSequences(this.project, this.PAGE_SIZE).subscribe((sequences: Sequence[]) => {
        this.sequenceStates = sequences;
      });
    }
  }

  updateDataSource(): void {
    this.dataSource = new DtTableDataSource(this.sequenceStates.slice(0, this.PAGE_SIZE) || []);
  }

  selectSequence(event: { sequence: Sequence; stage?: string }): void {
    const stage = event.stage || event.sequence.getStages().pop();
    this.router.navigate([
      '/project',
      event.sequence.project,
      'sequence',
      event.sequence.shkeptncontext,
      ...(stage ? ['stage', stage] : []),
    ]);
  }

  ngOnDestroy(): void {
    this._timer.unsubscribe();
  }
}
