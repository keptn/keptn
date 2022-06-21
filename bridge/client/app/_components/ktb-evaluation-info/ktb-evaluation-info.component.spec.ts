import { ComponentFixture, fakeAsync, TestBed, tick } from '@angular/core/testing';
import { KtbEvaluationInfoComponent } from './ktb-evaluation-info.component';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { EventTypes } from '../../../../shared/interfaces/event-types';
import { KeptnService } from '../../../../shared/models/keptn-service';
import { Trace } from '../../_models/trace';
import { ResultTypes } from '../../../../shared/models/result-types';
import { KtbEvaluationInfoModule } from './ktb-evaluation-info.module';

const evaluationTrace = Trace.fromJSON({
  id: 'myID2',
  data: {
    project: 'sockshop',
    service: 'carts',
    stage: 'dev',
  },
  type: EventTypes.EVALUATION_FINISHED,
  source: KeptnService.LIGHTHOUSE_SERVICE,
  shkeptncontext: 'myContext2',
});

describe('KtbEvaluationDetailsComponent', () => {
  let component: KtbEvaluationInfoComponent;
  let fixture: ComponentFixture<KtbEvaluationInfoComponent>;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [KtbEvaluationInfoModule, HttpClientTestingModule],
      providers: [],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbEvaluationInfoComponent);
    component = fixture.componentInstance;
    httpMock = TestBed.inject(HttpTestingController);
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should should load history if enabled and trace is provided', fakeAsync(() => {
    const project = 'sockshop';
    const stage = 'dev';
    const service = 'carts';
    const limit = 6;

    expect(component.evaluationsLoaded).toBe(false);
    component.evaluationInfo = {
      trace: Trace.fromJSON({
        id: 'myID',
        data: {
          project: project,
          service: service,
          stage: stage,
        },
        type: EventTypes.EVALUATION_FINISHED,
        source: KeptnService.LIGHTHOUSE_SERVICE,
        shkeptncontext: 'myContext',
      }),
      showHistory: true,
      data: {
        project: project,
        service: service,
        stage: stage,
      },
    };
    tick();
    httpMock
      .expectOne(
        `./api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?filter=data.project:${project}%20AND%20data.service:${service}%20AND%20data.stage:${stage}%20AND%20source:lighthouse-service&excludeInvalidated=true&limit=${limit}`
      )
      .flush({
        events: [evaluationTrace],
      });
    expect(component.evaluationsLoaded).toBe(true);
    expect(component.evaluation?.data.evaluationHistory?.length).toBe(1);
    expect(component.evaluationHistory).toEqual([evaluationTrace]);
    component.ngOnDestroy();
  }));

  it('should should load history if enabled and trace is not provided', fakeAsync(() => {
    const project = 'sockshop';
    const stage = 'dev';
    const service = 'carts';
    const limit = 5;

    expect(component.evaluationsLoaded).toBe(false);
    component.evaluationInfo = {
      trace: undefined,
      showHistory: true,
      data: {
        project: project,
        service: service,
        stage: stage,
      },
    };
    tick();
    httpMock
      .expectOne(
        `./api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?filter=data.project:${project}%20AND%20data.service:${service}%20AND%20data.stage:${stage}%20AND%20source:lighthouse-service&excludeInvalidated=true&limit=${limit}`
      )
      .flush({
        events: [evaluationTrace],
      });
    expect(component.evaluationsLoaded).toBe(true);
    expect(component.evaluationHistory).toEqual([evaluationTrace]);
    component.ngOnDestroy();
  }));

  it('should should show current evaluation in history', fakeAsync(() => {
    const project = 'sockshop';
    const stage = 'dev';
    const service = 'carts';
    const limit = 6;

    expect(component.evaluationsLoaded).toBe(false);
    component.evaluationInfo = {
      trace: Trace.fromJSON({
        id: 'myID',
        data: {
          project: project,
          service: service,
          stage: stage,
        },
        type: EventTypes.EVALUATION_FINISHED,
        source: KeptnService.LIGHTHOUSE_SERVICE,
        shkeptncontext: 'myContext',
      }),
      showHistory: true,
      data: {
        project: project,
        service: service,
        stage: stage,
      },
    };
    tick();
    httpMock
      .expectOne(
        `./api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?filter=data.project:${project}%20AND%20data.service:${service}%20AND%20data.stage:${stage}%20AND%20source:lighthouse-service&excludeInvalidated=true&limit=${limit}`
      )
      .flush({
        events: [
          {
            ...evaluationTrace,
            id: 'myID',
          },
        ],
      });
    expect(component.evaluationsLoaded).toBe(true);
    expect(component.evaluationHistory).toEqual([]);
    component.ngOnDestroy();
  }));

  it('should correctly set warning status with evaluationResult', () => {
    component.evaluationResult = {
      result: ResultTypes.WARNING,
      score: 0,
    };
    expect(component.isError).toBe(false);
    expect(component.isWarning).toBe(true);
    expect(component.isSuccess).toBe(false);
  });

  it('should correctly set error status with evaluationResult', () => {
    component.evaluationResult = {
      result: ResultTypes.FAILED,
      score: 0,
    };
    expect(component.isError).toBe(true);
    expect(component.isWarning).toBe(false);
    expect(component.isSuccess).toBe(false);
  });

  it('should correctly set success status with evaluationResult', () => {
    component.evaluationResult = {
      result: ResultTypes.PASSED,
      score: 0,
    };
    expect(component.isError).toBe(false);
    expect(component.isWarning).toBe(false);
    expect(component.isSuccess).toBe(true);
  });
});
