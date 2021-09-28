import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbCreateServiceComponent } from './ktb-create-service.component';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, convertToParamMap, ParamMap, Router } from '@angular/router';
import { BehaviorSubject, of, throwError } from 'rxjs';
import { AppModule } from '../../app.module';
import { NotificationsService } from '../../_services/notifications.service';
import { NotificationType } from '../../_models/notification';
import { HttpErrorResponse } from '@angular/common/http';
import { DataServiceMock } from '../../_services/data.service.mock';

describe('KtbCreateServiceComponent', () => {
  let component: KtbCreateServiceComponent;
  let fixture: ComponentFixture<KtbCreateServiceComponent>;
  const projectName = 'sockshop';
  let queryParams: BehaviorSubject<ParamMap>;

  beforeEach(async () => {
    queryParams = new BehaviorSubject<ParamMap>(convertToParamMap({}));
    await TestBed.configureTestingModule({
      imports: [AppModule],
      providers: [
        { provide: DataService, useClass: DataServiceMock },
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
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
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

  it('should show required error', () => {
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
    expect(notificationSpy).toHaveBeenCalledWith(NotificationType.SUCCESS, 'Service successfully created!', 5_000);
  });

  it('should not create service', () => {
    // given
    const notificationService = TestBed.inject(NotificationsService);
    const notificationSpy = jest.spyOn(notificationService, 'addNotification');
    const dataService = TestBed.inject(DataService);
    dataService.createService = jest
      .fn()
      .mockReturnValue(throwError(new HttpErrorResponse({ error: 'service already exists' })));

    // when
    component.createService(projectName);

    // then
    expect(notificationSpy).toHaveBeenCalledWith(NotificationType.ERROR, 'service already exists', 5_000);
  });

  it('should go back', () => {
    // given
    const redirectTo = '%2Fproject%2Fsockshop%2Fservice';
    const cancelButton = fixture.nativeElement.querySelector('button[type=reset]');
    const router = TestBed.inject(Router);
    const routerNavigateSpy = jest.spyOn(router, 'navigateByUrl');
    updateQueryParams(redirectTo);

    // when
    cancelButton.click();
    fixture.detectChanges();

    // then
    expect(routerNavigateSpy).toHaveBeenCalledWith(redirectTo);
  });

  it('should go back to service overview', () => {
    // given
    const router = TestBed.inject(Router);
    const routerNavigateSpy = jest.spyOn(router, 'navigate');
    const route = TestBed.inject(ActivatedRoute);
    const cancelButton = fixture.nativeElement.querySelector('button[type=reset]');

    // when
    cancelButton.click();
    fixture.detectChanges();

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

  function updateQueryParams(redirectTo: string): void {
    queryParams.next(convertToParamMap({ redirectTo }));
    fixture = TestBed.createComponent(KtbCreateServiceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }
});
