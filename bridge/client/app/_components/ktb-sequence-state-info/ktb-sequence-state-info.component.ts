import { Component, EventEmitter, HostBinding, Input, Output, ViewEncapsulation } from '@angular/core';
import { Sequence } from '../../_models/sequence';

@Component({
  selector: 'ktb-sequence-state-info',
  templateUrl: './ktb-sequence-state-info.component.html',
  styleUrls: ['./ktb-sequence-state-info.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class KtbSequenceStateInfoComponent {
  @HostBinding('class') cls = 'ktb-sequence-state-info';
  private _sequence?: Sequence;
  private _showStages = true;

  @Output() readonly stageClicked = new EventEmitter<{ sequence: Sequence; stage?: string }>();

  @Input()
  get sequence(): Sequence | undefined {
    return this._sequence;
  }
  set sequence(sequence: Sequence | undefined) {
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

  getServiceLink(sequence: Sequence): (string | undefined)[] {
    return [
      '/project',
      sequence.project,
      'service',
      sequence.service,
      'context',
      sequence.shkeptncontext,
      'stage',
      sequence.getLastStage(),
    ];
  }

  getSequenceLink(sequence: Sequence): (string | undefined)[] {
    return ['/project', sequence.project, 'sequence', sequence.shkeptncontext, 'stage', sequence.getLastStage()];
  }

  stageClick(sequence: Sequence, stage: string): void {
    this.stageClicked.emit({ sequence, stage });
  }
}
