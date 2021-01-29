import {Component, Input, OnInit} from '@angular/core';
import {Project} from '../../_models/project';

@Component({
  selector: 'ktb-environment-view',
  templateUrl: './ktb-environment-view.component.html',
  styleUrls: ['./ktb-environment-view.component.scss']
})
export class KtbEnvironmentViewComponent implements OnInit {
  private _project: Project;

  @Input()
  get project() {
    return this._project;
  }
<<<<<<< HEAD
  set project(project: Project) {
    if (this._project !== project) {
      this._project = project;
    }
  }

  constructor() {
  }

  ngOnInit(): void {
  }
=======

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

>>>>>>> directory change, environment layout adjustment
}
