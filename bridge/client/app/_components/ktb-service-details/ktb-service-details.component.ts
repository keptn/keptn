import {Location} from '@angular/common';
import {ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {Deployment} from '../../_models/deployment';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute, Router} from '@angular/router';
import {defaultIfEmpty, filter, takeUntil} from 'rxjs/operators';
import {forkJoin, Subject} from 'rxjs';
import {Trace} from '../../_models/trace';

@Component({
  selector: 'ktb-service-details',
  templateUrl: './ktb-service-details.component.html',
  styleUrls: ['./ktb-service-details.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbServiceDetailsComponent implements OnInit, OnDestroy{
  private _deployment: Deployment;
  private readonly unsubscribe$: Subject<void> = new Subject<void>();

  public projectName: string;
  public selectedStage: string;

  get deployment(): Deployment {
    return this._deployment;
  }

  set deployment(deployment: Deployment) {
    if (this._deployment !== deployment) {
      const selectLast = !!this._deployment;
      this._deployment = deployment;
      if (!this._deployment.sequence) {
        this.loadSequence(selectLast);
      } else {
        this.selectLastStage();
      }
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute, private router: Router, private location: Location) {

  }

  ngOnInit(): void {
    this.route.params.pipe(
      takeUntil(this.unsubscribe$)
    ).subscribe(params => {
      this.projectName = params.projectName;
      this.selectedStage = params.stage;
      this._changeDetectorRef.markForCheck();
    });
  }

  private loadSequence(selectLast: boolean) {
    if (this.deployment) {
      this.dataService.getRoot(this.projectName, this.deployment.shkeptncontext).subscribe(sequence => {
        this.deployment.sequence = sequence;
        const evaluations$ = [];
        for (const stage of this.deployment.stages) {
          if (!stage.evaluation && stage.evaluationContext) {
            evaluations$.push(this.dataService.getEvaluationResult(stage.evaluationContext));
          }
        }
        forkJoin(evaluations$)
          .pipe(defaultIfEmpty(null))
          .subscribe((evaluations: Trace[] | null) => {
            if (evaluations) {
              for (const evaluation of evaluations){
                this.deployment.getStage(evaluation.getStage()).evaluation = evaluation;
              }
            }

            if (selectLast || !this.selectedStage) {
              this.selectLastStage();
            }
            this._changeDetectorRef.markForCheck();
        });
      });
    }
  }

  private selectLastStage() {
    const stages = this.deployment.sequence.getStages();
    this.selectStage(stages[stages.length - 1]);
  }

  public selectStage(stageName: string) {
    this.selectedStage = stageName;
    const routeUrl = this.router.createUrlTree(['/project', this.projectName, 'service', this.deployment.service, 'context', this.deployment.sequence.shkeptncontext, 'stage', stageName]);
    this.location.go(routeUrl.toString());
    this._changeDetectorRef.markForCheck();
  }

  public isUrl(value: string): boolean {
    try {
      new URL(value);
    } catch (_) {
      return false;
    }
    return true;
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
