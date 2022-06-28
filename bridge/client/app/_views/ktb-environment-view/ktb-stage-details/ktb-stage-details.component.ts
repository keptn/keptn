import { Component, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { DtToggleButtonChange, DtToggleButtonItem } from '@dynatrace/barista-components/toggle-button-group';
import { DtOverlayConfig } from '@dynatrace/barista-components/overlay';
import { Project } from '../../../_models/project';
import { Stage } from '../../../_models/stage';
import { Service } from '../../../_models/service';
import { DataService } from '../../../_services/data.service';
import { takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';

export type ServiceFilterType = 'evaluation' | 'problem' | 'approval' | undefined;

@Component({
  selector: 'ktb-stage-details',
  templateUrl: './ktb-stage-details.component.html',
  styleUrls: ['./ktb-stage-details.component.scss'],
})
export class KtbStageDetailsComponent implements OnInit, OnDestroy {
  public _project?: Project;
  public selectedStage?: Stage;
  public filterEventType: ServiceFilterType;
  public overlayConfig: DtOverlayConfig = {
    pinnable: true,
  };
  public isQualityGatesOnly = false;
  public filteredServices: string[] = [];
  private readonly unsubscribe$ = new Subject<void>();

  @ViewChild('problemFilterEventButton') public problemFilterEventButton?: DtToggleButtonItem<string>;
  @ViewChild('evaluationFilterEventButton') public evaluationFilterEventButton?: DtToggleButtonItem<string>;
  @ViewChild('approvalFilterEventButton') public approvalFilterEventButton?: DtToggleButtonItem<string>;

  @Input()
  get project(): Project | undefined {
    return this._project;
  }

  set project(project: Project | undefined) {
    if (this._project !== project) {
      this._project = project;
      this.selectedStage = undefined;
    }
  }

  constructor(private dataService: DataService) {}

  ngOnInit(): void {
    this.dataService.isQualityGatesOnly.pipe(takeUntil(this.unsubscribe$)).subscribe((isQualityGatesOnly) => {
      this.isQualityGatesOnly = isQualityGatesOnly;
    });
  }

  selectStage($event: { stage: Stage; filterType: ServiceFilterType }): void {
    this.selectedStage = $event.stage;
    if (this.filterEventType !== $event.filterType) {
      this.resetFilter($event.filterType);
    }
  }

  private resetFilter(eventType: ServiceFilterType): void {
    this.problemFilterEventButton?.deselect();
    this.evaluationFilterEventButton?.deselect();
    this.approvalFilterEventButton?.deselect();
    this.filterEventType = eventType;
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  selectFilterEvent($event: DtToggleButtonChange<any>): void {
    if ($event.isUserInput) {
      this.filterEventType = $event.source.selected ? $event.value : null;
    }
  }

  getServiceLink(service: Service): string[] {
    return ['service', service.serviceName, 'context', service.deploymentContext ?? '', 'stage', service.stage];
  }

  public filterServices(services: Service[], type: ServiceFilterType): Service[] {
    const filteredServices =
      this.filteredServices.length === 0
        ? services
        : services.filter((service) => this.filteredServices.includes(service.serviceName));
    if (this.filterEventType && filteredServices.length === 0 && this.filterEventType === type) {
      this.resetFilter(undefined);
    }
    return filteredServices;
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
