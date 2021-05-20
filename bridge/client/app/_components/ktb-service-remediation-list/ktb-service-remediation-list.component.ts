import {Component, Input, OnDestroy, OnInit, TemplateRef, ViewChild} from '@angular/core';
import {Sequence} from '../../_models/sequence';
import {MatDialog, MatDialogRef} from '@angular/material/dialog';
import {ClipboardService} from '../../_services/clipboard.service';
import {DateUtil} from '../../_utils/date.utils';
import {takeUntil} from 'rxjs/operators';
import {ActivatedRoute} from '@angular/router';
import {Subject} from 'rxjs';

@Component({
  selector: 'ktb-service-remediation-list',
  templateUrl: './ktb-service-remediation-list.component.html',
  styleUrls: ['./ktb-service-remediation-list.component.scss']
})
export class KtbServiceRemediationListComponent implements OnInit, OnDestroy {
  @Input() stage: {stageName: string, remediations: Sequence[], config: string };

  @ViewChild('remediationDialog')
  public remediationDialog: TemplateRef<any>;
  public remediationDialogRef: MatDialogRef<any, any>;
  private unsubscribe$: Subject<void> = new Subject<void>();
  public projectName: string;

  constructor(private dialog: MatDialog, private clipboard: ClipboardService, public dateUtil: DateUtil, private route: ActivatedRoute) { }

  ngOnInit(): void {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        this.projectName = params.projectName;
      });
  }

  public showDialog(): void {
    this.remediationDialogRef = this.dialog.open(this.remediationDialog, {data: this.stage.config});
  }

  public closeDialog(): void {
    if (this.remediationDialogRef) {
      this.remediationDialogRef.close();
    }
  }

  public copyPayload(plainEvent: string): void {
    this.clipboard.copy(plainEvent, 'remediation payload');
  }

  public getRemediationLink(remediation: Sequence) {
    return ['/', 'project', this.projectName, 'sequence', remediation.shkeptncontext, 'event', remediation.getStage(this.stage.stageName).latestEvent.id];
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
