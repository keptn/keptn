import { ComponentFixture, TestBed } from '@angular/core/testing';
import { DashboardLegacyComponent } from './dashboard-legacy.component';
import { DataService } from '../_services/data.service';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ApiService } from '../_services/api.service';
import { ApiServiceMock } from '../_services/api.service.mock';
import { finalize, skip, take } from 'rxjs/operators';
import { ProjectSequences } from '../_components/ktb-project-list/ktb-project-list.component';
import { POLLING_INTERVAL_MILLIS } from '../_utils/app.utils';
import { RouterTestingModule } from '@angular/router/testing';
import { Navigation, Router } from '@angular/router';

describe('DashboardLegacyComponent', () => {
  let component: DashboardLegacyComponent;
  let fixture: ComponentFixture<DashboardLegacyComponent>;
  let dataService: DataService;
  let router: Router;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HttpClientTestingModule, RouterTestingModule],
      providers: [
        { provide: ApiService, useClass: ApiServiceMock },
        { provide: POLLING_INTERVAL_MILLIS, useValue: 0 },
      ],
      declarations: [DashboardLegacyComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(DashboardLegacyComponent);
    component = fixture.componentInstance;
    dataService = TestBed.inject(DataService);
    router = TestBed.inject(Router);
  });

  it('should create', () => {
    // given, when, then
    expect(component).toBeTruthy();
  });

  it('should not load projects when navigated immediately to it', () => {
    // given
    const loadProjectSpy = jest.spyOn(dataService, 'loadProjects');

    // when, then
    expect(loadProjectSpy).not.toHaveBeenCalled();
  });

  it('should load projects if navigated from any other page', () => {
    // given
    const loadProjectSpy = jest.spyOn(dataService, 'loadProjects');
    jest.spyOn(router, 'getCurrentNavigation').mockReturnValue({ previousNavigation: {} } as Navigation);

    // when
    new DashboardLegacyComponent(router, dataService, 0);

    // then
    expect(loadProjectSpy).toHaveBeenCalled();
  });

  it('should load sequences one after another', (done) => {
    // given
    dataService.loadKeptnInfo();
    component.refreshProjects();
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
    component.refreshProjects();

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
});
