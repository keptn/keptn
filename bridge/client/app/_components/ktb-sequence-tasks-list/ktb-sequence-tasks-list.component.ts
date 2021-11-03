import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  HostBinding,
  Input,
  OnDestroy,
  OnInit,
  ViewEncapsulation,
} from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Location } from '@angular/common';
import { Trace } from '../../_models/trace';
import { DateUtil } from '../../_utils/date.utils';
import { takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';

@Component({
  selector: 'ktb-sequence-tasks-list[tasks][stage]',
  templateUrl: './ktb-sequence-tasks-list.component.html',
  styleUrls: ['./ktb-sequence-tasks-list.component.scss'],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSequenceTasksListComponent implements OnInit, OnDestroy {
  @HostBinding('class') cls = 'ktb-sequence-tasks-list';
  public _tasks: Trace[] = [];
  public _stage?: string;
  public _focusedEventId?: string;
  private readonly unsubscribe$ = new Subject<void>();
  private currentScrollElement?: HTMLDivElement;

  @Input()
  get tasks(): Trace[] {
    return this._tasks;
  }
  set tasks(value: Trace[]) {
    if (this._tasks !== value) {
      this._tasks = value;
      this._changeDetectorRef.markForCheck();
      this.focusLastSequence();
    }
  }

  @Input()
  get stage(): string | undefined {
    return this._stage;
  }
  set stage(value: string | undefined) {
    if (this._stage !== value) {
      this._stage = value;
      this._changeDetectorRef.markForCheck();
      this.focusLastSequence();
    }
  }

  @Input()
  get focusedEventId(): string | undefined {
    return this._focusedEventId;
  }
  set focusedEventId(value: string | undefined) {
    if (this._focusedEventId !== value) {
      this._focusedEventId = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(
    private router: Router,
    private location: Location,
    private _changeDetectorRef: ChangeDetectorRef,
    public dateUtil: DateUtil,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    this.route.params.pipe(takeUntil(this.unsubscribe$)).subscribe((params) => {
      if (params.eventId) {
        this.focusedEventId = params.eventId;
      } else {
        this.focusLastSequence();
      }
    });
  }

  identifyEvent(index: number, item: Trace): Date | undefined {
    return item?.time ? new Date(item.time) : undefined;
  }

  scrollIntoView(element: HTMLDivElement): boolean {
    if (element !== this.currentScrollElement) {
      this.currentScrollElement = element;
      setTimeout(() => {
        element.scrollIntoView({ behavior: 'smooth' });
      }, 0);
    }
    return true;
  }

  focusEvent(event: Trace): void {
    if (event.project) {
      const routeUrl = this.router.createUrlTree([
        '/project',
        event.project,
        'sequence',
        event.shkeptncontext,
        'event',
        event.id,
      ]);
      this.location.go(routeUrl.toString());
    }
  }

  focusLastSequence(): void {
    if (
      this.stage &&
      !this.getTasksByStage(this.tasks, this.stage).some(
        (seq) => seq.id === this.focusedEventId || seq.findTrace((t) => t.id === this.focusedEventId)
      )
    ) {
      this.focusedEventId = this.tasks
        .slice()
        .reverse()
        .find((t) => t.stage === this.stage)?.id;
    }
  }

  getTasksByStage(tasks: Trace[], stage: string): Trace[] {
    return tasks.filter((t) => t.data?.stage === stage);
  }

  isInvalidated(event: Trace): boolean {
    return !!this.tasks.find((e) => e.isEvaluationInvalidation() && e.triggeredid === event.id);
  }

  isFocusedTask(task: Trace): boolean {
    return task.id === this.focusedEventId || !!task.findTrace((t) => t.id === this.focusedEventId);
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
