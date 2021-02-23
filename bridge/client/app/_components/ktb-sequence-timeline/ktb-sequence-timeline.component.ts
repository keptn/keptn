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
  public selectedStage: String;

  @Output() selectedStageChange: EventEmitter<{ stageName: String, triggerByEvent: boolean }> = new EventEmitter();

  @Input()
  get currentSequence(): Root {
    return this._currentSequence;
  }
  set currentSequence(root: Root) {
    if (this._currentSequence !== root) {
      if (!this._currentSequence) {
        let stage = this.route.snapshot.params.stage;
        let triggerByEvent = false;
        if (this.route.snapshot.params.eventId) {
          triggerByEvent = true;
          const trace = root.traces.find(t => t.id === this.route.snapshot.params.eventId);
          if (trace) {
            stage = trace.getStage();
          }
        }
        this.setSequence(root, stage, triggerByEvent);
      } else {
        this.setSequence(root);
      }
    }
  }

  setSequence(root: Root, stage?: string, triggerByEvent = false) {
    this._currentSequence = root;
    const stages = this._currentSequence.getStages();
    if (stage && stages.includes(stage)) {
      this.stageChanged(stage, triggerByEvent);
    }
    else {
      this.stageChanged(stages[stages.length - 1]);
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

  constructor(private _changeDetectorRef: ChangeDetectorRef, private route: ActivatedRoute) {
  }
}
