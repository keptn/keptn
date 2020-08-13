import { ChangeDetectorRef, Component, Input } from '@angular/core';
import { Trace } from '../../_models/trace';
import {Observable} from "rxjs";
import {Project} from "../../_models/project";
import {map} from "rxjs/operators";
import {DataService} from "../../_services/data.service";
import {DtOverlayConfig} from "@dynatrace/barista-components/overlay";

@Component({
  selector: 'ktb-approval-item',
  templateUrl: './ktb-approval-item.component.html',
  styleUrls: ['./ktb-approval-item.component.scss'],
})
export class KtbApprovalItemComponent {

  public project$: Observable<Project>;
  public _event: Trace;

  public overlayConfig: DtOverlayConfig = {
    pinnable: true
  };

  @Input()
  get event(): Trace {
    return this._event;
  }

  set event(value: Trace) {
    if (this._event !== value) {
      this._event = value;
      this.changeDetectorRef.markForCheck();
    }
  }

  constructor(private changeDetectorRef: ChangeDetectorRef, private dataService: DataService) {
    this.project$ = this.dataService.projects.pipe(
      map(projects => projects ? projects.find(project => {
        return project.projectName === this._event.getProject();
      }) : null)
    );
  }

  approveDeployment(approval) {
    this.dataService.sendApprovalEvent(approval, true);
  }

  declineDeployment(approval) {
    this.dataService.sendApprovalEvent(approval, false);
  }

}
