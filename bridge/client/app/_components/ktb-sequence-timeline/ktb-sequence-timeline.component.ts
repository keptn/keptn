import {ChangeDetectorRef, Component, Input, Output, EventEmitter} from '@angular/core';
import {Root} from '../../_models/root';
import {ActivatedRoute} from '@angular/router';
import {Subject} from 'rxjs';

@Component({
  selector: 'ktb-sequence-timeline',
  templateUrl: './ktb-sequence-timeline.component.html',
  styleUrls: ['./ktb-sequence-timeline.component.scss']
})
export class KtbSequenceTimelineComponent{
  private _currentSequence: Root;
  public _selectedStage: String;

  @Output() selectedStageChange: EventEmitter<{ stageName: String, triggerByEvent: boolean }> = new EventEmitter();

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
  get currentSequence(): Root {
    return this._currentSequence;
  }
  set currentSequence(root: Root) {
    if (this._currentSequence !== root) {
      this._currentSequence = root;
    }
  }

  selectStage(stage: String) {
    if (this.selectedStage !== stage) {
      this.stageChanged(stage);
    }
  }

  stageChanged(stageName: String, triggerByEvent = false) {
    this.selectedStage = stageName;
    this._changeDetectorRef.markForCheck();
    this.selectedStageChange.emit({stageName, triggerByEvent});
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) {
  }
}
