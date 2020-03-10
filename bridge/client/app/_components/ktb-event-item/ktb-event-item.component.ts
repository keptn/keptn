import {ChangeDetectorRef, Component, Directive, Input, OnInit} from '@angular/core';
import {Trace} from "../../_models/trace";
import DateUtil from "../../_utils/date.utils";

@Directive({
  selector: `ktb-event-item-detail, [ktb-event-item-detail], [ktbEventItemDetail]`,
  exportAs: 'ktbEventItemDetail',
})
export class KtbEventItemDetail {}

@Component({
  selector: 'ktb-event-item',
  templateUrl: './ktb-event-item.component.html',
  styleUrls: ['./ktb-event-item.component.scss']
})
export class KtbEventItemComponent implements OnInit {

  public _event: Trace;

  @Input()
  get event(): Trace {
    return this._event;
  }
  set event(value: Trace) {
    if (this._event !== value) {
      this._event = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit() {
  }

  getEventLabel(event: Trace): string {
    let label = event.type;
    switch(event.type) {
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
        if(event.data.State === "RESOLVED")
          label = "Problem resolved";
        else
          label = "Problem detected";
        break;
      }
      case "sh.keptn.event.problem.close": {
        label = "Problem closed";
        break;
      }
      default: {
        //statements;
        break;
      }
    }

    return label;
  }

  getCalendarFormat() {
    return DateUtil.getCalendarFormats().sameElse;
  }

}
