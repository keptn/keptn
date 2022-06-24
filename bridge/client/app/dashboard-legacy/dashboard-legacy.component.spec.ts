import { TestBed } from '@angular/core/testing';
import { DashboardLegacyComponent } from './dashboard-legacy.component';
import { DataService } from '../_services/data.service';
import { of } from 'rxjs';
import { KeptnInfo } from '../_models/keptn-info';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ApiService } from '../_services/api.service';
import { ApiServiceMock } from '../_services/api.service.mock';

describe('DashboardLegacyComponent', () => {
  let component: DashboardLegacyComponent;
  let dataService: DataService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [{ provide: ApiService, useClass: ApiServiceMock }],
    }).compileComponents();

    dataService = TestBed.inject(DataService);
  });

  it('should create', () => {
    // given, when
    createComponent();

    // then
    expect(component).toBeTruthy();
  });

  it('should not load projects if keptnInfo is not loaded', () => {
    // given
    const loadProjectSpy = jest.spyOn(dataService, 'loadProjects');

    // when
    createComponent();

    // then
    expect(loadProjectSpy).not.toHaveBeenCalled();
  });

  it('should load projects if keptnInfo is loaded', () => {
    // given
    const loadProjectSpy = jest.spyOn(dataService, 'loadProjects');
    jest.spyOn(dataService, 'keptnInfo', 'get').mockReturnValue(of({} as unknown as KeptnInfo));

    // when
    createComponent();

    // then
    expect(loadProjectSpy).toHaveBeenCalled();
  });

  it('should load a maximum of 5 sequences per project', (done) => {
    // given
    dataService.loadKeptnInfo();
    createComponent();
    component.loadProjects();

    // when
    component.latestSequences$.subscribe((projectSequences) => {
      if (Object.keys(projectSequences).length != 3) {
        return;
      }

      // then
      expect(Object.keys(projectSequences)).toEqual(['sockshop', 'sockshop-approve', 'sockshop-carts-db']);
      expect(projectSequences.sockshop.length).toEqual(5);
      done();
    });
  });

  function createComponent(): void {
    component = new DashboardLegacyComponent(dataService, 0);
  }
});
