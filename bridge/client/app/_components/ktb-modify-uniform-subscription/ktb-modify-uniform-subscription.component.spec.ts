import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbModifyUniformSubscriptionComponent } from './ktb-modify-uniform-subscription.component';
import { AppModule } from '../../app.module';
import { ActivatedRoute, convertToParamMap, ParamMap, Router } from '@angular/router';
import { UniformRegistrationsMock } from '../../_services/_mockData/uniform-registrations.mock';
import { BehaviorSubject, of, throwError } from 'rxjs';
import { DataService } from '../../_services/data.service';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { UniformRegistrationLocations } from '../../../../shared/interfaces/uniform-registration-locations';
import { UniformRegistrationInfo } from '../../../../shared/interfaces/uniform-registration-info';
import { WebhookConfig } from '../../../../shared/models/webhook-config';
import { HttpErrorResponse } from '@angular/common/http';
import { AbstractControl } from '@angular/forms';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';

describe('KtbModifyUniformSubscriptionComponent', () => {
  let component: KtbModifyUniformSubscriptionComponent;
  let fixture: ComponentFixture<KtbModifyUniformSubscriptionComponent>;
  let paramMap: BehaviorSubject<ParamMap>;

  beforeEach(async () => {
    paramMap = new BehaviorSubject<ParamMap>(
      convertToParamMap({
        projectName: 'sockshop',
        integrationId: UniformRegistrationsMock[0].id,
      })
    );
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
      providers: [
        { provide: ApiService, useClass: ApiServiceMock },
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: paramMap.asObservable(),
          },
        },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbModifyUniformSubscriptionComponent);
    component = fixture.componentInstance;
    TestBed.inject(DataService).loadProjects();
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should have disabled button if first and second control is invalid', async () => {
    // given
    fixture.detectChanges();

    // then
    assertIsUpdateButtonEnabled(false);
  });

  it('should have disabled button if first control is valid and second control is invalid', async () => {
    // given
    fixture.detectChanges();
    const taskControl = getTaskPrefix();
    taskControl.setValue('deployment');
    fixture.detectChanges();
    // then
    assertIsUpdateButtonEnabled(false);
  });

  it('should have disabled button if first control is invalid and second control is valid', async () => {
    // given
    fixture.detectChanges();
    const taskControl = getTaskPrefix();
    taskControl.setValue('');
    const taskSuffixControl = getTaskSuffix();
    taskSuffixControl.setValue('triggered');
    fixture.detectChanges();
    // then
    assertIsUpdateButtonEnabled(false);
  });

  it('should have disabled button if filter contains service but not a stage', async () => {
    // given
    const data = await component.data$?.toPromise();
    fixture.detectChanges();

    // when
    const taskControl = getTaskPrefix();
    taskControl.setValue('deployment');
    if (data) {
      data.subscription.filter = {
        projects: ['sockshop'],
        stages: [],
        services: ['carts'],
      };
    }
    fixture.detectChanges();
    // then
    expect(data ? component.isFormValid(data.subscription) : undefined).toBe(false);
  });

  it('should have disabled button if loading', () => {
    // given
    fixture.detectChanges();
    // when
    const taskControl = getTaskPrefix();
    taskControl.setValue('deployment');
    const taskSuffixControl = getTaskSuffix();
    taskSuffixControl.setValue('triggered');
    component.updating = true;
    fixture.detectChanges();

    // then
    assertIsUpdateButtonEnabled(false);
  });

  it('should have enabled button if task is valid', () => {
    // given
    fixture.detectChanges();
    const taskControl = getTaskPrefix();
    taskControl.setValue('deployment');
    const taskSuffixControl = getTaskSuffix();
    taskSuffixControl.setValue('triggered');
    fixture.detectChanges();

    // then
    assertIsUpdateButtonEnabled(true);
  });

  it('should have enabled button if filter contains a stage and a service', async () => {
    // given
    const data = await component.data$?.toPromise();
    fixture.detectChanges();
    // when
    const taskControl = getTaskPrefix();
    taskControl.setValue('deployment');
    const taskSuffixControl = getTaskSuffix();
    taskSuffixControl.setValue('triggered');

    if (data) {
      data.subscription.filter = {
        projects: ['sockshop'],
        stages: ['staging'],
        services: ['carts'],
      };
    }
    fixture.detectChanges();

    // then
    expect(data ? component.isFormValid(data.subscription) : undefined).toBe(true);
  });

  it('should have enabled button if filter contains just a stages', async () => {
    // given
    const data = await component.data$?.toPromise();
    fixture.detectChanges();
    // when
    const taskControl = getTaskPrefix();
    taskControl.setValue('deployment');
    const taskSuffixControl = getTaskSuffix();
    taskSuffixControl.setValue('triggered');

    // when
    if (data) {
      data.subscription.filter = {
        projects: ['sockshop'],
        stages: ['staging'],
        services: [],
      };
    }
    fixture.detectChanges();

    // then
    expect(data ? component.isFormValid(data.subscription) : undefined).toBe(true);
  });

  it('should have a disabled button if the webhook form is invalid', () => {
    // given
    setSubscription(10);
    fixture.detectChanges();
    component.isWebhookFormValid = false;
    fixture.detectChanges();

    // then
    assertIsUpdateButtonEnabled(false);
  });

  it('should have a enabled button if the webhook form is valid', () => {
    // given
    setSubscription(10, 0);
    fixture.detectChanges();
    component.isWebhookFormValid = true;
    fixture.detectChanges();

    // then
    assertIsUpdateButtonEnabled(true);
  });

  it('should set the right properties and enable the button when a global subscription is set', () => {
    // given
    setSubscription(1, 0);
    fixture.detectChanges();

    // then
    const isGlobalControl = getIsGlobalControl();
    expect(isGlobalControl.value).toEqual(true);
    const taskControl = getTaskPrefix();
    expect(taskControl.value).toEqual('deployment');
    const taskSuffixControl = getTaskSuffix();
    expect(taskSuffixControl.value).toEqual('triggered');

    const filterPairs: HTMLElement[] = Array.from(
      fixture.nativeElement.querySelectorAll('.dt-filter-field-tag-container')
    );
    expect(filterPairs.length).toEqual(0);
    assertIsUpdateButtonEnabled(true);
  });

  it('should set the right properties and enable the button when a subscription is set', () => {
    // given
    const subscription = setSubscription(2, 0);
    fixture.detectChanges();

    // then
    const isGlobalControl = getIsGlobalControl();
    expect(isGlobalControl.value).toEqual(false);
    const taskControl = getTaskPrefix();
    expect(taskControl.value).toEqual('test');
    const taskSuffixControl = getTaskSuffix();
    expect(taskSuffixControl.value).toEqual('triggered');

    const filterPairs: HTMLElement[] = Array.from(
      fixture.nativeElement.querySelectorAll('.dt-filter-field-tag-container')
    );
    expect(
      subscription.filter.stages?.every((stage) => filterPairs.some((pair) => pair.textContent === `Stage${stage}`))
    ).toEqual(true);
    expect(
      subscription.filter.services?.every((service) =>
        filterPairs.some((pair) => pair.textContent === `Service${service}`)
      )
    ).toEqual(true);
    expect(filterPairs.length).toEqual(
      (subscription.filter.stages?.length ?? 0) + (subscription.filter.services?.length ?? 0)
    );
    assertIsUpdateButtonEnabled(true);
  });

  it('should update subscription', () => {
    // given
    const subscription = setSubscription(2, 0);
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').textContent.trim()).toEqual(
      'Update subscription'
    );
    const dataService = TestBed.inject(DataService);
    const updateSpy = jest.spyOn(dataService, 'updateUniformSubscription');
    // when
    fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').click();
    fixture.detectChanges();

    // then
    expect(updateSpy).toHaveBeenCalledWith(UniformRegistrationsMock[2].id, subscription, undefined);
    expect(subscription.filter.projects?.includes('sockshop')).toEqual(true);
  });

  it('should update subscription for all keptn events with keptn.sh.>', () => {
    // given
    const subscription = setSubscription(2, 0);
    /* eslint-disable @typescript-eslint/ban-ts-comment */
    // @ts-ignore //Ignore private property
    component.taskControl.setValue('sh.keptn');
    component.taskSuffixControl.setValue('>');
    // @ts-ignore //Ignore private property
    component.isGlobalControl.setValue(true);
    component.editMode = true;
    /* eslint-enable */

    // when
    component.updateSubscription('sockshop', UniformRegistrationsMock[2].id, subscription, undefined);

    // then
    expect(subscription.event).toEqual('sh.keptn.>');
  });

  it('should update subscription for deplyoment keptn wildcard events with keptn.sh.event.approval.>', () => {
    // given
    const subscription = setSubscription(2, 0);
    /* eslint-disable @typescript-eslint/ban-ts-comment */
    // @ts-ignore //Ignore private property
    component.taskControl.setValue('deployment');
    component.taskSuffixControl.setValue('>');
    // @ts-ignore //Ignore private property
    component.isGlobalControl.setValue(true);
    component.editMode = true;
    /* eslint-enable */

    // when
    component.updateSubscription('sockshop', UniformRegistrationsMock[2].id, subscription, undefined);

    // then
    expect(subscription.event).toEqual('sh.keptn.event.deployment.>');
  });

  it('should update subscription and webhook', () => {
    // given
    const subscription = setSubscription(10, 0);
    const dataService = TestBed.inject(DataService);
    const updateSpy = jest.spyOn(dataService, 'updateUniformSubscription');
    const webhookConfig = new WebhookConfig();
    fixture.detectChanges();

    webhookConfig.method = 'POST';
    webhookConfig.url = 'https://keptn.sh';
    webhookConfig.payload = '{}';
    webhookConfig.header = [{ name: 'Content-Type', value: 'application/json' }];

    // when
    component.updateSubscription('sockshop', UniformRegistrationsMock[10].id, subscription, webhookConfig);
    fixture.detectChanges();

    webhookConfig.type = subscription.event;
    webhookConfig.filter = subscription.filter;
    webhookConfig.prevConfiguration = {
      filter: subscription.filter,
      type: subscription.event,
    };

    // then
    expect(updateSpy).toHaveBeenCalledWith(UniformRegistrationsMock[10].id, subscription, webhookConfig);
  });

  it('should revert loading if request fails', () => {
    // given
    const subscription = setSubscription(2, 0);
    const dataService = TestBed.inject(DataService);
    dataService.updateUniformSubscription = jest.fn().mockReturnValue(throwError(new HttpErrorResponse({ error: '' })));
    fixture.detectChanges();

    // when
    component.updateSubscription('sockshop', UniformRegistrationsMock[2].id, subscription);

    // then
    expect(component.updating).toEqual(false);
  });

  it('should create global subscription', () => {
    // given
    const subscription = setSubscription(1);
    fixture.detectChanges();
    const dataService = TestBed.inject(DataService);
    const route = TestBed.inject(Router);
    const updateSpy = jest.spyOn(dataService, 'createUniformSubscription');
    const routerSpy = jest.spyOn(route, 'navigate');
    // when
    const taskControl = getTaskPrefix();
    taskControl.setValue('deployment');
    const taskSuffixControl = getTaskSuffix();
    taskSuffixControl.setValue('triggered');
    const isGlobalControl = getIsGlobalControl();
    isGlobalControl.setValue(true);
    fixture.detectChanges();
    fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').click();
    fixture.detectChanges();

    // then
    expect(updateSpy).toHaveBeenCalledWith(
      UniformRegistrationsMock[1].id,
      Object.assign(subscription, {
        event: 'sh.keptn.event.deployment.triggered',
        _filter: [],
        filter: {
          projects: [],
          services: [],
          stages: [],
        },
      }),
      undefined
    );
    expect(routerSpy).toHaveBeenCalledWith([
      '/',
      'project',
      'sockshop',
      'uniform',
      'services',
      UniformRegistrationsMock[1].id,
    ]);
    expect(fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').textContent.trim()).toEqual(
      'Create subscription'
    );
  });

  it('should only have triggered suffix', () => {
    // given
    setSubscription(6, 0);
    fixture.detectChanges();

    expect(component.suffixes).toEqual([{ displayValue: 'triggered', value: 'triggered' }]);
  });

  it('should show all suffixes', () => {
    // given
    setSubscription(3);
    fixture.detectChanges();

    expect(component.suffixes).toEqual([
      {
        value: '>',
        displayValue: '*',
      },
      {
        value: 'triggered',
        displayValue: 'triggered',
      },
      {
        value: 'started',
        displayValue: 'started',
      },
      {
        value: 'finished',
        displayValue: 'finished',
      },
    ]);
  });

  it('should show webhook form', () => {
    // given
    setSubscription(10, 0);
    fixture.detectChanges();

    // then
    const webhookForm = fixture.nativeElement.querySelector('ktb-webhook-settings');
    expect(webhookForm).toBeTruthy();
  });

  it('should not show webhook form', () => {
    // given
    setSubscription(1, 0);
    fixture.detectChanges();

    // then
    const webhookForm = fixture.nativeElement.querySelector('ktb-webhook-settings');
    expect(webhookForm).toBeFalsy();
  });

  it('should show project checkbox', () => {
    // given
    setSubscription(1, 0);
    fixture.detectChanges();
    const checkbox = fixture.nativeElement.querySelector('[uitestid=ktb-modify-subscription-project-checkbox]');

    // then
    expect(checkbox).toBeTruthy();
  });

  it('should not show project checkbox', () => {
    // given
    setSubscription(10, 0);
    fixture.detectChanges();
    const checkbox = fixture.nativeElement.querySelector('[uitestid=ktb-modify-subscription-project-checkbox]');

    // then
    expect(checkbox).toBeFalsy();
  });

  it('it should enable "use for all projects" checkbox if filter is cleared', () => {
    // given
    setSubscription(2, 0);
    fixture.detectChanges();

    // when
    fixture.nativeElement.querySelector('.dt-filter-field-clear-all-button').click();
    fixture.detectChanges();
    // then
    const isGlobalControl = getIsGlobalControl();
    expect(isGlobalControl.enabled).toEqual(true);
  });

  it('it should disable "use for all projects" checkbox and set to false if filter is set', () => {
    // given
    setSubscription(3, 0);
    fixture.detectChanges();

    // then
    const isGlobalControl = getIsGlobalControl();
    expect(isGlobalControl.disabled).toEqual(true);
    expect(isGlobalControl.value).toEqual(false);
  });

  it('should initially load intersected events', () => {
    const eventPayload = { data: {} };
    const dataService = TestBed.inject(DataService);
    dataService.getIntersectedEvent = jest.fn().mockReturnValue(of(eventPayload));
    setSubscription(10, 0);
    fixture.detectChanges();
    expect(component.eventPayload).toEqual(eventPayload);
  });

  function assertIsUpdateButtonEnabled(isEnabled: boolean): void {
    const element = expect(
      fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').getAttribute('disabled')
    );
    (isEnabled ? element : element.not).toBeNull();
  }

  function setSubscription(integrationIndex: number, subscriptionIndex?: number): UniformSubscription {
    const dataService = TestBed.inject(DataService);
    const uniformRegistration = UniformRegistrationsMock[integrationIndex];
    const subscription =
      subscriptionIndex !== undefined
        ? uniformRegistration.subscriptions[subscriptionIndex]
        : new UniformSubscription('sockshop');
    dataService.getUniformSubscription = jest.fn().mockReturnValue(of(subscription));
    dataService.getUniformRegistrationInfo = jest.fn().mockReturnValue(
      of({
        isControlPlane: uniformRegistration.metadata.location === UniformRegistrationLocations.CONTROL_PLANE,
        isWebhookService: uniformRegistration.isWebhookService,
      } as UniformRegistrationInfo)
    );
    paramMap.next(
      convertToParamMap({
        projectName: 'sockshop',
        integrationId: uniformRegistration.id,
        subscriptionId: subscription.id,
      })
    );
    // set it again because of paramMap change
    fixture = TestBed.createComponent(KtbModifyUniformSubscriptionComponent);
    component = fixture.componentInstance;
    return subscription;
  }

  function getTaskPrefix(): AbstractControl {
    return component.subscriptionForm.get('taskPrefix') as AbstractControl;
  }

  function getTaskSuffix(): AbstractControl {
    return component.subscriptionForm.get('taskSuffix') as AbstractControl;
  }

  function getIsGlobalControl(): AbstractControl {
    return component.subscriptionForm.get('isGlobal') as AbstractControl;
  }
});
