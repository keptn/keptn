import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  HostBinding,
  Input,
  Output,
  ViewEncapsulation,
} from '@angular/core';
import { DateUtil } from '../../../_utils/date.utils';
import { Project } from '../../../_models/project';
import { Sequence } from '../../../_models/sequence';
import { ISequence } from '../../../../../shared/interfaces/sequence';

@Component({
  selector: 'ktb-root-events-list',
  templateUrl: './ktb-root-events-list.component.html',
  styleUrls: ['./ktb-root-events-list.component.scss'],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbRootEventsListComponent {
  @HostBinding('class') cls = 'ktb-root-events-list';
  public project?: Project;
  public _events: Sequence[] = [];
  public _selectedEvent?: Sequence;

  @Output() readonly selectedEventChange = new EventEmitter<{ sequence: Sequence; stage?: string }>();

  @Input()
  get events(): Sequence[] {
    return this._events;
  }

  set events(value: Sequence[]) {
    if (this._events !== value) {
      this._events = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  get selectedEvent(): Sequence | undefined {
    return this._selectedEvent;
  }

  set selectedEvent(value: Sequence | undefined) {
    if (this._selectedEvent !== value) {
      this._selectedEvent = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, public dateUtil: DateUtil) {}

  selectEvent(sequence: ISequence, stage?: string): void {
    // Refactor without using cast to Sequence
    // use ISequence instead
    this.selectedEvent = <Sequence>sequence;
    this.selectedEventChange.emit({ sequence: <Sequence>sequence, stage });
  }

  identifyEvent(index: number, item: Sequence): string | undefined {
    return item?.time;
  }

  public getShortType(type: string | undefined): string | undefined {
    return type ? Sequence.getShortType(type) : undefined;
  }
}
