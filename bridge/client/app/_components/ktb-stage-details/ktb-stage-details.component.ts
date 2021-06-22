import {ChangeDetectorRef, Component, Input, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {DtToggleButtonItem} from '@dynatrace/barista-components/toggle-button-group';
import {DtOverlayConfig} from '@dynatrace/barista-components/overlay';

import {Project} from '../../_models/project';
import {Stage} from '../../_models/stage';
import {Service} from '../../_models/service';
import {Root} from '../../_models/root';

import {DataService} from '../../_services/data.service';
import {filter, takeUntil} from 'rxjs/operators';
import {Subject} from 'rxjs';

@Component({
  selector: 'ktb-stage-details',
  templateUrl: './ktb-stage-details.component.html',
  styleUrls: ['./ktb-stage-details.component.scss']
})
export class KtbStageDetailsComponent implements OnInit, OnDestroy {
  public _project: Project;
  public selectedStage: Stage = null;
  public filterEventType: string = null;
  public overlayConfig: DtOverlayConfig = {
    pinnable: true
  };
  public isQualityGatesOnly: boolean;
  private _filteredServices: string[];
  private readonly unsubscribe$ = new Subject<void>();

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
      this.selectedStage = null;
      this._changeDetectorRef.markForCheck();
    }
  }

  get filteredServices(): string[] {
    return this._filteredServices;
  }

  set filteredServices(services: string[]) {
    this._filteredServices = services;
    this.resetFilter();
    this._changeDetectorRef.markForCheck();
  }

  constructor(private dataService: DataService, private _changeDetectorRef: ChangeDetectorRef) {
  }

  ngOnInit(): void {
    this.dataService.isQualityGatesOnly.pipe(
      takeUntil(this.unsubscribe$)
    ).subscribe(isQualityGatesOnly => {this.isQualityGatesOnly = isQualityGatesOnly});
  }

  selectStage($event) {
    this.selectedStage = $event.stage;
    if (this.filterEventType !== $event.filterType) {
      this.resetFilter($event.filterType);
    }
  }

  private resetFilter(eventType = null): void {
    this.problemFilterEventButton?.deselect();
    this.evaluationFilterEventButton?.deselect();
    this.approvalFilterEventButton?.deselect();
    this.filterEventType = eventType;
  }

  selectFilterEvent($event) {
    if ($event.isUserInput) {
      this.filterEventType = $event.source.selected ? $event.value : null;
    }
  }

  hasProblemEvent(problemEvents: Root[], service: Service): boolean {
    return problemEvents.filter(root => root?.data.service === service.serviceName).length > 0;
  }

  findProblemEvent(problemEvents: Root[], service: Service): Root[] {
    return problemEvents.filter(root => root?.data.service === service.serviceName);
  }

  findFailedRootEvent(failedRootEvents: Root[], service: Service): Root {
    return failedRootEvents.find(root => root.data.service === service.serviceName);
  }

  getServiceLink(service) {
    return ['service', service.serviceName, 'context', service.deploymentContext, 'stage', service.stage];
  }

  public filterServices(services: Service[]): Service[] {
    return this.filteredServices.length === 0 ? services : services.filter(service => this.filteredServices.includes(service.serviceName));
  }

  public filterRoots(roots: Root[]): Root[] {
    return this.filteredServices.length === 0 ? roots : roots?.filter(root => this.filteredServices.includes(root.getService()));

  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
