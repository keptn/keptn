import {TestBed} from '@angular/core/testing';
import {HttpClientTestingModule} from "@angular/common/http/testing";

import {AppModule} from "../app.module";
import {DataServiceMock} from './data.service.mock';

describe('MockDataService', () => {
  beforeEach(() => TestBed.configureTestingModule({
    declarations: [],
    imports: [
      AppModule,
      HttpClientTestingModule,
    ]
  }));

  it('should return 1 project', () => {
    const service: DataServiceMock = TestBed.get(DataServiceMock);
    service.loadProjects();
    service.projects.subscribe(projects => {
      expect(projects.length).toBe(2);
    });
  });
});
