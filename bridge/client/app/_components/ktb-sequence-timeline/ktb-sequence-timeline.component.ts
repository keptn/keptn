import {ChangeDetectorRef, Component, Input, Output, EventEmitter} from '@angular/core';
import {Root} from '../../_models/root';

@Component({
  selector: 'ktb-sequence-timeline',
  templateUrl: './ktb-sequence-timeline.component.html',
  styleUrls: ['./ktb-sequence-timeline.component.scss']
})
export class KtbSequenceTimelineComponent{
  private _currentSequence: Root;
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
  get currentSequence(): Root {
    return this._currentSequence;
  }
  set currentSequence(root: Root) {
    if (this._currentSequence !== root) {
      this._currentSequence = root;
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
