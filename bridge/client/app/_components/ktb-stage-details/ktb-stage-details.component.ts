import {ChangeDetectorRef, Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {Project} from '../../_models/project';
import {Stage} from '../../_models/stage';
import {Observable} from 'rxjs';
import {Trace} from '../../_models/trace';
import {DataService} from '../../_services/data.service';
import {DtToggleButtonItem} from '@dynatrace/barista-components/toggle-button-group';
import {DtOverlayConfig} from '@dynatrace/barista-components/overlay';
import {Service} from '../../_models/service';
import {Root} from '../../_models/root';

@Component({
  selector: 'ktb-stage-details',
  templateUrl: './ktb-stage-details.component.html',
  styleUrls: ['./ktb-stage-details.component.scss']
})
export class KtbStageDetailsComponent implements OnInit {
  public _project: Project;
  public selectedStage: Stage = null;
  public openApprovals$: Observable<Trace[]>;
  public filterEventType: string = null;
  public overlayConfig: DtOverlayConfig = {
    pinnable: true
  };

  @ViewChild('problemFilterEventButton') public problemFilterEventButton: DtToggleButtonItem<string>;
  @ViewChild('evaluationFilterEventButton') public evaluationFilterEventButton: DtToggleButtonItem<string>;
  @ViewChild('approvalFilterEventButton') public approvalFilterEventButton: DtToggleButtonItem<string>;

  @Input()
  get project() {
    return this._project;
  }

  set project(project: Project) {
    if (this._project !== project) {
      this._project = project;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private dataService: DataService, private _changeDetectorRef: ChangeDetectorRef) {
  }

  ngOnInit(): void {
    this.openApprovals$ = this.dataService.openApprovals;
  }

  selectStage($event) {
    this.problemFilterEventButton?.deselect();
    this.evaluationFilterEventButton?.deselect();
    this.approvalFilterEventButton?.deselect();
    this.selectedStage = $event.stage;
    this.filterEventType = $event.filterType;
  }

  selectFilterEvent($event) {
    if ($event.isUserInput) {
      this.filterEventType = $event.source.selected ? $event.value : null;
    }
  }

  countOpenApprovals(project: Project, stage: Stage, service?: Service): number {
    return this.getOpenApprovals(project, stage, service).length;
  }

  getOpenApprovals(project: Project, stage: Stage, service?: Service): Trace[] {
    return this.dataService.getOpenApprovals(project, stage, service);
  }

  findProblemEvent(problemEvents: Root[], service: Service): Root {
    return problemEvents.find(root => root?.data.service === service.serviceName);
  }

  findFailedRootEvent(failedRootEvents: Root[], service: Service): Root {
    return failedRootEvents.find(root => root.data.service === service.serviceName);
  }

}
