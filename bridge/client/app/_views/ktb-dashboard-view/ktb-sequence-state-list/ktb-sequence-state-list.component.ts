import { ChangeDetectionStrategy, Component, Input, ViewEncapsulation } from '@angular/core';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { DateUtil } from '../../../_utils/date.utils';
import { Router } from '@angular/router';
import { ISequenceState } from '../../../../../shared/interfaces/sequence';
import { getLastStageName } from '../../../_models/sequenceState';

@Component({
  selector: 'ktb-sequence-state-list',
  templateUrl: './ktb-sequence-state-list.component.html',
  styleUrls: [],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSequenceStateListComponent {
  private _sequenceStates: ISequenceState[] = [];
  public dataSource: DtTableDataSource<ISequenceState> = new DtTableDataSource();

  @Input()
  get sequenceStates(): ISequenceState[] {
    return this._sequenceStates;
  }

  set sequenceStates(value: ISequenceState[]) {
    if (this._sequenceStates !== value) {
      this._sequenceStates = value;
      this.updateDataSource();
    }
  }

  constructor(public dateUtil: DateUtil, private router: Router) {}

  updateDataSource(): void {
    this.dataSource = new DtTableDataSource(this.sequenceStates);
  }

  selectSequence(event: { sequence: ISequenceState; stage?: string }): void {
    const stage = event.stage || getLastStageName(event.sequence);
    this.router.navigate([
      '/project',
      event.sequence.project,
      'sequence',
      event.sequence.shkeptncontext,
      ...(stage ? ['stage', stage] : []),
    ]);
  }

  toSequence(value: unknown): ISequenceState {
    return value as ISequenceState;
  }
}
