import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbProjectSettingsComponent } from './ktb-project-settings.component';
import { DataService } from '../../_services/data.service';
import { BehaviorSubject, of } from 'rxjs';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { DataServiceMock } from '../../_services/data.service.mock';
import { ActivatedRoute, Router } from '@angular/router';

describe('KtbProjectSettingsComponent', () => {
  let component: KtbProjectSettingsComponent;
  let fixture: ComponentFixture<KtbProjectSettingsComponent>;
  const UNSAVED_DIALOG_STATE = 'unsaved';
  let dataService: DataService;
  const routeDataSubject = new BehaviorSubject<{ isCreateMode: boolean }>({ isCreateMode: false });

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
      providers: [
        { provide: DataService, useClass: DataServiceMock },
        {
          provide: ActivatedRoute,
          useValue: {
            params: of({ projectName: 'sockshop' }),
            data: routeDataSubject.asObservable(),
            queryParams: of({}),
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbProjectSettingsComponent);
    component = fixture.componentInstance;
    dataService = fixture.debugElement.injector.get(DataService);

    const notifications = document.getElementsByTagName('dt-confirmation-dialog-state');
    if (notifications.length > 0) {
      // eslint-disable-next-line @typescript-eslint/prefer-for-of
      for (let i = 0; i < notifications.length; i++) {
        notifications[i].remove();
      }
    }

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should have create mode enabled when routed to route data contains {isCreateMode: true}', () => {
    // given
    routeDataSubject.next({ isCreateMode: true });
    fixture.detectChanges();

    // then
    expect(component.isCreateMode).toBe(true);
  });

  it('should have a validation error if project name already exists in projects', async () => {
    // given
    component.isCreateMode = true;
    fixture.detectChanges();

    // when
    await dataService.loadProjects();
    component.projectNameControl.setValue('sockshop');

    // then
    expect(component.projectNameControl.hasError('duplicate')).toBe(true);
  });

  it('should navigate to created project', async () => {
    // given
    component.isCreateMode = true;
    component.projectName = 'sockshop';
    fixture.detectChanges();

    // when
    const router = TestBed.inject(Router);
    const routeSpy = jest.spyOn(router, 'navigate');
    await dataService.loadProjects();

    // then
    expect(routeSpy).toHaveBeenCalled();
  });

  it('should call DataService.setGitUpstreamUrl on setGitUpstream', () => {
    // given
    const gitData = {
      remoteURI: 'https://test.git',
      gitUser: 'username',
      gitToken: 'token',
    };
    component.projectName = 'sockshop';

    // when
    const spy = jest.spyOn(dataService, 'setGitUpstreamUrl');
    component.updateGitData(gitData);
    component.setGitUpstream();

    // then
    expect(spy).toHaveBeenCalled();
  });

  it('should have create mode disabled when routed to route data contains {isCreateMode: false}', () => {
    // given
    routeDataSubject.next({ isCreateMode: false });
    fixture.detectChanges();

    expect(component.isCreateMode).toBe(false);
  });

  it('should set project name to projectName retrieved by route', () => {
    expect(component.projectName).toEqual('sockshop');
  });

  it('should have an pattern validation error when project name does not match: first letter lowercase, only lowercase, numbers and hyphens allowed', () => {
    // given
    component.isCreateMode = true;
    fixture.detectChanges();

    component.projectNameControl.setValue('Sockshop');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.hasError('pattern'));

    component.projectNameControl.setValue('1ockshop');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.hasError('pattern'));

    component.projectNameControl.setValue('-ockshop');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.hasError('pattern'));

    component.projectNameControl.setValue('$ockshop');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.hasError('pattern'));

    component.projectNameControl.setValue('soCkshop');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.hasError('pattern'));

    component.projectNameControl.setValue('sock_shop');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.hasError('pattern'));

    component.projectNameControl.setValue('sockshop-1');
    component.projectNameForm.updateValueAndValidity();
    expect(component.projectNameForm.errors).toBeNull();
  });

  it('should delete a project and navigate to dashboard', () => {
    // given
    component.isCreateMode = false;
    component.projectName = 'sockshop';
    fixture.detectChanges();

    // when
    const router = TestBed.inject(Router);
    const routeSpy = jest.spyOn(router, 'navigate');

    dataService.loadProjects = jest.fn().mockImplementation(() => {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      dataService._projects.next(
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        dataService._projects.getValue().filter((project) => project.projectName !== 'sockshop')
      );
    });
    component.deleteProject('sockshop');

    // then
    expect(routeSpy).toHaveBeenCalled();
  });

  it('should not show a notification when the component is initialized', () => {
    // given
    component.isCreateMode = false;
    fixture.detectChanges();

    // then
    // Has to be retrieved by document, as it is not created at component level
    const notifications = document.getElementsByTagName('dt-confirmation-dialog-state');
    expect(component.unsavedDialogState).toBeNull();
    expect(notifications.length).toEqual(0);
  });

  it('should show a notification when "unsaved" is set', () => {
    // given
    component.isCreateMode = false;
    component.unsavedDialogState = UNSAVED_DIALOG_STATE;
    fixture.detectChanges();

    // then
    const notification = document.getElementsByTagName('dt-confirmation-dialog-state');
    expect(notification.length).toEqual(1);

    // We have to reset the state, as the dt-confirmation-dialog component has some pending timer open
    // and the test will not complete
    component.unsavedDialogState = null;
    fixture.detectChanges();
  });

  it('should show a notification for unsaved changes when git data is changed in update mode', () => {
    // given
    component.isCreateMode = true;
    fixture.detectChanges();

    // when
    component.updateGitData({ gitUser: 'someUser', remoteURI: 'someUri', gitToken: 'someToken', gitFormValid: true });
    fixture.detectChanges();

    // then
    expect(component.unsavedDialogState).toEqual(UNSAVED_DIALOG_STATE);
  });

  it('should not show a notification for unsaved changes when git data is changed in create mode', () => {
    // given
    component.isCreateMode = true;
    fixture.detectChanges();

    // when
    component.updateGitData({ gitUser: 'someUser', remoteURI: 'someUri', gitToken: 'someToken', gitFormValid: true });
    fixture.detectChanges();

    // then
    expect(component.unsavedDialogState).toEqual(UNSAVED_DIALOG_STATE);
  });

  it('should not show a notification when not all git data fields are set', () => {
    // given
    component.isCreateMode = false;
    fixture.detectChanges();

    // when
    component.updateGitData({ gitUser: 'someUser', remoteURI: 'someUri' });
    fixture.detectChanges();

    // then
    expect(component.unsavedDialogState).toBeNull();

    // when
    component.updateGitData({ gitUser: 'someUser', gitToken: 'someToken' });
    fixture.detectChanges();

    // then
    expect(component.unsavedDialogState).toBeNull();

    // when
    component.updateGitData({ remoteURI: 'someUri', gitToken: 'someToken' });
    fixture.detectChanges();

    // then
    expect(component.unsavedDialogState).toBeNull();
  });

  it('should not show a notification when the notification was closed', () => {
    // given
    component.isCreateMode = false;
    fixture.detectChanges();

    // given
    component.unsavedDialogState = UNSAVED_DIALOG_STATE;
    fixture.detectChanges();

    // when
    component.unsavedDialogState = null;
    fixture.detectChanges();

    // then
    const notification = document.getElementsByTagName('dt-confirmation-dialog-state')[0];
    // It still exists in the dom but is hidden - so we test for aria-hidden
    expect(notification.getAttribute('aria-hidden')).toEqual('true');
  });
});
