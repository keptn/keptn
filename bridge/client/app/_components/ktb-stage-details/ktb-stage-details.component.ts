import {ChangeDetectorRef, Component, Input, OnDestroy, OnInit, ViewChild} from '@angular/core';
import { DtToggleButtonChange, DtToggleButtonItem } from '@dynatrace/barista-components/toggle-button-group';
import {DtOverlayConfig} from '@dynatrace/barista-components/overlay';
import {Project} from '../../_models/project';
import {Stage} from '../../_models/stage';
import {Service} from '../../_models/service';
import {DataService} from '../../_services/data.service';
import {takeUntil} from 'rxjs/operators';
import {Subject} from 'rxjs';
import { Sequence } from '../../_models/sequence';

@Component({
  selector: 'ktb-stage-details',
  templateUrl: './ktb-stage-details.component.html',
  styleUrls: ['./ktb-stage-details.component.scss']
})
export class KtbStageDetailsComponent implements OnInit, OnDestroy {
  public _project?: Project;
  public selectedStage?: Stage;
  public filterEventType?: string;
  public overlayConfig: DtOverlayConfig = {
    pinnable: true
  };
  public isQualityGatesOnly = false;
  private _filteredServices: string[] = [];
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
    ).subscribe(isQualityGatesOnly => {
      this.isQualityGatesOnly = isQualityGatesOnly;
    });
  }

  selectStage($event: {stage: Stage, filterType?: string}) {
    this.selectedStage = $event.stage;
    if (this.filterEventType !== $event.filterType) {
      this.resetFilter($event.filterType);
    }
  }

  private resetFilter(eventType?: string): void {
    this.problemFilterEventButton?.deselect();
    this.evaluationFilterEventButton?.deselect();
    this.approvalFilterEventButton?.deselect();
    this.filterEventType = eventType;
  }

  // tslint:disable-next-line:no-any
  selectFilterEvent($event: DtToggleButtonChange<any>) {
    if ($event.isUserInput) {
      this.filterEventType = $event.source.selected ? $event.value : null;
    }
  }

  hasProblemEvent(problemEvents: Sequence[], service: Service): boolean {
    return problemEvents.some(p => service.openRemediations.some(r => r.service === p.service));
  }

  findProblemEvent(problemEvents: Sequence[], service: Service): Sequence[] {
    return service.openRemediations;
  }

  findFailedRootEvent(failedRootEvents: Sequence[], service: Service): Sequence | undefined {
    return failedRootEvents.find(sequence => sequence.service === service.serviceName);
  }

  getServiceLink(service: Service) {
    return ['service', service.serviceName, 'context', service.deploymentContext, 'stage', service.stage];
  }

  public filterServices(services: Service[]): Service[] {
    return this.filteredServices.length === 0 ? services : services.filter(service => this.filteredServices.includes(service.serviceName));
  }

  public filterSequences(sequences: Sequence[]): Sequence[] {
    return this.filteredServices.length === 0
          ? sequences
          : sequences?.filter(sequence => sequence.service ? this.filteredServices.includes(sequence.service) : false);
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
