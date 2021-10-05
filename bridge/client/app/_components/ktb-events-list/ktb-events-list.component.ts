import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  HostBinding,
  Input,
  ViewEncapsulation,
} from '@angular/core';
import { Trace } from '../../_models/trace';
import { Router } from '@angular/router';
import { Location } from '@angular/common';

@Component({
  selector: 'ktb-events-list',
  templateUrl: './ktb-events-list.component.html',
  styleUrls: ['./ktb-events-list.component.scss'],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbEventsListComponent {
  @HostBinding('class') cls = 'ktb-events-list';
  public _events: Trace[] = [];
  public _focusedEventId?: string;
  private currentScrollElement?: HTMLDivElement;

  @Input()
  get events(): Trace[] {
    return this._events;
  }
  set events(value: Trace[]) {
    if (this._events !== value) {
      this._events = value;
      this._changeDetectorRef.markForCheck();
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

  constructor(private router: Router, private location: Location, private _changeDetectorRef: ChangeDetectorRef) {}

  identifyEvent(index: number, item: Trace): Date | undefined | null {
    return item ? item.time : null;
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
    if (event.project && event.service) {
      const routeUrl = this.router.createUrlTree([
        '/project',
        event.project,
        event.service,
        event.shkeptncontext,
        event.id,
      ]);
      this.location.go(routeUrl.toString());
    }
  }

  isInvalidated(event: Trace): boolean {
    return !!this.events.find((e) => e.isEvaluationInvalidation() && e.triggeredid === event.id);
  }
}
