import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbSettingsViewComponent } from './ktb-settings-view.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';
import { BehaviorSubject, of } from 'rxjs';
import { ActivatedRoute, Router } from '@angular/router';


describe('KtbSettingsViewComponent', () => {
  const UNSAVED_DIALOG_STATE = 'unsaved';
  let component: KtbSettingsViewComponent;
  let fixture: ComponentFixture<KtbSettingsViewComponent>;
  let dataService: DataService;
  const routeDataSubject = new BehaviorSubject<{ isCreateMode: boolean }>({isCreateMode: false});

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [KtbSettingsViewComponent],
      imports: [AppModule, HttpClientTestingModule],
      providers: [
        {provide: DataService, useClass: DataServiceMock},
        {
          provide: ActivatedRoute,
          useValue: {
            params: of({projectName: 'sockshop'}),
            data: routeDataSubject.asObservable(),
            queryParams: of({})
          }
        }
      ]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSettingsViewComponent);
    component = fixture.componentInstance;
    dataService = fixture.debugElement.injector.get(DataService);
    fixture.detectChanges();
  });

  it('should create settings view component', () => {
    expect(component).toBeTruthy();
  });

  it('should have create mode enabled when routed to route data contains {isCreateMode: true}', () => {
    // given
    routeDataSubject.next({isCreateMode: true});
    fixture.detectChanges();

    // then
    expect(component.isCreateMode).toBeTrue();
  });

  it('should have a validation error if project name already exists in projects', async () => {
    // given
    component.isCreateMode = true;
    fixture.detectChanges();

    // when
    await dataService.loadProjects();
    component.projectNameControl.setValue('sockshop');

    // then
    expect(component.projectNameControl.hasError('projectName')).toBeTrue();
  });

  it('should navigate to created project', async () => {
    // given
    component.isCreateMode = true;
    component.projectName = 'sockshop';

    // when
    const router = TestBed.inject(Router);
    const routeSpy = spyOn(router, 'navigate');
    await dataService.loadProjects();

    // then
    expect(routeSpy).toHaveBeenCalled();
  });

  it('should call DataService.setGitUpstreamUrl on setGitUpstream', () => {
    // given
    const gitData = {
      remoteURI: 'https://test.git',
      gitUser: 'username',
      gitToken: 'token'
    };
    component.projectName = 'sockshop';

    // when
    const spy = spyOn(dataService, 'setGitUpstreamUrl').and.callThrough();
    component.updateGitData(gitData);
    component.setGitUpstream();

    // then
    expect(spy).toHaveBeenCalled();
  });

  it('should have create mode disabled when routed to route data contains {isCreateMode: false}', () => {
    expect(component.isCreateMode).toBeFalse();
  });

  it('should set project name to projectName retrieved by route', () => {
    expect(component.projectName).toEqual('sockshop');
  });

  it('should have an pattern validation error when project name does not match: first letter lowercase, only lowercase, numbers and hyphens allowed', () => {
    // given
    component.isCreateMode = true;

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
    component.projectName = 'sockshop';

    // when
    const router = TestBed.inject(Router);
    const routeSpy = spyOn(router, 'navigate');
    component.deleteProject('sockshop');

    // then
    expect(routeSpy).toHaveBeenCalled();
  });

  it('should not show a notification when the component is initialized', () => {
    // given
    // Has to be retrieved by document, as it is not created at component level
    const notifications = document.getElementsByTagName('dt-confirmation-dialog-state');

    // then
    expect(component.unsavedDialogState).toBeNull();
    expect(notifications.length).toEqual(0);
  });

  it('should show a notification for unsaved changes when project name is changed', () => {
    // given
    component.isCreateMode = true;
    fixture.detectChanges();

    // when
    const inputEl = fixture.nativeElement.querySelector('#projectNameInput');
    inputEl.value = 'sockshop';
    inputEl.dispatchEvent(new InputEvent('input'));
    fixture.detectChanges();

    // then
    const notification = document.getElementsByTagName('dt-confirmation-dialog-state');
    expect(component.unsavedDialogState).toEqual(UNSAVED_DIALOG_STATE);
    expect(notification.length).toEqual(1);
  });

  it('should show a notification for unsaved changes when git data is changed in create mode', () => {
    // given
    component.isCreateMode = true;
    fixture.detectChanges();

    // when
    component.updateGitData({gitUser: 'someUser'});
    fixture.detectChanges();

    // then
    const notification = document.getElementsByTagName('dt-confirmation-dialog-state');
    expect(component.unsavedDialogState).toEqual(UNSAVED_DIALOG_STATE);
    expect(notification.length).toEqual(1);
  });

  it('should show a notification for unsaved changes when git data is changed in update mode', () => {
    // given
    component.isCreateMode = false;
    fixture.detectChanges();

    // when
    component.updateGitData({gitUser: 'someUser'});
    fixture.detectChanges();

    // then
    const notification = document.getElementsByTagName('dt-confirmation-dialog-state');
    expect(component.unsavedDialogState).toEqual(UNSAVED_DIALOG_STATE);
    expect(notification.length).toEqual(1);
  });

  it('should show a notification for unsaved changes when shipyard file is changed', () => {
    // given
    component.isCreateMode = true;
    fixture.detectChanges();

    // when
    component.updateShipyardFile(new File(['test'], 'test.yaml'));
    fixture.detectChanges();

    // then
    const notification = document.getElementsByTagName('dt-confirmation-dialog-state');
    expect(component.unsavedDialogState).toEqual(UNSAVED_DIALOG_STATE);
    expect(notification.length).toEqual(1);
  });

  it('should not show a notification when the notification was closed', () => {
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


