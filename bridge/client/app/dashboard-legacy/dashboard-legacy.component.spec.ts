import { TestBed } from '@angular/core/testing';
import { DashboardLegacyComponent } from './dashboard-legacy.component';
import { DataService } from '../_services/data.service';
import { of } from 'rxjs';
import { KeptnInfo } from '../_models/keptn-info';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ApiService } from '../_services/api.service';
import { ApiServiceMock } from '../_services/api.service.mock';
import { finalize, skip, take } from 'rxjs/operators';
import { ProjectSequences } from '../_components/ktb-project-list/ktb-project-list.component';

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

  it('should load sequences one after another', (done) => {
    // given
    dataService.loadKeptnInfo();
    createComponent();
    component.loadProjects();
    const emitTimes = 3;
    let emitted = 0;

    // when
    component.latestSequences$
      .pipe(take(emitTimes), finalize(done))
      .subscribe((projectSequences: ProjectSequences): void => {
        emitted++;
        // then
        // For every project the last sequences are loaded lazy time by time
        // So the record is growing by one each emit
        expect(Object.keys(projectSequences).length).toBe(emitted);
        expect(emitted).toBeLessThanOrEqual(emitTimes);
      });
  });

  it('should create a reacord with all sequences for projects', (done) => {
    // given
    dataService.loadKeptnInfo();
    createComponent();
    component.loadProjects();

    // when
    component.latestSequences$
      .pipe(skip(2), take(1), finalize(done))
      .subscribe((projectSequences: ProjectSequences): void => {
        // then
        expect(Object.keys(projectSequences)).toEqual(['sockshop', 'sockshop-approve', 'sockshop-carts-db']);
        expect(projectSequences.sockshop.length).toEqual(5);
        expect(projectSequences['sockshop-approve'].length).toEqual(5);
        expect(projectSequences['sockshop-carts-db'].length).toEqual(5);
      });
  });

  function createComponent(): void {
    component = new DashboardLegacyComponent(dataService, 0);
  }
});
