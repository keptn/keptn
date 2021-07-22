import {Location} from '@angular/common';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component, Input,
  OnDestroy,
  TemplateRef,
  ViewChild
} from '@angular/core';
import {Deployment} from '../../_models/deployment';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute, Router} from '@angular/router';
import {takeUntil} from 'rxjs/operators';
import { forkJoin, Observable, Subject } from 'rxjs';
import {Trace} from '../../_models/trace';
import {MatDialog, MatDialogRef} from '@angular/material/dialog';
import {ClipboardService} from '../../_services/clipboard.service';

@Component({
  selector: 'ktb-service-details',
  templateUrl: './ktb-service-details.component.html',
  styleUrls: ['./ktb-service-details.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbServiceDetailsComponent implements OnDestroy{
  private _deployment?: Deployment;
  private readonly unsubscribe$: Subject<void> = new Subject<void>();
  private _selectedStage?: string;
  @ViewChild('remediationDialog')
  // tslint:disable-next-line:no-any
  public remediationDialog?: TemplateRef<any>;
  // tslint:disable-next-line:no-any
  public remediationDialogRef?: MatDialogRef<any, any>;
  public projectName?: string;

  @Input()
  get selectedStage(): string | undefined {
    return this._selectedStage;
  }
  set selectedStage(stageName: string | undefined) {
    this.selectStage(stageName);
  }

  @Input()
  get deployment(): Deployment | undefined {
    return this._deployment;
  }

  set deployment(deployment: Deployment | undefined) {
    if (deployment && this._deployment !== deployment) {
      if (!deployment.sequence) {
        this.loadSequence(deployment);
      } else {
        this._deployment = deployment;
      }
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute,
              private router: Router, private location: Location, private dialog: MatDialog, private clipboard: ClipboardService) {
    this.route.params.pipe(
      takeUntil(this.unsubscribe$)
    ).subscribe(params => {
      this.projectName = params.projectName;
      this._changeDetectorRef.markForCheck();
    });
  }

  private loadSequence(deployment: Deployment) {
    if (this.projectName) {
      this.dataService.getRoot(this.projectName, deployment.shkeptncontext).subscribe(sequence => {
        deployment.sequence = sequence;
        const evaluations$: Observable<Trace | undefined>[] = this.fetchEvaluations(deployment);
        if (evaluations$.length !== 0) {
          forkJoin(evaluations$)
            .subscribe((evaluations: (Trace | undefined)[]) => {
              for (const evaluation of evaluations) {
                this.setEvaluation(deployment, evaluation);
              }
              this._deployment = deployment;
              this._changeDetectorRef.markForCheck();
            });
        } else {
          this._deployment = deployment;
          this._changeDetectorRef.markForCheck();
        }
      });
    }
  }

  private fetchEvaluations(deployment: Deployment) {
    const evaluations$: Observable<Trace | undefined>[] = [];
    for (const stage of deployment.stages) {
      if (!stage.evaluation && stage.evaluationContext) {
        evaluations$.push(this.dataService.getEvaluationResult(stage.evaluationContext));
      }
    }
    return evaluations$;
  }

  private setEvaluation(deployment: Deployment, evaluation: Trace | undefined) {
    if (evaluation?.stage) {
      const stage = deployment.getStage(evaluation.stage);
      if (stage) {
        stage.evaluation = evaluation;
      }
    }
  }

  public selectStage(stageName: string | undefined) {
    this._selectedStage = stageName;
    if (this.deployment?.sequence) {
      const routeUrl = this.router.createUrlTree(['/project', this.projectName, 'service', this.deployment.service, 'context',
                                                            this.deployment.sequence.shkeptncontext, 'stage', stageName]);
      this.location.go(routeUrl.toString());
    }
    this._changeDetectorRef.markForCheck();
  }

  public isUrl(value: string): boolean {
    try {
      // tslint:disable-next-line:no-unused-expression
      new URL(value);
    } catch (_) {
      return false;
    }
    return true;
  }

  public showRemediationConfigDialog(): void {
    if (this.remediationDialog && this.deployment && this.selectedStage) {
      this.remediationDialogRef = this.dialog.open(this.remediationDialog, {data: this.deployment.getStage(this.selectedStage)?.config});
    }
  }

  public closeRemediationConfigDialog(): void {
    if (this.remediationDialogRef) {
      this.remediationDialogRef.close();
    }
  }

  public copyPayload(plainEvent: string): void {
    this.clipboard.copy(plainEvent, 'remediation payload');
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
