import {ChangeDetectorRef, Component, Input, Output, EventEmitter} from '@angular/core';
import {Sequence} from '../../_models/sequence';

@Component({
  selector: 'ktb-sequence-timeline',
  templateUrl: './ktb-sequence-timeline.component.html',
  styleUrls: ['./ktb-sequence-timeline.component.scss']
})
export class KtbSequenceTimelineComponent{
  private _currentSequence: Sequence;
  public _selectedStage: string;

  @Output() selectedStageChange: EventEmitter<string> = new EventEmitter();

  @Input()
  get selectedStage(): string {
    return this._selectedStage;
  }
  set selectedStage(stage: string) {
    if(this._selectedStage !== stage) {
      this._selectedStage = stage;
    }
  }

  @Input()
  get currentSequence(): Sequence {
    return this._currentSequence;
  }
  set currentSequence(sequence: Sequence) {
    if (this._currentSequence !== sequence) {
      this._currentSequence = sequence;
    }
  }

  selectStage(stage: string) {
    if (this.selectedStage !== stage) {
      this.stageChanged(stage);
    }
  }

  stageChanged(stageName: string) {
    this.selectedStage = stageName;
    this._changeDetectorRef.markForCheck();
    this.selectedStageChange.emit(stageName);
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) {
  }
}
