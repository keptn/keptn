import { TestBed } from '@angular/core/testing';
import { DataService } from './data.service';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../app.module';
import { APIService } from './api.service';
import { TriggerSequenceData } from '../_models/trigger-sequence';
import moment from 'moment';

describe('DataService', () => {
  let dataService: DataService;
  let apiService: APIService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
    });
    dataService = TestBed.inject(DataService);
    apiService = TestBed.inject(APIService);
  });

  it('should be an instance', () => {
    expect(dataService).toBeTruthy();
  });

  it('should trigger a delivery', () => {
    // given
    const spy = jest.spyOn(apiService, 'triggerSequence');
    const data: TriggerSequenceData = {
      project: 'podtato-head',
      stage: 'hardening',
      service: 'helloservice',
      configurationChange: {
        values: {
          image: 'docker.io/keptn:v0.1.2',
        },
      },
    };

    // when
    dataService.triggerDelivery(data);

    // then
    expect(spy).toHaveBeenCalledWith('sh.keptn.event.hardening.delivery.triggered', data);
  });

  it('should trigger an evaluation', () => {
    // given
    const spy = jest.spyOn(apiService, 'triggerSequence');
    const date = moment().toISOString();
    const data: TriggerSequenceData = {
      project: 'podtato-head',
      stage: 'hardening',
      service: 'helloservice',
      evaluation: {
        timeframe: '1h15m',
        start: date,
      },
    };

    // when
    dataService.triggerEvaluation(data);

    // then
    expect(spy).toHaveBeenCalledWith('sh.keptn.event.hardening.evaluation.triggered', data);
  });

  it('should trigger a custom sequence', () => {
    // given
    const spy = jest.spyOn(apiService, 'triggerSequence');
    const data: TriggerSequenceData = {
      project: 'podtato-head',
      stage: 'hardening',
      service: 'helloservice',
      labels: {
        key1: 'val1',
      },
    };

    // when
    dataService.triggerCustomSequence(data, 'testsequence');

    // then
    expect(spy).toHaveBeenCalledWith('sh.keptn.event.hardening.testsequence.triggered', data);
  });
});
