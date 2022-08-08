import { Component, EventEmitter, Input, Output } from '@angular/core';
import { EvaluationBadgeVariant } from '../../../_components/ktb-evaluation-badge/ktb-evaluation-badge.utils';
import { createStageDeploymentStateInfo, Deployment } from '../../../_models/deployment';

@Component({
  selector: 'ktb-deployment-timeline[deployment]',
  templateUrl: './ktb-deployment-stage-timeline.component.html',
  styleUrls: ['./ktb-deployment-stage-timeline.component.scss'],
})
export class KtbDeploymentStageTimelineComponent {
  public createStageDeploymentStateInfo = createStageDeploymentStateInfo;
  public EvaluationBadgeFillState = EvaluationBadgeVariant;

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
