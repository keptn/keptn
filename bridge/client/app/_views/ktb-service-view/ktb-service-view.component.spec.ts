import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbServiceViewComponent } from './ktb-service-view.component';
import { ActivatedRoute, convertToParamMap, ParamMap } from '@angular/router';
import { BehaviorSubject } from 'rxjs';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { AppUtils, POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import {
  ServiceStateResponse,
  ServiceStateResponseWithoutRemediationsMock,
} from '../../../../shared/fixtures/service-state-response.mock';
import {
  ServiceDeploymentMock,
  ServiceDeploymentWithApprovalMock,
} from '../../../../shared/fixtures/service-deployment-response.mock';
import { Deployment } from '../../_models/deployment';
import { Location } from '@angular/common';
import { ServiceState } from '../../_models/service-state';
import { DataService } from '../../_services/data.service';
import { KtbServiceViewModule } from './ktb-service-view.module';

describe('KtbServiceViewComponent', () => {
  let component: KtbServiceViewComponent;
  let fixture: ComponentFixture<KtbServiceViewComponent>;
  const projectName = 'sockshop';
  let paramsSubject: BehaviorSubject<ParamMap>;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    paramsSubject = new BehaviorSubject(convertToParamMap({}));
    await TestBed.configureTestingModule({
      declarations: [],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: paramsSubject.asObservable(),
          },
        },
        {
          provide: POLLING_INTERVAL_MILLIS,
          useValue: 0,
        },
      ],
      imports: [KtbServiceViewModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbServiceViewComponent);
    httpMock = TestBed.inject(HttpTestingController);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load and select service states', () => {
    const selectedDeployment = ServiceStateResponse[0].deploymentInformation[2];
    paramsSubject.next(
      convertToParamMap({
        projectName,
        serviceName: 'carts',
        stage: 'staging',
        shkeptncontext: selectedDeployment.keptnContext,
      })
    );
    loadDeployment(
      selectedDeployment.keptnContext,
      selectedDeployment.stages.some((st) => st.hasOpenRemediations)
    );

    expect(component.selectedDeployment?.deploymentInformation.deployment).not.toBeUndefined();
    expect(component.selectedDeployment?.stage).toEqual('staging');
  });

  it('should load and select service states and select latest stage if param does not exist', () => {
    const selectedDeployment = ServiceStateResponse[0].deploymentInformation[2];
    const serviceName = 'carts';
    const location = TestBed.inject(Location);
    const locationSpy = jest.spyOn(location, 'go');
    paramsSubject.next(
      convertToParamMap({
        projectName,
        serviceName,
        stage: 'notExisting',
        shkeptncontext: selectedDeployment.keptnContext,
      })
    );
    loadDeployment(
      selectedDeployment.keptnContext,
      selectedDeployment.stages.some((st) => st.hasOpenRemediations)
    );

    expect(component.selectedDeployment?.deploymentInformation.deployment).not.toBeUndefined();
    expect(component.selectedDeployment?.stage).toEqual('production');
    expect(locationSpy).toHaveBeenCalledWith(
      `/project/${projectName}/service/${serviceName}/context/${selectedDeployment.keptnContext}/stage/production`
    );
  });

  it('should load and remove context param, if it does not exist', () => {
    const serviceName = 'carts';
    const location = TestBed.inject(Location);
    const locationSpy = jest.spyOn(location, 'go');
    paramsSubject.next(
      convertToParamMap({
        projectName,
        serviceName,
        stage: 'staging',
        shkeptncontext: 'notExisting',
      })
    );
    httpMock.expectOne(`./api/project/${projectName}/serviceStates`).flush(AppUtils.copyObject(ServiceStateResponse));

    expect(component.selectedDeployment).toBeUndefined();
    expect(locationSpy).toHaveBeenCalledWith(`/project/${projectName}/service/${serviceName}`);
  });

  it('should load and remove service param, if it does not exist', () => {
    const selectedDeployment = ServiceStateResponse[0].deploymentInformation[2];
    const serviceName = 'notExisting';
    const location = TestBed.inject(Location);
    const locationSpy = jest.spyOn(location, 'go');
    paramsSubject.next(
      convertToParamMap({
        projectName,
        serviceName,
        stage: 'staging',
        shkeptncontext: selectedDeployment.keptnContext,
      })
    );
    httpMock.expectOne(`./api/project/${projectName}/serviceStates`).flush(AppUtils.copyObject(ServiceStateResponse));
    expect(component.selectedDeployment).toBeUndefined();
    expect(locationSpy).toHaveBeenCalledWith(`/project/${projectName}/service`);
  });

  it('should only update remediations if deployment is finished', () => {
    const selectedDeployment = ServiceStateResponse[0].deploymentInformation[2];
    const deployment = Deployment.fromJSON(ServiceDeploymentMock);
    const serviceName = 'carts';
    const updateRemediationSpy = jest.spyOn(deployment, 'updateRemediations');

    component.deploymentSelected(
      {
        deploymentInformation: {
          ...selectedDeployment,
          deployment,
        },
        stage: 'production',
      },
      'sockshop'
    );

    httpMock.expectOne(`./api/project/${projectName}/service/${serviceName}/openRemediations`).flush({ stages: [] });
    fixture.detectChanges();
    expect(updateRemediationSpy).toHaveBeenCalled();
  });

  it('should not update remediations if hasRemediations is false', () => {
    const selectedDeployment = ServiceStateResponseWithoutRemediationsMock[0].deploymentInformation[0];
    const deployment = Deployment.fromJSON(ServiceDeploymentMock);
    const updateRemediationSpy = jest.spyOn(TestBed.inject(DataService), 'getOpenRemediationsOfService');

    component.deploymentSelected(
      {
        deploymentInformation: {
          ...selectedDeployment,
          deployment,
        },
        stage: 'production',
      },
      'sockshop'
    );
    fixture.detectChanges();
    expect(updateRemediationSpy).not.toHaveBeenCalled();
  });

  it('should update service states', () => {
    const serviceStateUpdateSpy = jest.spyOn(ServiceState, 'update');
    component.serviceStates = ServiceStateResponse.map((st) => ServiceState.fromJSON(st));
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    component.updateServiceStates(ServiceStateResponse.map((st) => ServiceState.fromJSON(st)));
    expect(serviceStateUpdateSpy).toHaveBeenCalled();
  });

  it('should update deployment', () => {
    const selectedDeployment = ServiceStateResponse[0].deploymentInformation[2];
    const deployment = Deployment.fromJSON(ServiceDeploymentWithApprovalMock);
    const updateDeploymentSpy = jest.spyOn(deployment, 'update');

    component.selectedDeployment = {
      deploymentInformation: {
        ...selectedDeployment,
        deployment,
      },
      stage: 'production',
    };

    component.deploymentSelected(
      {
        deploymentInformation: {
          ...selectedDeployment,
          deployment: deployment,
        },
        stage: 'production',
      },
      'sockshop'
    );

    httpMock
      .expectOne(
        `./api/project/${projectName}/deployment/${
          selectedDeployment.keptnContext
        }?fromTime=2021-10-13T10:54:43.315Z&includeRemediations=${selectedDeployment.stages.some(
          (st) => st.hasOpenRemediations
        )}`
      )
      .flush(AppUtils.copyObject(ServiceDeploymentWithApprovalMock));
    expect(updateDeploymentSpy).toHaveBeenCalled();
  });

  function loadDeployment(keptnContext: string, includeRemediations: boolean): void {
    httpMock.expectOne(`./api/project/${projectName}/serviceStates`).flush(AppUtils.copyObject(ServiceStateResponse));
    httpMock
      .expectOne(`./api/project/${projectName}/deployment/${keptnContext}?includeRemediations=${includeRemediations}`)
      .flush(AppUtils.copyObject(ServiceDeploymentMock));
  }
});
