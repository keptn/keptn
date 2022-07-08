import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { switchMap, takeUntil } from 'rxjs/operators';
import { ActivatedRoute, Router } from '@angular/router';
import { combineLatest, Subject } from 'rxjs';
import { DataService } from '../../../_services/data.service';
import { DtTableDataSource } from '@dynatrace/barista-components/table';
import { Location } from '@angular/common';
import { ServiceDeploymentInformation as sdi, ServiceState } from '../../../../../shared/models/service-state';
import { DeploymentInformationSelection } from '../../../_interfaces/deployment-selection';

class DeploymentInformation implements sdi {
  keptnContext!: string;
  name!: string;
  stages!: { name: string; hasOpenRemediations: boolean; time: string }[];
  version?: string;
}

@Component({
  selector: 'ktb-deployment-list[service]',
  templateUrl: './ktb-deployment-list.component.html',
  styleUrls: ['./ktb-deployment-list.component.scss'],
})
export class KtbDeploymentListComponent implements OnInit, OnDestroy {
  private _service?: ServiceState;
  private projectName?: string;
  private readonly unsubscribe$ = new Subject<void>();
  public _selectedDeploymentInfo?: DeploymentInformationSelection;
  public dataSource = new DtTableDataSource<DeploymentInformation>();
  public loading = false;
  public DeploymentClass = DeploymentInformation;

  @Output() selectedDeploymentInfoChange: EventEmitter<DeploymentInformationSelection> = new EventEmitter();

  @Input()
  get service(): ServiceState | undefined {
    return this._service;
  }

  set service(service: ServiceState | undefined) {
    if (this._service !== service) {
      this._service = service;
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

  public ngOnInit(): void {
    const params$ = this.route.params.pipe(takeUntil(this.unsubscribe$));

    const project$ = params$.pipe(
      switchMap((params) => this.dataService.getProject(params.projectName)),
      takeUntil(this.unsubscribe$)
    );

    combineLatest([project$, params$]).subscribe(([project]) => {
      this.projectName = project?.projectName;
      this.updateDataSource();
    });
  }

  private updateDataSource(count = -1): void {
    this.dataSource.data =
      (count !== -1 ? this.service?.deploymentInformation.slice(0, count) : this.service?.deploymentInformation) ?? [];
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

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
