import { TestBed } from '@angular/core/testing';
import { ApiService } from './api.service';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { AppModule } from '../app.module';
import { EventResult } from '../_interfaces/event-result';
import { HttpRequest } from '@angular/common/http';

describe('ApiService', () => {
  let apiService: ApiService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
    });
    apiService = TestBed.inject(ApiService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  it('should be an instance', () => {
    expect(apiService).toBeTruthy();
  });

  it('should fetch traces by ID', () => {
    apiService.getTracesByIds('myProject', ['id1', 'id2']).subscribe();
    httpMock.expectOne((req: HttpRequest<EventResult>) => {
      expect(req.url).toBe('./api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished');
      expect(req.params.get('filter')).toBe('data.project:myProject AND source:lighthouse-service AND id:id1,id2');
      expect(req.params.get('excludeInvalidated')).toBe('true');
      expect(req.params.get('limit')).toBe('2');
      expect(req.params.keys().length).toBe(3);
      return true;
    });
  });
});
