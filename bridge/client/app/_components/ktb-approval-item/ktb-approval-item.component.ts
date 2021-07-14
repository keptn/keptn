import {ChangeDetectorRef, Component, Input} from '@angular/core';
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
  public approvalResult: boolean = null;

  public overlayConfig: DtOverlayConfig = {
    pinnable: true
  };

  @Input() isSequence = false;

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

  handleApproval(approval, result) {
    this.dataService.sendApprovalEvent(approval, result);
    this.approvalResult = result;
  }

  public getDeploymentEvaluation(project: Project): Trace | undefined {
    return this.isSequence ? project.getDeploymentEvaluationOfSequence(this.event) : project.getDeploymentEvaluation(this.event);
  }

}
