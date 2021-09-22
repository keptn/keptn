import { Component, EventEmitter, Input, Output, ViewEncapsulation } from '@angular/core';
import {Sequence} from '../../_models/sequence';

@Component({
  selector: 'ktb-sequence-state-info',
  templateUrl: './ktb-sequence-state-info.component.html',
  styleUrls: ['./ktb-sequence-state-info.component.scss'],
  host: {
    class: 'ktb-sequence-state-info',
  },
  encapsulation: ViewEncapsulation.None,
})
export class KtbSequenceStateInfoComponent {

  private _sequence?: Sequence;
  private _showOnlyLastStage = false;

  @Output() readonly stageClicked = new EventEmitter<{ sequence: Sequence, stage?: string }>();

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
  get showOnlyLastStage(): boolean {
    return this._showOnlyLastStage;
  }
  set showOnlyLastStage(showOnlyLastStage: boolean) {
    if (this._showOnlyLastStage !== showOnlyLastStage) {
      this._showOnlyLastStage = showOnlyLastStage;
    }
  }

  constructor() {
  }

  getStages(): (string | undefined)[] | undefined {
    return this.showOnlyLastStage ? [this.sequence?.getLastStage()] : this.sequence?.getStages();
  }

  getServiceLink(sequence: Sequence): (string | undefined)[] {
    return ['/project', sequence.project, 'service', sequence.service, 'context', sequence.shkeptncontext, 'stage', sequence.getLastStage()];
  }

  getSequenceLink(sequence: Sequence): (string | undefined)[] {
    return ['/project', sequence.project, 'sequence', sequence.shkeptncontext, 'stage', sequence.getLastStage()];
  }

  stageClick(sequence: Sequence, stage: string): void {
    this.stageClicked.emit({sequence, stage});
  }
}
