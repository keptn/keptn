import {TestBed} from '@angular/core/testing';
import {HttpClientTestingModule} from "@angular/common/http/testing";

import {AppModule} from "../app.module";
import {MockDataService} from './mock-data.service';

describe('MockDataService', () => {
  beforeEach(() => TestBed.configureTestingModule({
    declarations: [],
    imports: [AppModule, HttpClientTestingModule]
  }));

  it('should return 1 project', () => {
    const service: MockDataService = TestBed.get(MockDataService);
    service.loadProjects();
    service.projects.subscribe(projects => {
      expect(projects.length).toBe(2);
    });
  });
});
