import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Input, OnDestroy, OnInit,
  ViewEncapsulation
} from '@angular/core';
import {DtTableDataSource} from "@dynatrace/barista-components/table";

import {DateUtil} from "../../_utils/date.utils";
import {DataService} from "../../_services/data.service";
import {Sequence} from "../../_models/sequence";
import {Subject} from "rxjs";
import {takeUntil} from "rxjs/operators";

@Component({
  selector: 'ktb-sequence-state-list',
  templateUrl: './ktb-sequence-state-list.component.html',
  styleUrls: ['./ktb-sequence-state-list.component.scss'],
  host: {
    class: 'ktb-sequence-state-list'
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSequenceStateListComponent implements OnInit, OnDestroy {

  private readonly unsubscribe$ = new Subject<void>();

  public _sequenceStates: Sequence[] = [];
  public dataSource: DtTableDataSource<Sequence> = new DtTableDataSource();

  public PAGE_SIZE = 5;

  @Input()
  get sequenceStates(): Sequence[] {
    return this._sequenceStates;
  }
  set sequenceStates(value: Sequence[]) {
    if (this._sequenceStates !== value) {
      this._sequenceStates = value;
      this.updateDataSource();
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, public dataService: DataService, public dateUtil: DateUtil) { }

  ngOnInit() {
    this.dataService.sequences
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(sequences => {
        this.updateDataSource();
        this._changeDetectorRef.markForCheck();
      });
  }

  updateDataSource() {
    this.dataSource.data = this.sequenceStates.slice(0, this.PAGE_SIZE) || [];
  }

  getServiceLink(sequence: Sequence) {
    return ['/project', sequence.project, 'service', sequence.service, 'context', sequence.shkeptncontext, 'stage', sequence.getLastStage()];
  }

  getSequenceLink(sequence: Sequence) {
    return ['/project', sequence.project, 'sequence', sequence.shkeptncontext, 'stage', sequence.getLastStage()];
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete()
  }

}
