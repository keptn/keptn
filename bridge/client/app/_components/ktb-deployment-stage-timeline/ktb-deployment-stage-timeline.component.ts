import { Component, EventEmitter, Input, OnDestroy, Output } from '@angular/core';
import { Deployment } from '../../_models/deployment';
import { Subject } from 'rxjs';

@Component({
  selector: 'ktb-deployment-timeline[deployment]',
  templateUrl: './ktb-deployment-stage-timeline.component.html',
  styleUrls: ['./ktb-deployment-stage-timeline.component.scss'],
})
export class KtbDeploymentStageTimelineComponent implements OnDestroy {
  @Input() deployment?: Deployment;
  @Input() selectedStage?: string;
  @Output() selectedStageChange: EventEmitter<string> = new EventEmitter<string>();
  private unsubscribe$: Subject<void> = new Subject();

  public selectStage(stage: string): void {
    if (this.selectedStage !== stage) {
      this.selectedStage = stage;
      this.selectedStageChange.emit(this.selectedStage);
    }
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
