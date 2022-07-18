import {
  Component,
  Directive,
  EventEmitter,
  Input,
  OnDestroy,
  OnInit,
  Output,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { Trace } from '../../../_models/trace';
import { ClipboardService } from '../../../_services/clipboard.service';
import { DataService } from '../../../_services/data.service';
import { DateUtil } from '../../../_utils/date.utils';
import { ActivatedRoute } from '@angular/router';
import { AppUtils } from '../../../_utils/app.utils';

@Directive({
  selector: `ktb-task-item-detail, [ktb-task-item-detail], [ktbTaskItemDetail]`,
  exportAs: 'ktbTaskItemDetail',
})
export class KtbTaskItemDetailDirective {}

@Component({
  selector: 'ktb-task-item[task]',
  templateUrl: './ktb-task-item.component.html',
  styleUrls: ['./ktb-task-item.component.scss'],
})
export class KtbTaskItemComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  public _task?: Trace;
  public latestDeployment?: string;
  public isUrl = AppUtils.isValidUrl;
  private _expanded = false;
  @Input()
  public set isExpanded(isExpanded: boolean) {
    this._expanded = isExpanded;
    this.setLatestDeployment();
  }
  public get isExpanded(): boolean {
    return this._expanded;
  }
  @Input() public isSubtask = false;

  @ViewChild('taskPayloadDialog')
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public taskPayloadDialog?: TemplateRef<any>;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public taskPayloadDialogRef?: MatDialogRef<any, any>;

  @Output() itemClicked: EventEmitter<Trace> = new EventEmitter();

  @Input()
  get task(): Trace | undefined {
    return this._task;
  }

  set task(trace: Trace | undefined) {
    if (this._task !== trace) {
      this._task = trace;
      this.latestDeployment = undefined;
      this.setLatestDeployment();
    }
  }

  constructor(
    private dataService: DataService,
    private dialog: MatDialog,
    private clipboard: ClipboardService,
    public dateUtil: DateUtil,
    private route: ActivatedRoute
  ) {}

  private setLatestDeployment(): void {
    if (!this.isExpanded || !this.task || this.latestDeployment !== undefined || !this.task.isApproval()) {
      return;
    }

    const { project, stage, service } = this.task;
    if (!project || !stage || !service) {
      this.latestDeployment = '';
      return;
    }

    this.dataService.getService(project, stage, service).subscribe((svc) => {
      this.latestDeployment = svc.deployedImage ?? '';
    });
  }

  showEventPayloadDialog(event: MouseEvent, task: Trace): void {
    event.stopPropagation();
    if (this.taskPayloadDialog) {
      this.taskPayloadDialogRef = this.dialog.open(this.taskPayloadDialog, { data: task.plainEvent });
    }
  }

  closeEventPayloadDialog(): void {
    if (this.taskPayloadDialogRef) {
      this.taskPayloadDialogRef.close();
    }
  }

  copyEventPayload(plainEvent: string): void {
    this.clipboard.copy(plainEvent, 'payload');
  }

  onClick(item: Trace): void {
    this.itemClicked.emit(item);
  }

  ngOnInit(): void {
    this.route.params.pipe(takeUntil(this.unsubscribe$)).subscribe((params) => {
      if (params.eventId === this.task?.id) {
        this.isExpanded = true;
      }
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
