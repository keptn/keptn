import {
  ChangeDetectorRef,
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
import { Observable, of, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { Trace } from '../../_models/trace';
import { Project } from '../../_models/project';
import { ClipboardService } from '../../_services/clipboard.service';
import { DataService } from '../../_services/data.service';
import { DateUtil } from '../../_utils/date.utils';
import { ActivatedRoute } from '@angular/router';

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
  public project$: Observable<Project | undefined> = of(undefined);
  public _task?: Trace;
  @Input() public isExpanded = false;

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

  set task(value: Trace | undefined) {
    if (this._task !== value) {
      this._task = value;
      this.changeDetectorRef.markForCheck();
    }
  }

  constructor(
    private changeDetectorRef: ChangeDetectorRef,
    private dataService: DataService,
    private dialog: MatDialog,
    private clipboard: ClipboardService,
    public dateUtil: DateUtil,
    private route: ActivatedRoute
  ) {}

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

  isUrl(value: string): boolean {
    try {
      // tslint:disable-next-line:no-unused-expression
      new URL(value);
    } catch (_) {
      return false;
    }
    return true;
  }

  onClick(item: Trace): void {
    this.itemClicked.emit(item);
  }

  ngOnInit(): void {
    if (this.task?.project) {
      this.project$ = this.dataService.getProject(this.task?.project);
    }
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
