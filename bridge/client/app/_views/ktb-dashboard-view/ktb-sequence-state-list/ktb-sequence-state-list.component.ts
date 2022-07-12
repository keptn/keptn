import { ChangeDetectionStrategy, Component, Input, ViewEncapsulation } from '@angular/core';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { DateUtil } from '../../../_utils/date.utils';
import { Router } from '@angular/router';
import { ISequence } from '../../../../../shared/interfaces/sequence';
import { getLastStageName } from '../../../_models/sequence';

@Component({
  selector: 'ktb-sequence-state-list',
  templateUrl: './ktb-sequence-state-list.component.html',
  styleUrls: [],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSequenceStateListComponent {
  private _sequenceStates: ISequence[] = [];
  public dataSource: DtTableDataSource<ISequence> = new DtTableDataSource();

  @Input()
  get sequenceStates(): ISequence[] {
    return this._sequenceStates;
  }

  set sequenceStates(value: ISequence[]) {
    if (this._sequenceStates !== value) {
      this._sequenceStates = value;
      this.updateDataSource();
    }
  }

  constructor(public dateUtil: DateUtil, private router: Router) {}

  updateDataSource(): void {
    this.dataSource = new DtTableDataSource(this.sequenceStates);
  }

  selectSequence(event: { sequence: ISequence; stage?: string }): void {
    const stage = event.stage || getLastStageName(event.sequence);
    this.router.navigate([
      '/project',
      event.sequence.project,
      'sequence',
      event.sequence.shkeptncontext,
      ...(stage ? ['stage', stage] : []),
    ]);
  }

  toSequence(value: unknown): ISequence {
    return value as ISequence;
  }
}
