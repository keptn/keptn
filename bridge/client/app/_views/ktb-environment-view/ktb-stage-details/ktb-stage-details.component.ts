import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { DtToggleButtonChange, DtToggleButtonItem } from '@dynatrace/barista-components/toggle-button-group';
import { DtOverlayConfig } from '@dynatrace/barista-components/overlay';
import { Project } from '../../../_models/project';
import { Stage } from '../../../_models/stage';
import { Service } from '../../../_models/service';
import { DataService } from '../../../_services/data.service';
import { Observable } from 'rxjs';
import { ISelectedStageInfo } from '../ktb-environment-view.component';

export type ServiceFilterType = 'evaluation' | 'problem' | 'approval' | undefined;

@Component({
  selector: 'ktb-stage-details',
  templateUrl: './ktb-stage-details.component.html',
  styleUrls: ['./ktb-stage-details.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbStageDetailsComponent {
  public filterEventType: ServiceFilterType;
  public overlayConfig: DtOverlayConfig = {
    pinnable: true,
  };
  public isQualityGatesOnly$: Observable<boolean>;
  public filteredServices: string[] = [];
  private _selectedStageInfo?: ISelectedStageInfo;

  @Input() project?: Project;
  @Input() set selectedStageInfo(stageInfo: ISelectedStageInfo | undefined) {
    this._selectedStageInfo = stageInfo;
    if (stageInfo && this.filterEventType !== stageInfo.filterType) {
      this.resetFilter(stageInfo.filterType);
    }
  }
  get selectedStageInfo(): ISelectedStageInfo | undefined {
    return this._selectedStageInfo;
  }
  @Output() selectedStageInfoChange = new EventEmitter<ISelectedStageInfo>();

  @ViewChild('problemFilterEventButton') public problemFilterEventButton?: DtToggleButtonItem<string>;
  @ViewChild('evaluationFilterEventButton') public evaluationFilterEventButton?: DtToggleButtonItem<string>;
  @ViewChild('approvalFilterEventButton') public approvalFilterEventButton?: DtToggleButtonItem<string>;

  constructor(private dataService: DataService) {
    this.isQualityGatesOnly$ = this.dataService.isQualityGatesOnly;
  }

  private resetFilter(eventType: ServiceFilterType): void {
    this.problemFilterEventButton?.deselect();
    this.evaluationFilterEventButton?.deselect();
    this.approvalFilterEventButton?.deselect();
    this.filterEventType = eventType;
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  selectFilterEvent(stage: Stage, $event: DtToggleButtonChange<any>): void {
    if ($event.isUserInput) {
      this.filterEventType = $event.source.selected ? $event.value : null;

      this._selectedStageInfo = { stage, filterType: this.filterEventType };
      this.selectedStageInfoChange.emit(this.selectedStageInfo);
    }
  }

  getServiceLink(service: Service): string[] {
    return [
      '/project',
      this.project?.projectName ?? '',
      'service',
      service.serviceName,
      'context',
      service.deploymentContext ?? '',
      'stage',
      service.stage,
    ];
  }

  public filterServices(stage: Stage, services: Service[], type: ServiceFilterType): Service[] {
    const filteredServices =
      this.filteredServices.length === 0
        ? services
        : services.filter((service) => this.filteredServices.includes(service.serviceName));
    if (this.filterEventType && filteredServices.length === 0 && this.filterEventType === type) {
      this.resetFilter(undefined);

      this._selectedStageInfo = { stage, filterType: undefined };
      this.selectedStageInfoChange.emit(this.selectedStageInfo);
    }
    return filteredServices;
  }
}
