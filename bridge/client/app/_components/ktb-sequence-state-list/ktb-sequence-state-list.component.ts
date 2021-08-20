import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Inject, Input, NgZone, OnDestroy, ViewEncapsulation } from '@angular/core';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { DateUtil } from '../../_utils/date.utils';
import { DataService } from '../../_services/data.service';
import { Sequence } from '../../_models/sequence';
import { Subscription } from 'rxjs';
import { Project } from '../../_models/project';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';

@Component({
  selector: 'ktb-sequence-state-list',
  templateUrl: './ktb-sequence-state-list.component.html',
  styleUrls: ['./ktb-sequence-state-list.component.scss'],
  host: {
    class: 'ktb-sequence-state-list',
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
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
      this.ngZone.runOutsideAngular(() => {
        this._timer = AppUtils.createTimer(0, this.initialDelayMillis)
          .subscribe(() => {
            this.loadLatestSequences();
          });
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
      this._changeDetectorRef.detectChanges();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, public dataService: DataService, public dateUtil: DateUtil, private ngZone: NgZone, @Inject(POLLING_INTERVAL_MILLIS) private initialDelayMillis: number) {
  }

  loadLatestSequences() {
    if (this.project) {
      this.dataService.loadLatestSequences(this.project, this.PAGE_SIZE)
        .subscribe((sequences: Sequence[]) => {
          this.sequenceStates = sequences;
        });
    }
  }

  updateDataSource() {
    this.dataSource = new DtTableDataSource(this.sequenceStates.slice(0, this.PAGE_SIZE) || []);
  }

  getServiceLink(sequence: Sequence) {
    return ['/project', sequence.project, 'service', sequence.service, 'context', sequence.shkeptncontext, 'stage', sequence.getLastStage()];
  }

  getSequenceLink(sequence: Sequence) {
    return ['/project', sequence.project, 'sequence', sequence.shkeptncontext, 'stage', sequence.getLastStage()];
  }

  ngOnDestroy(): void {
    this._timer.unsubscribe();
  }

}
