import {ComponentFixture, TestBed} from '@angular/core/testing';

import {KtbSequenceStateListComponent} from './ktb-sequence-state-list.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";
import {ResultTypes} from "../../_models/result-types";
import {EvaluationResult} from "../../_models/evaluation-result";
import {Sequence} from "../../_models/sequence";

describe('KtbSequenceStateListComponent', () => {
  let component: KtbSequenceStateListComponent;
  let fixture: ComponentFixture<KtbSequenceStateListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [KtbSequenceStateListComponent],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSequenceStateListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set datasource', () => {
    // given
    component.sequenceStates = [{
      "name": "delivery",
      "service": "carts",
      "project": "sockshop",
      "time": "2021-07-20T08:36:11.208Z",
      "shkeptncontext": "d40a970f-7ffb-459e-a960-cabe7ca89c9c",
      "state": "triggered",
      "stages": [{
        "image": "docker.io/keptnexamples/carts:0.12.1",
        "latestEvaluation": {"result": ResultTypes.PASSED, "score": 0} as EvaluationResult,
        "latestEvent": {
          "type": "sh.keptn.event.dev.delivery.finished",
          "id": "2b2ef01f-4663-4588-9d52-1480fd67e249",
          "time": "2021-07-20T08:37:27.134Z"
        },
        "name": "dev"
      }, {
        "image": "docker.io/keptnexamples/carts:0.12.1",
        "latestEvaluation": {"result": ResultTypes.PASSED, "score": 0} as EvaluationResult,
        "latestEvent": {
          "type": "sh.keptn.event.release.started",
          "id": "73597657-14c7-4865-bc23-82064ea42532",
          "time": "2021-07-20T08:38:37.221Z"
        },
        "name": "staging"
      }]
    }];

    // when
    fixture.detectChanges();

    // then
    expect(component.dataSource.data.length).toEqual(4);
  });

  afterEach(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  });
});
