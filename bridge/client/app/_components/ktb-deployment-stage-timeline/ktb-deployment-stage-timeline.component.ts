import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Deployment } from '../../_models/deployment';

@Component({
  selector: 'ktb-deployment-timeline[deployment]',
  templateUrl: './ktb-deployment-stage-timeline.component.html',
  styleUrls: ['./ktb-deployment-stage-timeline.component.scss'],
})
export class KtbDeploymentStageTimelineComponent {
  @Input() deployment?: Deployment;
  @Input() selectedStage?: string;
  @Output() selectedStageChange: EventEmitter<string> = new EventEmitter<string>();

  public selectStage(stage: string): void {
    if (this.selectedStage !== stage) {
      this.selectedStage = stage;
      this.selectedStageChange.emit(this.selectedStage);
    }
  }
}
