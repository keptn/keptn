import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbProjectViewComponent } from './ktb-project-view.component';
import { ActivatedRoute, convertToParamMap, ParamMap, UrlSegment } from '@angular/router';
import { BehaviorSubject, of } from 'rxjs';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { KtbProjectViewCommonModule } from './ktb-project-view-common.module';
import { RouterTestingModule } from '@angular/router/testing';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('ProjectBoardComponent', () => {
  let component: KtbProjectViewComponent;
  let fixture: ComponentFixture<KtbProjectViewComponent>;
  let paramsSubject: BehaviorSubject<ParamMap>;

  beforeEach(async () => {
    paramsSubject = new BehaviorSubject(convertToParamMap({ projectName: 'sockshop' }));

    await TestBed.configureTestingModule({
      declarations: [],
      imports: [KtbProjectViewCommonModule, BrowserAnimationsModule, RouterTestingModule, HttpClientTestingModule],
      providers: [
        { provide: ApiService, useClass: ApiServiceMock },
        { provide: POLLING_INTERVAL_MILLIS, useValue: 0 },
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: paramsSubject.asObservable(),
            snapshot: { url: [{ path: 'project' }, { path: 'sockshop' }] },
            url: of([new UrlSegment('project', {}), new UrlSegment('sockshop', {})]),
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbProjectViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should have "project" error when project can not be found', (done) => {
    // given
    const projectName = 'wrong-project';

    // when
    paramsSubject.next(convertToParamMap({ projectName }));
    fixture.detectChanges();

    // then
    component.error$.subscribe((err) => {
      expect(err).toEqual('project');
      done();
    });
  });

  it("should show a project doesn't exists message when error is project", () => {
    // given
    const projectName = 'wrong-project';

    // when
    paramsSubject.next(convertToParamMap({ projectName }));
    fixture.detectChanges();

    // then
    const elem = fixture.nativeElement.querySelector('dt-empty-state-item-title');
    expect(elem.textContent).toEqual("Project doesn't exist");
  });

  it('should have an undefined error when no error occurred', (done) => {
    // when
    fixture.detectChanges();

    // then
    component.error$.subscribe((err) => {
      expect(err).toBeUndefined();
      done();
    });
  });

  it('should have a hasProject set to true when no error occurred', (done) => {
    // when
    fixture.detectChanges();

    // then
    component.hasProject$.subscribe((hasProject) => {
      expect(hasProject).toBe(true);
      done();
    });
  });

  it('should have a hasProject set to false when project is not found occurred', (done) => {
    // given
    const projectName = 'wrong-project';

    // when
    paramsSubject.next(convertToParamMap({ projectName }));
    fixture.detectChanges();

    // then
    component.hasProject$.subscribe((hasProject) => {
      expect(hasProject).toBe(false);
      done();
    });
  });

  it('should not show notification indicator on component creation', () => {
    expect(fixture.nativeElement.querySelector('.notification-indicator')).toBeFalsy();
  });

  it('should not show notification indicator after changing hasUnreadLogs from true to false', () => {
    component.hasUnreadLogs$ = of(true);
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('.notification-indicator')).toBeTruthy();

    component.hasUnreadLogs$ = of(false);
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('.notification-indicator')).toBeFalsy();
  });

  it('should show notification indicator if hasUnreadLogs is set to true', () => {
    component.hasUnreadLogs$ = of(true);
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('.notification-indicator')).toBeTruthy();
  });
});
