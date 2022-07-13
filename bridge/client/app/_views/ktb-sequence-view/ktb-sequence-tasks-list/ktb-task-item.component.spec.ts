import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbTaskItemComponent } from './ktb-task-item.component';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { Trace } from '../../../_models/trace';
import { RouterTestingModule } from '@angular/router/testing';
import { KtbSequenceViewModule } from '../ktb-sequence-view.module';

describe('KtbEventItemComponent', () => {
  let component: KtbTaskItemComponent;
  let fixture: ComponentFixture<KtbTaskItemComponent>;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbSequenceViewModule, HttpClientTestingModule, RouterTestingModule],
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
    // given, when
    component.isExpanded = true;
    component.task = getApprovalTrace();

    // then
    httpMock.expectOne('./api/controlPlane/v1/project/socksohp/stage/dev/service/carts').flush({
      deployedImage: 'docker.io/mongo:4.2.2',
    });
    expect(component.latestDeployment).toBe('docker.io/mongo:4.2.2');
  });

  it('should not set latestDeployment if task is not an approval', () => {
    // given, when
    component.isExpanded = true;
    component.task = Trace.fromJSON({});

    // then
    expect(component.latestDeployment).toBeUndefined();
  });

  it('should revert the latestDeployment if task is changed', () => {
    // given
    component.isExpanded = true;
    component.task = getApprovalTrace();

    httpMock.expectOne('./api/controlPlane/v1/project/socksohp/stage/dev/service/carts').flush({
      deployedImage: 'docker.io/mongo:4.2.2',
    });

    // when
    component.task = getApprovalTrace();

    // then
    httpMock.expectOne('./api/controlPlane/v1/project/socksohp/stage/dev/service/carts').flush({
      deployedImage: 'docker.io/mongo:4.2.3',
    });
    expect(component.latestDeployment).toBe('docker.io/mongo:4.2.3');
  });

  it('should not fetch the latestDeployment if it is not expanded', () => {
    // given, when
    component.task = getApprovalTrace();

    // then
    httpMock.expectNone('./api/controlPlane/v1/project/socksohp/stage/dev/service/carts');
  });

  it('should not fetch the latestDeployment if it is already fetched', () => {
    // given
    component.isExpanded = true;
    component.task = getApprovalTrace();

    httpMock.expectOne('./api/controlPlane/v1/project/socksohp/stage/dev/service/carts').flush({
      deployedImage: 'docker.io/mongo:4.2.2',
    });

    // when
    component.isExpanded = false;
    component.isExpanded = true;

    // then
    httpMock.expectNone('./api/controlPlane/v1/project/socksohp/stage/dev/service/carts');
  });

  it('isUrl should return true if given valid URL', () => {
    // given
    const url = 'https://keptn.sh';

    // when
    const isUrl = component.isUrl(url);

    // then
    expect(isUrl).toBeTruthy();
  });

  it('isUrl should return false if given invalid URL', () => {
    // given
    const url = 'keptn.sh';

    // when
    const isUrl = component.isUrl(url);

    // then
    expect(isUrl).toBeFalsy();
  });

  function getApprovalTrace(): Trace {
    return Trace.fromJSON({
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
  }
});
