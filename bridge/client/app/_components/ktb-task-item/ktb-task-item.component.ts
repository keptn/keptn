import { ChangeDetectorRef, Component, Directive, Input, TemplateRef, ViewChild } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import {Observable} from "rxjs";
import {map} from "rxjs/operators";

import {Trace} from '../../_models/trace';
import {Project} from "../../_models/project";
import {ClipboardService} from '../../_services/clipboard.service';
import {DataService} from "../../_services/data.service";

import {DateUtil} from '../../_utils/date.utils';

@Directive({
  selector: `ktb-task-item-detail, [ktb-task-item-detail], [ktbTaskItemDetail]`,
  exportAs: 'ktbTaskItemDetail',
})
export class KtbTaskItemDetail {
}

@Component({
  selector: 'ktb-task-item',
  templateUrl: './ktb-task-item.component.html',
  styleUrls: ['./ktb-task-item.component.scss'],
})
export class KtbTaskItemComponent {

  public project$: Observable<Project>;
  public _task: Trace;

  @ViewChild('taskPayloadDialog')
  public taskPayloadDialog: TemplateRef<any>;
  public taskPayloadDialogRef: MatDialogRef<any, any>;

  @Input()
  get task(): Trace {
    return this._task;
  }
  set task(value: Trace) {
    if (this._task !== value) {
      this._task = value;
      this.changeDetectorRef.markForCheck();
    }
  }

  constructor(private changeDetectorRef: ChangeDetectorRef,
              private dataService: DataService,
              private dialog: MatDialog,
              private clipboard: ClipboardService,
              public dateUtil: DateUtil) {
    this.project$ = this.dataService.projects.pipe(
      map(projects => projects ? projects.find(project => {
        return project.projectName === this._task.getProject();
      }) : null)
    );
  }

  showEventPayloadDialog(event, task) {
    event.stopPropagation();
    this.taskPayloadDialogRef = this.dialog.open(this.taskPayloadDialog, { data: task.plainEvent });
  }

  closeEventPayloadDialog() {
    if (this.taskPayloadDialogRef) {
      this.taskPayloadDialogRef.close();
    }
  }

  copyEventPayload(plainEvent: string): void {
    this.clipboard.copy(plainEvent, 'payload');
  }

  isUrl(value: string): boolean {
    try {
      new URL(value);
    } catch (_) {
      return false;
    }
    return true;
  }

}
