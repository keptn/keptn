import {ChangeDetectorRef, Component, Input, OnInit, Output, EventEmitter} from '@angular/core';
import {Root} from "../../_models/root";

@Component({
  selector: 'ktb-sequence-timeline',
  templateUrl: './ktb-sequence-timeline.component.html',
  styleUrls: ['./ktb-sequence-timeline.component.scss']
})
export class KtbSequenceTimelineComponent implements OnInit {
  private _currentRoot: Root;

  @Output() selectedStageChange: EventEmitter<String> = new EventEmitter();

  @Input()
  get currentRoot(): Root {
    return this._currentRoot;
  }
  set currentRoot(root: Root) {
    if (this._currentRoot !== root) {
      this._currentRoot = root;
      const stages = this._currentRoot.getStages();
      this.selectStage(stages[stages.length-1]);
      this._changeDetectorRef.markForCheck();
    }
  }


  selectStage(stage: String) {
    this.selectedStageChange.emit(stage);
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit(): void {
  }

}
