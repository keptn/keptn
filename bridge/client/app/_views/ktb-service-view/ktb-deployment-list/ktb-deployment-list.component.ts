import { Component, EventEmitter, Input, Output } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../../../_services/data.service';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { Location } from '@angular/common';
import { ServiceDeploymentInformation as sdi, ServiceState } from '../../../../../shared/models/service-state';
import { DeploymentInformationSelection } from '../../../_interfaces/deployment-selection';

class DeploymentInformation implements sdi {
  keptnContext!: string;
  name!: string;
  stages!: {
    name: string;
    time: string;
  }[];
  version?: string;
}

@Component({
  selector: 'ktb-deployment-list[service][projectName]',
  templateUrl: './ktb-deployment-list.component.html',
  styleUrls: ['./ktb-deployment-list.component.scss'],
})
export class KtbDeploymentListComponent {
  private _service?: ServiceState;
  public _selectedDeploymentInfo?: DeploymentInformationSelection;
  public dataSource = new DtTableDataSource<DeploymentInformation>();
  public loading = false;
  public DeploymentClass = DeploymentInformation;

  @Output() selectedDeploymentInfoChange: EventEmitter<DeploymentInformationSelection> = new EventEmitter();
  @Input() projectName = '';

  @Input()
  get service(): ServiceState | undefined {
    return this._service;
  }

  set service(service: ServiceState | undefined) {
    if (this._service !== service) {
      this._service = service;
      this.updateDataSource();
    }
  }
  @Input()
  get selectedDeploymentInfo(): DeploymentInformationSelection | undefined {
    return this._selectedDeploymentInfo;
  }
  set selectedDeploymentInfo(deployment: DeploymentInformationSelection | undefined) {
    if (this._selectedDeploymentInfo !== deployment) {
      this._selectedDeploymentInfo = deployment;
    }
  }

  constructor(
    private route: ActivatedRoute,
    private dataService: DataService,
    private router: Router,
    private location: Location
  ) {}

  private updateDataSource(): void {
    this.dataSource.data = this.service?.deploymentInformation ?? [];
  }

  public selectDeployment(deploymentInformation: DeploymentInformation, stageName?: string): void {
    if (
      this.selectedDeploymentInfo?.deploymentInformation.keptnContext !== deploymentInformation.keptnContext ||
      stageName
    ) {
      stageName ??= deploymentInformation.stages[deploymentInformation.stages.length - 1].name;
      const routeUrl = this.router.createUrlTree([
        '/project',
        this.projectName,
        'service',
        this.service?.name,
        'context',
        deploymentInformation.keptnContext,
        'stage',
        stageName,
      ]);
      this.location.go(routeUrl.toString());
      this.selectedDeploymentInfo = { deploymentInformation, stage: stageName };
      this.selectedDeploymentInfoChange.emit(this.selectedDeploymentInfo);
    }
  }

  public selectStage(deployment: DeploymentInformation, stageName: string, $event: MouseEvent): void {
    $event.stopPropagation();
    this.selectDeployment(deployment, stageName);
  }
}
