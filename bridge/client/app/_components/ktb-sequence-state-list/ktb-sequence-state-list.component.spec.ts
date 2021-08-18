import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing';
import { KtbSequenceStateListComponent } from './ktb-sequence-state-list.component';
import { AppModule } from '../../app.module';
import { ResultTypes } from '../../../../shared/models/result-types';
import { Sequence } from '../../_models/sequence';
import {INITIAL_DELAY_MILLIS} from '../../_utils/app.utils';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbSequenceStateListComponent', () => {
  let component: KtbSequenceStateListComponent;
  let fixture: ComponentFixture<KtbSequenceStateListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule
      ],
      providers: [
        {provide: INITIAL_DELAY_MILLIS, useValue: 0}
      ]
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbSequenceStateListComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set datasource', () => {
    // given
    component.sequenceStates = [Sequence.fromJSON({
      name: 'delivery',
      service: 'carts',
      project: 'sockshop',
      time: '2021-07-20T08:36:11.208Z',
      shkeptncontext: 'd40a970f-7ffb-459e-a960-cabe7ca89c9c',
      state: 'triggered',
      stages: [{
        image: 'docker.io/keptnexamples/carts:0.12.1',
        latestEvaluation: {result: ResultTypes.PASSED, score: 0},
        latestEvent: {
          type: 'sh.keptn.event.dev.delivery.finished',
          id: '2b2ef01f-4663-4588-9d52-1480fd67e249',
          time: '2021-07-20T08:37:27.134Z',
        },
        name: 'dev',
      }],
    })];

    // when
    fixture.detectChanges();

    // then
    expect(component.dataSource.data.length).toEqual(1);
  });

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
