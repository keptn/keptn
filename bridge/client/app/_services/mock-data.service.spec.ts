import { TestBed } from '@angular/core/testing';

import { MockDataService } from './mock-data.service';

describe('MockDataService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should return 1 project', () => {
    const service: MockDataService = TestBed.get(MockDataService);
    service.loadProjects();
    service.projects.subscribe(projects => {
      expect(projects.length).toBe(1);
    });
  });
});
