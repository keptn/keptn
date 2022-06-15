import { TestBed } from '@angular/core/testing';
import { DashboardLegacyComponent } from './dashboard-legacy.component';
import { DataService } from '../_services/data.service';
import { of } from 'rxjs';
import { KeptnInfo } from '../_models/keptn-info';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('DashboardComponent', () => {
  let component: DashboardLegacyComponent;
  let dataService: DataService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [],
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

  function createComponent(): void {
    component = new DashboardLegacyComponent(dataService);
  }
});
