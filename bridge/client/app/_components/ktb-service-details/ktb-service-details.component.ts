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

type DeploymentInfo = { deployment: Deployment, stage: string };
@Component({
  selector: 'ktb-service-details',
  templateUrl: './ktb-service-details.component.html',
  styleUrls: ['./ktb-service-details.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbServiceDetailsComponent implements OnDestroy{
  private _deploymentInfo?: DeploymentInfo;
  private readonly unsubscribe$: Subject<void> = new Subject<void>();
  @ViewChild('remediationDialog')
  // tslint:disable-next-line:no-any
  public remediationDialog?: TemplateRef<any>;
  // tslint:disable-next-line:no-any
  public remediationDialogRef?: MatDialogRef<any, any>;
  public projectName?: string;
  public isLoading = false;

  @Input()
  get deploymentInfo(): DeploymentInfo | undefined {
    return this._deploymentInfo;
  }

  set deploymentInfo(info: DeploymentInfo | undefined) {
    if (info && this._deploymentInfo !== info) {
      if (this.deploymentInfo?.deployment.shkeptncontext !== info.deployment.shkeptncontext) {
        this._deploymentInfo = undefined;
        if (!info.deployment.sequence) {
          this.isLoading = true;
        }
      }
      if (!info.deployment.sequence) {
        this.loadSequence(info);
      } else {
        this._deploymentInfo = info;
      }
      this._changeDetectorRef.markForCheck();
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

  private loadSequence(info: DeploymentInfo) {
    if (this.projectName) {
      this.dataService.getRoot(this.projectName, info.deployment.shkeptncontext).subscribe(sequence => {
        info.deployment.sequence = sequence;
        const evaluations$: Observable<Trace | undefined>[] = this.fetchEvaluations(info.deployment);
        if (evaluations$.length !== 0) {
          forkJoin(evaluations$)
            .subscribe((evaluations: (Trace | undefined)[]) => {
              for (const evaluation of evaluations) {
                info.deployment.setEvaluation(evaluation);
              }
              this.isLoading = false;
              this._deploymentInfo = info;
              this._changeDetectorRef.markForCheck();
            });
        } else {
          this.isLoading = false;
          this._deploymentInfo = info;
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

  public selectStage(stageName: string) {
    if (this.deploymentInfo) {
      this.deploymentInfo.stage = stageName;
      const routeUrl = this.router.createUrlTree(['/project', this.projectName, 'service', this.deploymentInfo.deployment.service, 'context',
        this.deploymentInfo.deployment.shkeptncontext, 'stage', stageName]);
      this.location.go(routeUrl.toString());
      this._changeDetectorRef.markForCheck();
    }
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
    if (this.remediationDialog && this.deploymentInfo) {
      this.remediationDialogRef = this.dialog.open(this.remediationDialog, {data: this.deploymentInfo.deployment.getStage(this.deploymentInfo.stage)?.config});
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
