import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, convertToParamMap, Router } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { BehaviorSubject, of, throwError } from 'rxjs';
import { PendingChangesGuard } from '../../_guards/pending-changes.guard';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { DataService } from '../../_services/data.service';
import { TestUtils } from '../../_utils/test.utils';
import { KtbProjectSettingsComponent } from './ktb-project-settings.component';
import { KtbProjectSettingsModule } from './ktb-project-settings.module';
import { getDefaultSshData } from './ktb-project-settings-git-extended/ktb-project-settings-git-extended.component.spec';
import { EventService } from '../../_services/event.service';
import { DeleteResult, DeleteType } from '../../_interfaces/delete';
import { NotificationsService } from '../../_services/notifications.service';
import { NotificationType } from '../../_models/notification';
import { KtbProjectCreateMessageComponent } from './ktb-project-create-message/ktb-project-create-message.component';

describe('KtbProjectSettingsComponent', () => {
  let component: KtbProjectSettingsComponent;
  let fixture: ComponentFixture<KtbProjectSettingsComponent>;
  let dataService: DataService;
  let routeParamsSubject: BehaviorSubject<{ projectName: string } | Record<string, never>>;
  let queryParamsSubject: BehaviorSubject<{ created?: boolean }>;

  beforeEach(async () => {
    queryParamsSubject = new BehaviorSubject<{ created?: boolean }>({});
    routeParamsSubject = new BehaviorSubject<{ projectName: string } | Record<string, never>>({
      projectName: 'sockshop',
    });
    await TestBed.configureTestingModule({
      imports: [
        KtbProjectSettingsModule,
        HttpClientTestingModule,
        RouterTestingModule.withRoutes([
          { path: 'dashboard', component: KtbProjectSettingsComponent },
          { path: 'project/:projectName/settings/project', component: KtbProjectSettingsComponent },
        ]),
      ],
      providers: [
        PendingChangesGuard,
        { provide: ApiService, useClass: ApiServiceMock },
        {
          provide: ActivatedRoute,
          useValue: {
            params: routeParamsSubject.asObservable(),
            queryParams: queryParamsSubject.asObservable(),
            queryParamMap: of(convertToParamMap({})),
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbProjectSettingsComponent);
    component = fixture.componentInstance;
    dataService = fixture.debugElement.injector.get(DataService);
    dataService.loadKeptnInfo();
    dataService.loadProjects();

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

  it('should have create mode enabled when projectName param is not set', () => {
    // given
    routeParamsSubject.next({});
    fixture.detectChanges();

    // then
    expect(component.isCreateMode).toBe(true);
  });

  it('should have a validation error if project name already exists in projects', async () => {
    // given
    routeParamsSubject.next({});
    fixture.detectChanges();

    // when
    component.projectNameControl.setValue('sockshop');
    component.projectNameControl.updateValueAndValidity();
    fixture.detectChanges();

    // then
    expect(component.projectNameControl.hasError('duplicate')).toBe(true);
  });

  it('should navigate to created project', async () => {
    // given
    routeParamsSubject.next({});
    component.gitDataExtended = null;
    component.projectNameControl.setValue('sockshop');
    component.shipyardFile = new File(['test content'], 'test1.yaml');
    fixture.detectChanges();

    // when
    const router = TestBed.inject(Router);
    const routeSpy = jest.spyOn(router, 'navigate');
    await component.createProject();

    // then
    expect(routeSpy).toHaveBeenCalled();
  });

  it('should have create mode disabled when projectName param is set', () => {
    routeParamsSubject.next({ projectName: 'sockshop' });
    expect(component.isCreateMode).toBe(false);
  });

  it('should set project name to projectName retrieved by route', () => {
    routeParamsSubject.next({ projectName: 'sockshop' });
    expect(component.projectName).toEqual('sockshop');
  });

  it('should have an pattern validation error when project name does not match: first letter lowercase, only lowercase, numbers and hyphens allowed', () => {
    // given
    routeParamsSubject.next({});
    component.isCreateMode = true;
    fixture.detectChanges();

    component.projectNameControl.setValue('Sockshop');
    expect(component.projectNameControl.hasError('pattern')).toBe(true);

    component.projectNameControl.setValue('1ockshop');
    expect(component.projectNameControl.hasError('pattern')).toBe(true);

    component.projectNameControl.setValue('-ockshop');
    expect(component.projectNameControl.hasError('pattern')).toBe(true);

    component.projectNameControl.setValue('$ockshop');
    expect(component.projectNameControl.hasError('pattern')).toBe(true);

    component.projectNameControl.setValue('soCkshop');
    expect(component.projectNameControl.hasError('pattern')).toBe(true);

    component.projectNameControl.setValue('sock_shop');
    expect(component.projectNameControl.hasError('pattern')).toBe(true);

    component.projectNameControl.setValue('sockshop-1');
    expect(component.projectNameControl.errors).toBeNull();
  });

  it('should delete a project and navigate to dashboard', () => {
    // given
    const eventService = TestBed.inject(EventService);
    const router = TestBed.inject(Router);
    const routeSpy = jest.spyOn(router, 'navigate');
    component.isCreateMode = false;
    component.projectName = 'sockshop';
    fixture.detectChanges();

    // when
    eventService.deletionTriggeredEvent.next({ type: DeleteType.PROJECT, name: 'sockshop' });

    // then
    expect(routeSpy).toHaveBeenCalled();
  });

  it('should try to delete a project and throw an error', () => {
    // given
    const eventService = TestBed.inject(EventService);
    const progressSpy = jest.spyOn(eventService.deletionProgressEvent, 'next');
    jest.spyOn(dataService, 'deleteProject').mockReturnValue(throwError(() => new Error('my error')));
    component.isCreateMode = false;
    component.projectName = 'sockshop';
    fixture.detectChanges();

    // when
    eventService.deletionTriggeredEvent.next({ type: DeleteType.PROJECT, name: 'sockshop' });

    // then
    expect(progressSpy).toHaveBeenCalledWith({
      error: 'Project could not be deleted: my error',
      isInProgress: false,
      result: DeleteResult.ERROR,
    });
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

  it('should not allow navigation for unsaved changes', () => {
    // given

    // when
    component.updateGitDataExtended({
      user: 'someUser',
      remoteURL: 'someUri',
      https: {
        token: 'someToken',
        insecureSkipTLS: false,
      },
    });
    fixture.detectChanges();

    // then
    expect(component.canDeactivate()).not.toEqual(true);
  });

  it('should show a dialog when showNotification is called', () => {
    // given
    component.updateGitDataExtended({
      user: 'someUser',
      remoteURL: 'someUri',
      https: {
        token: 'someToken',
        insecureSkipTLS: false,
      },
    });
    fixture.detectChanges();

    // when
    component.showNotification();
    fixture.detectChanges();

    // then
    const notification = document.getElementsByTagName('dt-confirmation-dialog-state');
    expect(notification.length).toEqual(1);

    // We have to reset the state, as the dt-confirmation-dialog component has some pending timer open
    // and the test will not complete
    component.hideNotification();
    fixture.detectChanges();
  });

  it('should not show dialog when the notification was closed', () => {
    // given
    component.updateGitDataExtended({
      user: 'someUser',
      remoteURL: 'someUri',
      https: {
        token: 'someToken',
        insecureSkipTLS: false,
      },
    });
    fixture.detectChanges();

    // given
    component.showNotification();
    fixture.detectChanges();

    // when
    component.hideNotification();
    fixture.detectChanges();

    // then
    const notification = document.getElementsByTagName('dt-confirmation-dialog-state')[0];
    // It still exists in the dom but is hidden - so we test for aria-hidden
    expect(notification.getAttribute('aria-hidden')).toEqual('true');
  });

  it('should update extended Git HTTPS data', async () => {
    // given
    const apiService = TestBed.inject(ApiService);
    const createExtendedSpy = jest.spyOn(apiService, 'createProjectExtended');
    TestUtils.enableResourceService();
    expect(component.resourceServiceEnabled).toBe(true);
    expect(component.isProjectFormTouched).toBe(false);
    component.shipyardFile = new File(['test content'], 'test1.yaml');
    component.projectNameControl.setValue('myProject');

    // when
    component.updateGitDataExtended({
      remoteURL: 'https://myurl.git',
      user: 'myUser',
      https: {
        token: '',
        insecureSkipTLS: false,
        proxy: {
          password: '',
          scheme: 'https',
          url: 'myUrl:5000',
          user: '',
        },
      },
    });
    expect(component.isProjectFormTouched).toBe(true);

    // when
    await component.createProject();

    // then
    expect(createExtendedSpy).toHaveBeenCalledWith('myProject', btoa('test content'), {
      remoteURL: 'https://myurl.git',
      user: 'myUser',
      https: {
        token: '',
        insecureSkipTLS: false,
        proxy: {
          password: '',
          scheme: 'https',
          url: 'myUrl:5000',
          user: '',
        },
      },
    });
  });

  it('should update extended Git SSH data', async () => {
    // given
    const apiService = TestBed.inject(ApiService);
    const createExtendedSpy = jest.spyOn(apiService, 'createProjectExtended');
    TestUtils.enableResourceService();
    expect(component.resourceServiceEnabled).toBe(true);
    expect(component.isProjectFormTouched).toBe(false);
    component.shipyardFile = new File(['test content'], 'test1.yaml');
    component.projectNameControl.setValue('myProject');

    // when
    component.updateGitDataExtended({
      remoteURL: 'https://my-git-url.com',
      ssh: {
        privateKeyPass: 'myPrivateKeyPass',
        privateKey: 'myPrivateKey',
      },
    });
    expect(component.isProjectFormTouched).toBe(true);

    // when
    await component.createProject();

    // then
    expect(createExtendedSpy).toHaveBeenCalledWith('myProject', btoa('test content'), {
      remoteURL: 'https://my-git-url.com',
      ssh: {
        privateKeyPass: 'myPrivateKeyPass',
        privateKey: 'myPrivateKey',
      },
    });
  });

  it('should reset input for extended git data to default', () => {
    // given
    TestUtils.enableResourceService();
    // is a reference and may be modified by child components
    component.gitInputDataExtended = {
      remoteURL: 'https://my-git-url.com',
      ssh: {
        privateKeyPass: 'myPrivateKeyPass',
        privateKey: 'myPrivateKey',
      },
    };

    // when
    component.reset();

    // then
    expect(component.gitInputDataExtended).toEqual({
      remoteURL: 'https://github.com/Kirdock/keptn-dynatrace',
      user: 'Kirdock',
    });
  });

  it('should update git upstream without user', () => {
    // given
    const updateSpy = jest.spyOn(dataService, 'updateGitUpstream');
    component.projectName = 'sockshop';
    component.gitDataExtended = {
      remoteURL: 'https://github.com/Kirdock/keptn-dynatrace',
      https: {
        token: 'myGitToken',
        insecureSkipTLS: false,
      },
    };

    // when
    component.updateGitUpstream();

    // then
    expect(updateSpy).toHaveBeenCalledWith('sockshop', {
      remoteURL: 'https://github.com/Kirdock/keptn-dynatrace',
      https: {
        token: 'myGitToken',
        insecureSkipTLS: false,
      },
    });
  });

  it('should not update gitUpstream if data is not set', () => {
    // given
    const updateUpstreamSpy = jest.spyOn(dataService, 'updateGitUpstream');
    fixture.detectChanges();

    // when
    component.updateGitUpstream();

    // then
    expect(updateUpstreamSpy).not.toHaveBeenCalled();
  });

  it('should update gitUpstream', () => {
    // given
    fixture.detectChanges();
    const updateUpstreamSpy = jest.spyOn(dataService, 'updateGitUpstream');
    const data = getDefaultSshData();

    // when
    component.updateGitDataExtended(data);
    component.updateGitUpstream();

    // then
    expect(updateUpstreamSpy).toHaveBeenCalledWith('sockshop', data);
  });

  it('should update gitUpstream and set inProgress to false on error', () => {
    // given
    fixture.detectChanges();
    jest.spyOn(dataService, 'updateGitUpstream').mockReturnValue(throwError(() => of('error')));
    const updateUpstreamSpy = jest.spyOn(dataService, 'updateGitUpstream');
    const inProgressSpy = jest.fn();
    Object.defineProperty(component, 'isGitUpstreamInProgress', {
      get: jest.fn(() => true),
      set: inProgressSpy,
    });

    // when
    component.updateGitDataExtended(getDefaultSshData());
    component.updateGitUpstream();

    // then
    expect(inProgressSpy).toHaveBeenCalledWith(false);
    expect(updateUpstreamSpy).toHaveBeenCalled();
  });

  it('should not create project if shipyard is not set', async () => {
    // given
    component.gitDataExtended = null;
    component.projectNameControl.setValue('sockshop');

    // when
    const createSpy = jest.spyOn(dataService, 'createProjectExtended');
    await component.createProject();

    // then
    expect(createSpy).not.toHaveBeenCalled();
  });

  it('should not create project if shipyard is empty', async () => {
    // given
    component.gitDataExtended = null;
    component.projectNameControl.setValue('sockshop');
    component.shipyardFile = new File([''], 'test1.yaml');

    // when
    const createSpy = jest.spyOn(dataService, 'createProjectExtended');
    await component.createProject();

    // then
    expect(createSpy).not.toHaveBeenCalled();
  });

  it('should not create project if data is not set', async () => {
    // given
    component.gitDataExtended = undefined;
    component.projectNameControl.setValue('sockshop');
    component.shipyardFile = new File(['asfd'], 'test1.yaml');

    // when
    const createSpy = jest.spyOn(dataService, 'createProjectExtended');
    await component.createProject();

    // then
    expect(createSpy).not.toHaveBeenCalled();
  });

  it('should show an error if project creation failed', async () => {
    // given
    const notificationService = TestBed.inject(NotificationsService);
    const notificationSpy = jest.spyOn(notificationService, 'addNotification');
    jest.spyOn(dataService, 'createProjectExtended').mockReturnValue(throwError(() => new Error()));
    component.gitDataExtended = null;
    component.projectNameControl.setValue('sockshop');
    component.shipyardFile = new File(['asdf'], 'test1.yaml');

    // when
    await component.createProject();

    // then
    expect(notificationSpy).toHaveBeenCalledWith(
      NotificationType.ERROR,
      `The project could not be created: please, check the logs of resource-service.`
    );
  });

  it('should save unsaved changes on create', () => {
    // given
    component.unsavedDialogState = 'unsaved';
    component.isCreateMode = true;
    const createSpy = jest.spyOn(component, 'createProject');

    // when
    component.saveAll();

    // then
    expect(component.unsavedDialogState).toBeNull();
    expect(createSpy).toHaveBeenCalled();
  });

  it('should save unsaved changes on update', () => {
    // given
    component.unsavedDialogState = 'unsaved';
    const createSpy = jest.spyOn(component, 'updateGitUpstream');

    // when
    component.saveAll();

    // then
    expect(component.unsavedDialogState).toBeNull();
    expect(createSpy).toHaveBeenCalled();
  });

  it('should update shipyard file', () => {
    // given
    const shipyardFile = new File([], 'myfile.yaml');
    const formTouchedSpy = jest.spyOn(component, 'projectFormTouched');

    // when
    component.updateShipyardFile(shipyardFile);

    // then
    expect(component.shipyardFile).toEqual(new File([], 'myfile.yaml'));
    expect(formTouchedSpy).toHaveBeenCalled();
  });

  it('should show create message and remove param', () => {
    // given
    const notificationService = TestBed.inject(NotificationsService);
    const notificationSpy = jest.spyOn(notificationService, 'addNotification');
    const routeSpy = jest.spyOn(TestBed.inject(Router), 'navigate');

    // when
    queryParamsSubject.next({ created: true });

    expect(notificationSpy).toHaveBeenCalledWith(
      NotificationType.SUCCESS,
      '',
      {
        component: KtbProjectCreateMessageComponent,
        data: {
          projectName: 'sockshop',
          routerLink: '/project/sockshop/settings/services/create',
        },
      },
      10_000
    );
    expect(routeSpy).toHaveBeenCalledWith(['/', 'project', 'sockshop', 'settings', 'project']);
  });
});
