import { ChangeDetectionStrategy, Component, HostBinding, Input } from '@angular/core';
import { Router } from '@angular/router';
import { Location } from '@angular/common';
import { Trace } from '../../../_models/trace';
import { DateUtil } from '../../../_utils/date.utils';
import { ApiService } from '../../../_services/api.service';

@Component({
  selector: 'ktb-sequence-tasks-list[tasks][focusedEventId]',
  templateUrl: './ktb-sequence-tasks-list.component.html',
  styleUrls: ['./ktb-sequence-tasks-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSequenceTasksListComponent {
  @HostBinding('class') cls = 'ktb-sequence-tasks-list';
  private _tasks: Trace[] = [];
  private _focusedEventId?: string;
  private currentScrollElement?: HTMLDivElement;

  @Input()
  get tasks(): Trace[] {
    return this._tasks;
  }

  set tasks(value: Trace[]) {
    if (this._tasks !== value) {
      this._tasks = value;
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
    }
    if (!value) {
      this.focusLastSequence();
    }
  }

  constructor(
    private router: Router,
    private location: Location,
    public dateUtil: DateUtil,
    private apiService: ApiService
  ) {}

  identifyEvent(_index: number, item: Trace): string {
    return item.id;
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
    if (!event.project) {
      return;
    }
    const sequenceFilters = this.apiService.getSequenceFilters(event.project);
    const routeUrl = this.router.createUrlTree(
      ['/project', event.project, 'sequence', event.shkeptncontext, 'event', event.id],
      { queryParams: sequenceFilters }
    );
    this._focusedEventId = event.id;
    this.location.go(routeUrl.toString());
  }

  focusLastSequence(): void {
    if (
      this.tasks.length &&
      (!this.focusedEventId || !this.tasks.some((seq) => seq.findTrace((t) => t.id === this.focusedEventId)))
    ) {
      this._focusedEventId = this.tasks[this.tasks.length - 1].id;
    }
  }

  isFocusedTask(task: Trace): boolean {
    return !!task.findTrace((t) => t.id === this.focusedEventId);
  }
}
