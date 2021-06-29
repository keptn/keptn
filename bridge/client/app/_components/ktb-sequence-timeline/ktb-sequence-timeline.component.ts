import {ChangeDetectorRef, Component, Input, Output, EventEmitter} from '@angular/core';
import {Root} from '../../_models/root';
import {ActivatedRoute} from '@angular/router';
import {Subject} from 'rxjs';
import {Sequence} from '../../_models/sequence';

@Component({
  selector: 'ktb-sequence-timeline',
  templateUrl: './ktb-sequence-timeline.component.html',
  styleUrls: ['./ktb-sequence-timeline.component.scss']
})
export class KtbSequenceTimelineComponent{
  private _currentSequence: Sequence;
  public _selectedStage: String;

  @Output() selectedStageChange: EventEmitter<String> = new EventEmitter();

  @Input()
  get selectedStage(): String {
    return this._selectedStage;
  }
  set selectedStage(stage: String) {
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

  selectStage(stage: String) {
    if (this.selectedStage !== stage) {
      this.stageChanged(stage);
    }
  }

  stageChanged(stageName: String) {
    this.selectedStage = stageName;
    this._changeDetectorRef.markForCheck();
    this.selectedStageChange.emit(stageName);
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) {
  }
}
