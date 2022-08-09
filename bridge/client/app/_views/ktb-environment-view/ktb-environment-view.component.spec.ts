import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbEnvironmentViewComponent } from './ktb-environment-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbEnvironmentViewModule } from './ktb-environment-view.module';
import { RouterTestingModule } from '@angular/router/testing';
import { POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { Location } from '@angular/common';
import { Stage } from '../../_models/stage';
import { ActivatedRoute, convertToParamMap, ParamMap } from '@angular/router';
import { BehaviorSubject } from 'rxjs';
import { DataService } from '../../_services/data.service';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('KtbEnvironmentViewComponent', () => {
  let component: KtbEnvironmentViewComponent;
  let fixture: ComponentFixture<KtbEnvironmentViewComponent>;
  let queryParamMap: BehaviorSubject<ParamMap>;
  let paramMap: BehaviorSubject<ParamMap>;
  let dataService: DataService;

  beforeEach(async () => {
    queryParamMap = new BehaviorSubject<ParamMap>(convertToParamMap({}));
    paramMap = new BehaviorSubject<ParamMap>(convertToParamMap({}));
    await TestBed.configureTestingModule({
      imports: [KtbEnvironmentViewModule, RouterTestingModule, BrowserAnimationsModule, HttpClientTestingModule],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            queryParamMap: queryParamMap.asObservable(),
            paramMap: paramMap.asObservable(),
          },
        },
        {
          provide: ApiService,
          useClass: ApiServiceMock,
        },
        { provide: POLLING_INTERVAL_MILLIS, useValue: 0 },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbEnvironmentViewComponent);
    component = fixture.componentInstance;
    dataService = TestBed.inject(DataService);
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should set stage and the right location', () => {
    // given
    const locationSpy = jest.spyOn(TestBed.inject(Location), 'go');

    // when
    component.setSelectedStageInfo('sockshop', { stage: { stageName: 'myStage' } as Stage, filterType: 'approval' });

    // then
    expect(locationSpy).toHaveBeenCalledWith('/project/sockshop/environment/stage/myStage?filterType=approval');
    expect(component.selectedStageInfo).toEqual({ stage: { stageName: 'myStage' } as Stage, filterType: 'approval' });
  });

  it('should set stage info and the right location without filterType', () => {
    // given
    const locationSpy = jest.spyOn(TestBed.inject(Location), 'go');

    // when
    component.setSelectedStageInfo('sockshop', { stage: { stageName: 'myStage' } as Stage, filterType: undefined });

    // then
    expect(locationSpy).toHaveBeenCalledWith('/project/sockshop/environment/stage/myStage');
    expect(component.selectedStageInfo).toEqual({ stage: { stageName: 'myStage' } as Stage, filterType: undefined });
  });

  it('should load the project', () => {
    // given
    const loadSpy = jest.spyOn(dataService, 'loadProject');

    // when
    paramMap.next(convertToParamMap({ projectName: 'sockshop' }));
    fixture.detectChanges();

    // then
    expect(loadSpy).toHaveBeenCalledWith('sockshop');
  });

  it('should set the selected stage through params', () => {
    // given, when
    paramMap.next(convertToParamMap({ projectName: 'sockshop', stageName: 'dev' }));

    // then
    expect(component.selectedStageInfo?.stage.stageName).toBe('dev');
    expect(component.selectedStageInfo?.filterType).toBeUndefined();
  });

  it('should set the selected stage with filterType through params', () => {
    // given, when
    queryParamMap.next(convertToParamMap({ filterType: 'approval' }));
    paramMap.next(convertToParamMap({ projectName: 'sockshop', stageName: 'dev' }));

    // then
    expect(component.selectedStageInfo?.stage.stageName).toBe('dev');
    expect(component.selectedStageInfo?.filterType).toBe('approval');
  });
});
