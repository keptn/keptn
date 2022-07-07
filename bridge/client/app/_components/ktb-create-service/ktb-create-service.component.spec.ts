import { ComponentFixture, fakeAsync, TestBed, tick } from '@angular/core/testing';
import { KtbCreateServiceComponent } from './ktb-create-service.component';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, convertToParamMap, ParamMap, Router } from '@angular/router';
import { BehaviorSubject, firstValueFrom, of, throwError } from 'rxjs';
import { NotificationsService } from '../../_services/notifications.service';
import { NotificationType } from '../../_models/notification';
import { HttpErrorResponse } from '@angular/common/http';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { KtbCreateServiceModule } from './ktb-create-service.module';
import { RouterTestingModule } from '@angular/router/testing';
import { take } from 'rxjs/operators';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('KtbCreateServiceComponent', () => {
  let component: KtbCreateServiceComponent;
  let fixture: ComponentFixture<KtbCreateServiceComponent>;
  const projectName = 'sockshop';
  let queryParams: BehaviorSubject<ParamMap>;

  beforeEach(async () => {
    queryParams = new BehaviorSubject<ParamMap>(convertToParamMap({}));
    await TestBed.configureTestingModule({
      imports: [KtbCreateServiceModule, RouterTestingModule, BrowserAnimationsModule],
      providers: [
        { provide: ApiService, useClass: ApiServiceMock },
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of(
              convertToParamMap({
                projectName,
              })
            ),
            queryParamMap: queryParams.asObservable(),
            snapshot: {},
          },
        },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbCreateServiceComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set the projectName', (done) => {
    component.projectName$.pipe(take(1)).subscribe((actualProjectName) => {
      expect(actualProjectName).toBe(projectName);
      done();
    });
  });

  it('should show duplicate error', () => {
    const serviceNames = ['carts', 'carts-db'];
    for (const serviceName of serviceNames) {
      component.serviceNameControl.setValue(serviceName);
      component.formGroup.updateValueAndValidity();
      expect(component.formGroup.hasError('duplicate'));
    }
  });

  it('should show pattern error', () => {
    const serviceNames = ['Service', '1service', '-service', '$service', 'serVice', 'ser_ice'];
    for (const serviceName of serviceNames) {
      component.serviceNameControl.setValue(serviceName);
      component.formGroup.updateValueAndValidity();
      expect(component.formGroup.hasError('pattern'));
    }
  });

  it('should show required error', async () => {
    await firstValueFrom(component.projectName$);
    await firstValueFrom(component.redirectTo$);

    const serviceNames = ['service', ''];
    for (const serviceName of serviceNames) {
      component.serviceNameControl.setValue(serviceName);
      component.formGroup.updateValueAndValidity();
    }
    fixture.detectChanges();
    expect(component.formGroup.hasError('required'));
    checkCreateButton(false);
  });

  it('should create service', () => {
    // given
    const serviceName = 'service-1';
    component.serviceNameControl.setValue(serviceName);
    component.formGroup.updateValueAndValidity();
    fixture.detectChanges();
    const notificationService = TestBed.inject(NotificationsService);
    const dataService = TestBed.inject(DataService);
    const notificationSpy = jest.spyOn(notificationService, 'addNotification');
    const loadProjectsSpy = jest.spyOn(dataService, 'loadProject');
    const createButton = getCreateButton();

    expect(component.formGroup.errors).toBeNull();
    checkCreateButton(true);

    // when
    createButton.click();
    fixture.detectChanges();

    // then
    expect(loadProjectsSpy).toHaveBeenCalledWith(projectName);
    expect(notificationSpy).toHaveBeenCalledWith(NotificationType.SUCCESS, 'Service successfully created!');
  });

  it('should not create service', fakeAsync(() => {
    // given
    fixture.detectChanges();
    const notificationService = TestBed.inject(NotificationsService);
    const notificationSpy = jest.spyOn(notificationService, 'addNotification');
    const inProgressSpy = jest.fn();
    Object.defineProperty(component, 'isCreating', {
      get: jest.fn(() => true),
      set: inProgressSpy,
    });
    const dataService = TestBed.inject(DataService);
    dataService.createService = jest
      .fn()
      .mockReturnValue(throwError(new HttpErrorResponse({ error: 'service already exists' })));

    // when
    component.createService(projectName, null);

    // then
    expect(inProgressSpy).toHaveBeenCalledWith(true);
    tick();
    expect(inProgressSpy).toHaveBeenCalledWith(false);
    expect(notificationSpy).toBeCalledTimes(0);
  }));

  it('should go back', async () => {
    // given
    const redirectTo = '%2Fproject%2Fsockshop%2Fservice';
    const router = TestBed.inject(Router);
    const routerNavigateSpy = jest.spyOn(router, 'navigateByUrl');

    // when
    try {
      await component.cancel(redirectTo);
    } catch (_err) {}

    // then
    expect(routerNavigateSpy).toHaveBeenCalledWith(redirectTo);
  });

  it('should go back to service overview', async () => {
    // given
    const router = TestBed.inject(Router);
    const routerNavigateSpy = jest.spyOn(router, 'navigate');
    const route = TestBed.inject(ActivatedRoute);

    /// when
    try {
      await component.cancel(null);
    } catch (_err) {}

    // then
    expect(routerNavigateSpy).toHaveBeenCalledWith(['../'], { relativeTo: route });
  });

  function checkCreateButton(isEnabled: boolean): void {
    const disabled: null | string = getCreateButton().getAttribute('disabled');
    if (isEnabled) {
      expect(disabled).toBeNull();
    } else {
      expect(disabled).not.toBeNull();
    }
  }

  function getCreateButton(): HTMLElement {
    return fixture.nativeElement.querySelector('button[uitestid=createServiceButton]');
  }
});
