import {Location} from '@angular/common';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component, Input,
  OnDestroy,
  OnInit,
  TemplateRef,
  ViewChild
} from '@angular/core';
import {Deployment} from '../../_models/deployment';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute, Router} from '@angular/router';
import {defaultIfEmpty, filter, takeUntil} from 'rxjs/operators';
import {forkJoin, Subject} from 'rxjs';
import {Trace} from '../../_models/trace';
import {MatDialog, MatDialogRef} from '@angular/material/dialog';
import {ClipboardService} from '../../_services/clipboard.service';

@Component({
  selector: 'ktb-service-details',
  templateUrl: './ktb-service-details.component.html',
  styleUrls: ['./ktb-service-details.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbServiceDetailsComponent implements OnInit, OnDestroy{
  private _deployment: Deployment;
  private readonly unsubscribe$: Subject<void> = new Subject<void>();
  @ViewChild('remediationDialog')
  public remediationDialog: TemplateRef<any>;
  public remediationDialogRef: MatDialogRef<any, any>;
  private _selectedStage: string;

  public projectName: string;

  @Input()
  get selectedStage(): string {
    return this._selectedStage;
  }
  set selectedStage(stageName: string) {
    this.selectStage(stageName);
  }

  @Input()
  get deployment(): Deployment {
    return this._deployment;
  }

  set deployment(deployment: Deployment) {
    if (this._deployment !== deployment) {
      if (deployment) {
        if (!deployment.sequence) {
          this.loadSequence(deployment);
        } else {
          this._deployment = deployment;
        }
      }
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute, private router: Router, private location: Location, private dialog: MatDialog, private clipboard: ClipboardService) {
    this.route.params.pipe(
      takeUntil(this.unsubscribe$)
    ).subscribe(params => {
      this.projectName = params.projectName;
      this._changeDetectorRef.markForCheck();
    });
  }

  ngOnInit(): void {

  }

  private loadSequence(deployment: Deployment) {
    if (deployment) {
      this.dataService.getRoot(this.projectName, deployment.shkeptncontext).subscribe(sequence => {
        deployment.sequence = sequence;
        const evaluations$ = [];
        for (const stage of deployment.stages) {
          if (!stage.evaluation && stage.evaluationContext) {
            evaluations$.push(this.dataService.getEvaluationResult(stage.evaluationContext));
          }
        }
        forkJoin(evaluations$)
          .pipe(defaultIfEmpty(null))
          .subscribe((evaluations: Trace[] | null) => {
            if (evaluations) {
              for (const evaluation of evaluations){
                deployment.getStage(evaluation.getStage()).evaluation = evaluation;
              }
            }
            this._deployment = deployment;
            this._changeDetectorRef.markForCheck();
        });
      });
    }
    else {
      this._deployment = deployment;
    }
  }

  public selectStage(stageName: string) {
    this._selectedStage = stageName;
    if (this.deployment?.sequence) {
      const routeUrl = this.router.createUrlTree(['/project', this.projectName, 'service', this.deployment.service, 'context', this.deployment.sequence.shkeptncontext, 'stage', stageName]);
      this.location.go(routeUrl.toString());
    }
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

  public showRemediationConfigDialog(): void {
    this.remediationDialogRef = this.dialog.open(this.remediationDialog, {data: this.deployment.getStage(this.selectedStage).config});
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
