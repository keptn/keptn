import { ComponentFixture, fakeAsync, TestBed, waitForAsync } from '@angular/core/testing';

import { ProjectBoardComponent } from './project-board.component';
import { AppModule, INITIAL_DELAY_MILLIS } from '../app.module';
import { DataServiceMock } from '../_services/data.service.mock';
import { ActivatedRoute, convertToParamMap, ParamMap, UrlSegment } from '@angular/router';
import { BehaviorSubject, of } from 'rxjs';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { DataService } from '../_services/data.service';



describe('ProjectBoardComponent', () => {
  let component: ProjectBoardComponent;
  let fixture: ComponentFixture<ProjectBoardComponent>;
  let paramsSubject: BehaviorSubject<ParamMap>;

  function setupTraceTest(paramMap: ParamMap) {
    const route = TestBed.inject(ActivatedRoute);
    route.snapshot.url[0].path = 'trace';
    paramsSubject.next(paramMap);
    fixture = TestBed.createComponent(ProjectBoardComponent);
    component = fixture.componentInstance;
  }

  beforeEach(waitForAsync(() => {
    paramsSubject = new BehaviorSubject(convertToParamMap({projectName: 'sockshop'}));
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule
      ],
      providers: [
        {provide: DataService, useClass: DataServiceMock},
        {provide: INITIAL_DELAY_MILLIS, useValue: 0},
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: paramsSubject.asObservable(),
            snapshot: {url: [{path: 'project'}, {path: 'sockshop'}]},
            url: of([new UrlSegment('project', {}), new UrlSegment('sockshop', {})])
          }
        }
      ]
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(ProjectBoardComponent);
        component = fixture.componentInstance;

        fixture.detectChanges();
      });
  }));

  it('should have "project" error when project can not be found', (done) => {
    // given
    const projectName = 'wrong-project';

    // when
    paramsSubject.next(convertToParamMap({projectName}));
    fixture.detectChanges();

    // then
    component.error$.subscribe(err => {
      expect(err).toEqual('project');
      done();
    });
  });

  it('should have a "trace" error when a trace can not be found with wrong shkeptncontext', (done) => {
    // given
    setupTraceTest(convertToParamMap({shkeptncontext: 'asdf123asdf456789'}));

    // when
    fixture.detectChanges();

    // then
    component.error$.subscribe(err => {
      expect(err).toEqual('trace');
      done();
    });
  });

  it('should have a "trace" error when a trace can not be found with right shkeptncontext but wrong eventselector', (done) => {
    // given
    setupTraceTest(convertToParamMap({shkeptncontext: '0bbaaa6b-fd89-4def-ad2c-975beda970cf', eventselector: 'some-wrong-selector'}));

    // when
    fixture.detectChanges();

    // then
    component.error$.subscribe(err => {
      expect(err).toEqual('trace');
      done();
    });
  });

  it('should show a project doesn\'t exists message when error is project', () => {
    // given
    const projectName = 'wrong-project';

    // when
    paramsSubject.next(convertToParamMap({projectName}));
    fixture.detectChanges();

    // then
    const elem = fixture.nativeElement.querySelector('dt-empty-state-item-title');
    expect(elem.textContent).toEqual('Project doesn\'t exist');
  });

  it('should show a trace not found message when error is trace', () => {
    // given
    setupTraceTest(convertToParamMap({shkeptncontext: 'asdf123asdf456789'}));

    // when
    fixture.detectChanges();

    // then
    const elem = fixture.nativeElement.querySelector('dt-empty-state-item-title');
    expect(elem.textContent).toContain('Traces for asdf123asdf456789 not found');
  });

  it('should have an undefined error when no error occurred', (done) => {
    // when
    fixture.detectChanges();

    // then
    component.error$.subscribe(err => {
      expect(err).toBeUndefined();
      done();
    });
  });

  it('should have a hasProject set to true when no error occurred', (done) => {
    // when
    fixture.detectChanges();

    // then
    component.hasProject$.subscribe(hasProject => {
      expect(hasProject).toBeTrue();
      done();
    });
  });

  it('should have a hasProject set to false when project is not found occurred', (done) => {
    // given
    const projectName = 'wrong-project';

    // when
    paramsSubject.next(convertToParamMap({projectName}));
    fixture.detectChanges();

    // then
    component.hasProject$.subscribe(hasProject => {
      expect(hasProject).toBeFalse();
      done();
    });
  });

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});


