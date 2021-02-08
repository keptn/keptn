import {ChangeDetectorRef, Component, Input, OnInit, Output, EventEmitter} from '@angular/core';
import {Root} from '../../_models/root';

@Component({
  selector: 'ktb-sequence-timeline',
  templateUrl: './ktb-sequence-timeline.component.html',
  styleUrls: ['./ktb-sequence-timeline.component.scss']
})
export class KtbSequenceTimelineComponent implements OnInit {
  private _currentSequence: Root;
  public selectedStage: String;

  @Output() selectedStageChange: EventEmitter<String> = new EventEmitter();

  @Input()
  get currentSequence(): Root {
    return this._currentSequence;
  }
  set currentSequence(root: Root) {
    if (this._currentSequence !== root) {
      this._currentSequence = root;
      const stages = this._currentSequence.getStages();
      this.stageChanged(stages[stages.length - 1]);
    }
  }

  selectStage(stage: String) {
    if (this.selectedStage !== stage) {
      this.stageChanged(stage);
    }
  }

  stageChanged(stage: String) {
    this.selectedStage = stage;
    this._changeDetectorRef.markForCheck();
    this.selectedStageChange.emit(stage);
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit(): void {
  }

}
