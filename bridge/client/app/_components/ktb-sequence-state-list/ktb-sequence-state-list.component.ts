import { ChangeDetectionStrategy, Component, Input, ViewEncapsulation } from '@angular/core';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { DateUtil } from '../../_utils/date.utils';
import { Sequence } from '../../_models/sequence';
import { Project } from '../../_models/project';
import { Router } from '@angular/router';

@Component({
  selector: 'ktb-sequence-state-list',
  templateUrl: './ktb-sequence-state-list.component.html',
  styleUrls: [],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSequenceStateListComponent {
  @Input() project?: Project;
  private _sequenceStates: Sequence[] = [];
  public dataSource: DtTableDataSource<Sequence> = new DtTableDataSource();
  public SequenceClass = Sequence;

  @Input()
  get sequenceStates(): Sequence[] {
    return this._sequenceStates;
  }

  set sequenceStates(value: Sequence[]) {
    if (this._sequenceStates !== value) {
      this._sequenceStates = value;
      this.updateDataSource();
    }
  }

  constructor(public dateUtil: DateUtil, private router: Router) {}

  updateDataSource(): void {
    this.dataSource = new DtTableDataSource(this.sequenceStates);
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
}
