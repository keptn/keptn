import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Directive, EventEmitter,
  Input,
  OnInit, Output,
  ViewEncapsulation
} from '@angular/core';
import {coerceArray} from "@angular/cdk/coercion";
import {Root} from "../../_models/root";

@Component({
  selector: 'ktb-root-events-list',
  templateUrl: './ktb-root-events-list.component.html',
  styleUrls: ['./ktb-root-events-list.component.scss'],
  host: {
    class: 'ktb-root-events-list'
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbRootEventsListComponent implements OnInit {

  public _events: Root[] = [];
  public _selectedEvent: Root = null;

  @Output() readonly selectedEventChange = new EventEmitter<Root>();

  @Input()
  get events(): Root[] {
    return this._events;
  }
  set events(value: Root[]) {
    const newValue = coerceArray(value);
    if (this._events !== newValue) {
      // TODO: provide correctly sorted list from API? why is the sorting changed on client side?
      newValue.sort((a, b) => new Date(b.time).getTime() - new Date(a.time).getTime());
      this._events = newValue;
      this._changeDetectorRef.markForCheck();
    }
  }

  @Input()
  get selectedEvent(): Root {
    return this._selectedEvent;
  }
  set selectedEvent(value: Root) {
    if (this._selectedEvent !== value) {
      this._selectedEvent = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit() {
  }

  selectEvent(event: Root) {
    this.selectedEvent = event;
    this.selectedEventChange.emit(this.selectedEvent);
  }

  getEventLabel(key: string): string {
    let label = key;
    switch(key) {
      case "sh.keptn.internal.event.service.create": {
        label = "Service create"
        break;
      }
      case "sh.keptn.event.configuration.change": {
        label = "Configuration change"
        break;
      }
      case "sh.keptn.event.monitoring.configure": {
        label = "Configure monitoring"
        break;
      }
      case "sh.keptn.events.deployment-finished": {
        label = "Deployment finished";
        break;
      }
      case "sh.keptn.events.tests-finished": {
        label = "Tests finished";
        break;
      }
      case "sh.keptn.events.evaluation-done": {
        label = "Evaluation done";
        break;
      }
      case "sh.keptn.internal.event.get-sli": {
        label = "Start SLI retrieval";
        break;
      }
      case "sh.keptn.internal.event.get-sli.done": {
        label = "SLI retrieval done";
        break;
      }
      case "sh.keptn.events.done": {
        label = "Done";
        break;
      }
      case "sh.keptn.events.done": {
        label = "Done";
        break;
      }


      case "sh.keptn.events.done": {
        label = "Done";
        break;
      }
      case "sh.keptn.event.problem.open": {
        label = "Problem open";
        break;
      }
      case "sh.keptn.events.problem": {
        label = "Problem detected";
        break;
      }
      default: {
        //statements;
        break;
      }
    }

    return label;
  }

}
