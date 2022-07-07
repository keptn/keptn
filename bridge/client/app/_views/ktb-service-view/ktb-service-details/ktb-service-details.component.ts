import { Location } from '@angular/common';
import { Component, Input, OnDestroy, TemplateRef, ViewChild } from '@angular/core';
import { DataService } from '../../../_services/data.service';
import { ActivatedRoute, Router } from '@angular/router';
import { takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { ClipboardService } from '../../../_services/clipboard.service';
import { DeploymentSelection } from '../../../_interfaces/deployment-selection';

@Component({
  selector: 'ktb-service-details',
  templateUrl: './ktb-service-details.component.html',
  styleUrls: ['./ktb-service-details.component.scss'],
})
export class KtbServiceDetailsComponent implements OnDestroy {
  private _deploymentInfo?: DeploymentSelection;
  private readonly unsubscribe$: Subject<void> = new Subject<void>();
  @ViewChild('remediationDialog')
  /* eslint-disable @typescript-eslint/no-explicit-any */
  public remediationDialog?: TemplateRef<any>;
  public remediationDialogRef?: MatDialogRef<any, any>;
  /* eslint-enable @typescript-eslint/no-explicit-any */
  public projectName?: string;
  public isLoading = false;

  @Input()
  get deploymentInfo(): DeploymentSelection | undefined {
    return this._deploymentInfo;
  }

  set deploymentInfo(info: DeploymentSelection | undefined) {
    this._deploymentInfo = info;
  }

  constructor(
    private dataService: DataService,
    private route: ActivatedRoute,
    private router: Router,
    private location: Location,
    private dialog: MatDialog,
    private clipboard: ClipboardService
  ) {
    this.route.paramMap.pipe(takeUntil(this.unsubscribe$)).subscribe((params) => {
      this.projectName = params.get('projectName') ?? undefined;
    });
  }

  public selectStage(stageName: string): void {
    if (this.deploymentInfo) {
      this.deploymentInfo.stage = stageName;
      const routeUrl = this.router.createUrlTree([
        '/project',
        this.projectName,
        'service',
        this.deploymentInfo.deployment.service,
        'context',
        this.deploymentInfo.deployment.keptnContext,
        'stage',
        stageName,
      ]);
      this.location.go(routeUrl.toString());
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

  public showRemediationConfigDialog(config: string): void {
    if (this.remediationDialog) {
      this.remediationDialogRef = this.dialog.open(this.remediationDialog, {
        data: atob(config),
      });
    }
  }

  public closeRemediationConfigDialog(): void {
    this.remediationDialogRef?.close();
  }

  public copyPayload(plainEvent: string): void {
    this.clipboard.copy(plainEvent, 'remediation payload');
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
