import { Component, EventEmitter, HostBinding, Input, Output, ViewEncapsulation } from '@angular/core';
import { ISequence } from '../../../../shared/interfaces/sequence';
import { createSequenceStateInfo, getLastStageName, getStageNames } from '../../_models/sequence';

@Component({
  selector: 'ktb-sequence-state-info',
  templateUrl: './ktb-sequence-state-info.component.html',
  styleUrls: ['./ktb-sequence-state-info.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class KtbSequenceStateInfoComponent {
  @HostBinding('class') cls = 'ktb-sequence-state-info';
  private _sequence?: ISequence;
  private _showStages = true;

  createSequenceStateInfo = createSequenceStateInfo;
  getStageNames = getStageNames;

  @Output() readonly stageClicked = new EventEmitter<{ sequence: ISequence; stage?: string }>();

  @Input()
  get sequence(): ISequence | undefined {
    return this._sequence;
  }

  set sequence(sequence: ISequence | undefined) {
    if (this._sequence !== sequence) {
      this._sequence = sequence;
    }
  }

  @Input()
  get showStages(): boolean {
    return this._showStages;
  }

  set showStages(showStages: boolean) {
    if (this._showStages !== showStages) {
      this._showStages = showStages;
    }
  }

  getServiceLink(sequence: ISequence): (string | undefined)[] {
    return [
      '/project',
      sequence.project,
      'service',
      sequence.service,
      'context',
      sequence.shkeptncontext,
      'stage',
      getLastStageName(sequence),
    ];
  }

  getSequenceLink(sequence: ISequence): (string | undefined)[] {
    return ['/project', sequence.project, 'sequence', sequence.shkeptncontext, 'stage', getLastStageName(sequence)];
  }

  stageClick(sequence: ISequence, stage: string): void {
    this.stageClicked.emit({ sequence, stage });
  }
}
