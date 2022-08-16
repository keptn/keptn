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
import { SequenceState } from '../../../_models/sequenceState';
import { ISequenceState } from '../../../../../shared/interfaces/sequence';

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
  public _events: SequenceState[] = [];
  public _selectedEvent?: SequenceState;

  @Output() readonly selectedEventChange = new EventEmitter<{ sequence: SequenceState; stage?: string }>();

  @Input()
  get events(): SequenceState[] {
    return this._events;
  }

  set events(value: SequenceState[]) {
    if (this._events !== value) {
      this._events = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  get selectedEvent(): SequenceState | undefined {
    return this._selectedEvent;
  }

  set selectedEvent(value: SequenceState | undefined) {
    if (this._selectedEvent !== value) {
      this._selectedEvent = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, public dateUtil: DateUtil) {}

  selectEvent(sequence: ISequenceState, stage?: string): void {
    // Refactor without using cast to Sequence
    // use ISequence instead
    this.selectedEvent = <SequenceState>sequence;
    this.selectedEventChange.emit({ sequence: <SequenceState>sequence, stage });
  }

  identifyEvent(_index: number, item: SequenceState): string | undefined {
    return item?.time;
  }

  public getShortType(type: string | undefined): string | undefined {
    return type ? SequenceState.getShortType(type) : undefined;
  }
}
