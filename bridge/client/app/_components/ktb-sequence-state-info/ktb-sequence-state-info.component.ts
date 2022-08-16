import { Component, EventEmitter, HostBinding, Input, Output, ViewEncapsulation } from '@angular/core';
import { ISequenceState } from '../../../../shared/interfaces/sequence';
import { createSequenceStateInfo, getLastStageName, getStageNames } from '../../_models/sequenceState';
import { EvaluationBadgeVariant } from '../ktb-evaluation-badge/ktb-evaluation-badge.utils';

@Component({
  selector: 'ktb-sequence-state-info',
  templateUrl: './ktb-sequence-state-info.component.html',
  styleUrls: ['./ktb-sequence-state-info.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class KtbSequenceStateInfoComponent {
  @HostBinding('class') cls = 'ktb-sequence-state-info';
  private _sequence?: ISequenceState;
  private _showStages = true;
  public EvaluationBadgeVariant = EvaluationBadgeVariant;
  createSequenceStateInfo = createSequenceStateInfo;
  getStageNames = getStageNames;

  @Output() readonly stageClicked = new EventEmitter<{ sequence: ISequenceState; stage?: string }>();

  @Input()
  get sequence(): ISequenceState | undefined {
    return this._sequence;
  }

  set sequence(sequence: ISequenceState | undefined) {
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

  getServiceLink(sequence: ISequenceState): (string | undefined)[] {
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

  getSequenceLink(sequence: ISequenceState): (string | undefined)[] {
    return ['/project', sequence.project, 'sequence', sequence.shkeptncontext, 'stage', getLastStageName(sequence)];
  }

  stageClick(sequence: ISequenceState, stage: string): void {
    this.stageClicked.emit({ sequence, stage });
  }
}
