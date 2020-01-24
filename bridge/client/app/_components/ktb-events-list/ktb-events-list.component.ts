import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Input,
  OnInit,
  ViewEncapsulation
} from '@angular/core';
import {Trace} from "../../_models/trace";
import DateUtil from "../../_utils/date.utils";

@Component({
  selector: 'ktb-events-list',
  templateUrl: './ktb-events-list.component.html',
  styleUrls: ['./ktb-events-list.component.scss'],
  host: {
    class: 'ktb-root-events-list'
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbEventsListComponent implements OnInit {

  public _events: Trace[] = [];

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

  constructor(private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit() {
  }

  getCalendarFormats() {
    return DateUtil.getCalendarFormats();
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
