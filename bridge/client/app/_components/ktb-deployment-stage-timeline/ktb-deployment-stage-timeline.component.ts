import {Component, EventEmitter, Input, OnDestroy, OnInit, Output} from '@angular/core';
import {Deployment} from '../../_models/deployment';
import {takeUntil} from 'rxjs/operators';
import {ActivatedRoute} from '@angular/router';
import {Subject} from 'rxjs';

@Component({
  selector: 'ktb-deployment-timeline',
  templateUrl: './ktb-deployment-stage-timeline.component.html',
  styleUrls: ['./ktb-deployment-stage-timeline.component.scss']
})
export class KtbDeploymentStageTimelineComponent implements OnInit, OnDestroy {
  @Input() deployment: Deployment;
  @Input() selectedStage: string;
  @Output() selectedStageChange: EventEmitter<string> = new EventEmitter<string>();
  private readonly unsubscribe$ = new Subject<void>();
  public projectName: string;

  constructor(private route: ActivatedRoute) { }

  public ngOnInit(): void {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        this.projectName = params.projectName;
      });
  }

  public selectStage(stage: string): void {
    if (this.selectedStage !== stage) {
      this.selectedStage = stage;
      this.selectedStageChange.emit(this.selectedStage);
    }
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
