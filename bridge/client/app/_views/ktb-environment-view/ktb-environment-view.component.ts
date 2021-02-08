import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Project} from '../../_models/project';
import {Stage} from '../../_models/stage';
import {Observable} from 'rxjs';
import {Trace} from '../../_models/trace';
import {Service} from '../../_models/service';
import {DataService} from '../../_services/data.service';
import {Root} from '../../_models/root';
import {DtOverlayConfig} from '@dynatrace/barista-components/overlay';
import {DtToggleButtonItem} from '@dynatrace/barista-components/toggle-button-group';

@Component({
  selector: 'ktb-environment-view',
  templateUrl: './ktb-environment-view.component.html',
  styleUrls: ['./ktb-environment-view.component.scss']
})
export class KtbEnvironmentViewComponent implements OnInit {
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
    }
  }

  constructor(private dataService: DataService) { }

  ngOnInit(): void {
    this.openApprovals$ = this.dataService.openApprovals;
  }

  selectStage($event, stage: Stage, filterType?: string) {
    this.problemFilterEventButton?.deselect();
    this.evaluationFilterEventButton?.deselect();
    this.approvalFilterEventButton?.deselect();

    this.selectedStage = stage;
    this.filterEventType = filterType;
    $event.stopPropagation();
  }

  trackStage(index: number, stage: Stage) {
    return stage.stageName;
  }

  countOpenApprovals(openApprovals: Trace[], project: Project, stage: Stage, service?: Service) {
    return this.getOpenApprovals(openApprovals, project, stage, service).length;
  }

  getOpenApprovals(openApprovals: Trace[], project: Project, stage: Stage, service?: Service) {
    return openApprovals.filter(approval => approval.data.project === project.projectName && approval.data.stage === stage.stageName && (!service || approval.data.service === service.serviceName));
  }

  findProblemEvent(problemEvents: Root[], service: Service) {
    return problemEvents.find(root => root?.data.service === service.serviceName);
  }

  findFailedRootEvent(failedRootEvents: Root[], service: Service) {
    return failedRootEvents.find(root => root.data.service === service.serviceName);
  }

  selectFilterEvent($event) {
    if ($event.isUserInput) {
      this.filterEventType = $event.source.selected ? $event.value : null;
    }
  }

}
