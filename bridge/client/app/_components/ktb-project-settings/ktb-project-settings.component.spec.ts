import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { BehaviorSubject, of } from 'rxjs';
import { PendingChangesGuard } from '../../_guards/pending-changes.guard';
import { IGitData } from '../../_interfaces/git-upstream';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { DataService } from '../../_services/data.service';
import { TestUtils } from '../../_utils/test.utils';
import { KtbProjectSettingsComponent } from './ktb-project-settings.component';
import { KtbProjectSettingsModule } from './ktb-project-settings.module';

describe('KtbProjectSettingsComponent', () => {
  let component: KtbProjectSettingsComponent;
  let fixture: ComponentFixture<KtbProjectSettingsComponent>;
  let dataService: DataService;
  let routeParamsSubject: BehaviorSubject<{ projectName: string } | Record<string, never>>;

  beforeEach(async () => {
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
            queryParams: of({}),
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbProjectSettingsComponent);
    component = fixture.componentInstance;
    dataService = fixture.debugElement.injector.get(DataService);
    dataService.loadKeptnInfo();
    dataService.loadProjects().subscribe();

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

  it('should call DataService.setGitUpstreamUrl on setGitUpstream', () => {
    // given
    const gitData: IGitData = {
      gitRemoteURL: 'https://test.git',
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
    component.isCreateMode = false;
    component.projectName = 'sockshop';
    fixture.detectChanges();

    // when
    const router = TestBed.inject(Router);
    const routeSpy = jest.spyOn(router, 'navigate');
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

  it('should not allow navigation for unsaved changes', () => {
    // given

    // when
    component.updateGitData({
      gitUser: 'someUser',
      gitRemoteURL: 'someUri',
      gitToken: 'someToken',
      gitFormValid: true,
    });
    fixture.detectChanges();

    // then
    expect(component.canDeactivate()).not.toEqual(true);
  });

  it('should show a dialog when showNotification is called', () => {
    // given
    component.updateGitData({
      gitUser: 'someUser',
      gitRemoteURL: 'someUri',
      gitToken: 'someToken',
      gitFormValid: true,
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
    component.updateGitData({
      gitUser: 'someUser',
      gitRemoteURL: 'someUri',
      gitToken: 'someToken',
      gitFormValid: true,
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
      https: {
        gitRemoteURL: 'https://myurl.git',
        gitToken: '',
        gitProxyInsecure: false,
        gitProxyPassword: '',
        gitProxyScheme: 'https',
        gitProxyUrl: '',
        gitProxyUser: '',
        gitUser: 'myUser',
      },
    });
    expect(component.isProjectFormTouched).toBe(true);

    // when
    await component.createProject();

    // then
    expect(createExtendedSpy).toHaveBeenCalledWith('myProject', btoa('test content'), {
      gitRemoteURL: 'https://myurl.git',
      gitToken: '',
      gitProxyInsecure: false,
      gitProxyPassword: '',
      gitProxyScheme: 'https',
      gitProxyUrl: '',
      gitProxyUser: '',
      gitUser: 'myUser',
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
      ssh: {
        gitPrivateKeyPass: 'myPrivateKeyPass',
        gitPrivateKey: 'myPrivateKey',
        gitRemoteURL: 'https://my-git-url.com',
      },
    });
    expect(component.isProjectFormTouched).toBe(true);

    // when
    await component.createProject();

    // then
    expect(createExtendedSpy).toHaveBeenCalledWith('myProject', btoa('test content'), {
      gitPrivateKeyPass: 'myPrivateKeyPass',
      gitPrivateKey: 'myPrivateKey',
      gitRemoteURL: 'https://my-git-url.com',
    });
  });

  it('should reset input for extended git data to default', () => {
    // given
    TestUtils.enableResourceService();
    // is a reference and may be modified by child components
    component.gitInputDataExtended = {
      ssh: {
        gitPrivateKeyPass: 'myPrivateKeyPass',
        gitPrivateKey: 'myPrivateKey',
        gitRemoteURL: 'https://my-git-url.com',
      },
    };

    // when
    component.reset();

    // then
    expect(component.gitInputDataExtended).toEqual({
      https: {
        gitProxyInsecure: false,
        gitProxyPassword: '',
        gitProxyScheme: 'https',
        gitProxyUrl: '',
        gitProxyUser: '',
        gitRemoteURL: 'https://github.com/Kirdock/keptn-dynatrace',
        gitToken: '',
        gitUser: 'Kirdock',
      },
    });
  });

  it('should update git upstream without user', () => {
    // given
    const updateSpy = jest.spyOn(dataService, 'setGitUpstreamUrl');
    component.projectName = 'sockshop';
    component.gitData = {
      gitRemoteURL: 'https://github.com/Kirdock/keptn-dynatrace',
      gitToken: 'myGitToken',
    };

    // when
    component.setGitUpstream();

    // then
    expect(updateSpy).toHaveBeenCalledWith(
      'sockshop',
      'https://github.com/Kirdock/keptn-dynatrace',
      'myGitToken',
      undefined
    );
  });
});
