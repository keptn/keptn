import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbTaskItemComponent } from './ktb-task-item.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { Trace } from '../../_models/trace';
import { AppUtils } from '../../_utils/app.utils';

const approvalTrace = Trace.fromJSON({
  data: {
    approval: {
      pass: 'manual',
      warning: 'manual',
    },
    project: 'socksohp',
    service: 'carts',
    stage: 'dev',
  },
  gitcommitid: '05f241d1c1f022d65ab6af60bd2ef1272c1fbfbb',
  id: 'c0275a97-fc4f-4279-bc3e-c892482c8dbd',
  shkeptncontext: '51b96d8c-ff67-4a75-859b-124bcdbb04ff',
  shkeptnspecversion: '0.2.3',
  source: 'shipyard-controller',
  specversion: '1.0',
  time: '2022-05-25T12:50:52.560Z',
  type: 'sh.keptn.event.approval.triggered',
  finished: false,
  started: true,
  label: 'approval',
  traces: [
    {
      traces: [],
      data: {
        message: "Approval strategy for result '': manual",
        project: 'socksohp',
        service: 'helloservice',
        stage: 'dev',
        status: 'succeeded',
      },
      id: '1ade2f8f-8bcd-4afc-a89f-c5f5fc6eaeca',
      shkeptncontext: '51b96d8c-ff67-4a75-859b-124bcdbb04ff',
      shkeptnspecversion: '0.2.4',
      source: 'approval-service',
      specversion: '1.0',
      time: '2022-05-25T12:50:52.763Z',
      triggeredid: 'c0275a97-fc4f-4279-bc3e-c892482c8dbd',
      type: 'sh.keptn.event.approval.started',
      icon: 'unknown',
    },
  ],
});

describe('KtbEventItemComponent', () => {
  let component: KtbTaskItemComponent;
  let fixture: ComponentFixture<KtbTaskItemComponent>;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbTaskItemComponent);
    component = fixture.componentInstance;
    httpMock = TestBed.inject(HttpTestingController);
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should set latestDeployment if task is an approval', () => {
    component.isExpanded = true;
    component.task = approvalTrace;

    httpMock.expectOne('./api/controlPlane/v1/project/socksohp/stage/dev/service/carts').flush({
      deployedImage: 'docker.io/mongo:4.2.2',
    });
    expect(component.latestDeployment).toBe('docker.io/mongo:4.2.2');
  });

  it('should not set latestDeployment if task is not an approval', () => {
    component.isExpanded = true;
    component.task = Trace.fromJSON({});

    expect(component.latestDeployment).toBeUndefined();
  });

  it('should revert the latestDeployment if task is changed', () => {
    component.isExpanded = true;
    component.task = approvalTrace;

    httpMock.expectOne('./api/controlPlane/v1/project/socksohp/stage/dev/service/carts').flush({
      deployedImage: 'docker.io/mongo:4.2.2',
    });

    component.task = Trace.fromJSON(AppUtils.copyObject(approvalTrace));
    httpMock.expectOne('./api/controlPlane/v1/project/socksohp/stage/dev/service/carts').flush({
      deployedImage: 'docker.io/mongo:4.2.3',
    });
    expect(component.latestDeployment).toBe('docker.io/mongo:4.2.3');
  });

  it('should not fetch the latestDeployment if it is not expanded', () => {
    component.task = approvalTrace;
    httpMock.expectNone('./api/controlPlane/v1/project/socksohp/stage/dev/service/carts');
  });

  it('should not fetch the latestDeployment if it is already fetched', () => {
    component.isExpanded = true;
    component.task = approvalTrace;

    httpMock.expectOne('./api/controlPlane/v1/project/socksohp/stage/dev/service/carts').flush({
      deployedImage: 'docker.io/mongo:4.2.2',
    });

    component.isExpanded = false;
    component.isExpanded = true;
    httpMock.expectNone('./api/controlPlane/v1/project/socksohp/stage/dev/service/carts');
  });
});
