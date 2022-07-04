import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbModifyUniformSubscriptionComponent } from './ktb-modify-uniform-subscription.component';
import { ActivatedRoute, convertToParamMap, ParamMap, Router } from '@angular/router';
import { UniformRegistrationsMock } from '../../_services/_mockData/uniform-registrations.mock';
import { BehaviorSubject, of, throwError } from 'rxjs';
import { DataService } from '../../_services/data.service';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { UniformRegistrationLocations } from '../../../../shared/interfaces/uniform-registration-locations';
import { UniformRegistrationInfo } from '../../../../shared/interfaces/uniform-registration-info';
import { HttpErrorResponse } from '@angular/common/http';
import { AbstractControl } from '@angular/forms';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { ProjectsMock } from '../../_services/_mockData/projects.mock';
import { KtbModifyUniformSubscriptionModule } from './ktb-modify-uniform-subscription.module';
import { RouterTestingModule } from '@angular/router/testing';
import { KtbKeptnServicesListComponent } from '../ktb-keptn-services-list/ktb-keptn-services-list.component';
import { IWebhookConfigClient } from '../../../../shared/interfaces/webhook-config';

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
      imports: [
        KtbModifyUniformSubscriptionModule,
        HttpClientTestingModule,
        RouterTestingModule.withRoutes([
          {
            path: 'project/:projectName/settings/uniform/integrations/:integrationId',
            component: KtbKeptnServicesListComponent,
          },
        ]),
      ],
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

  it('should update subscription', () => {
    // given
    const subscription = setSubscription(2, 0);
    fixture.detectChanges();
    const dataService = TestBed.inject(DataService);
    const updateSpy = jest.spyOn(dataService, 'updateUniformSubscription');
    component.updateSubscription('sockshop', UniformRegistrationsMock[2].id, subscription);

    // then
    expect(updateSpy).toHaveBeenCalledWith(UniformRegistrationsMock[2].id, subscription, undefined);
    expect(subscription.filter.projects?.includes('sockshop')).toEqual(true);
  });

  it('should update subscription for all keptn events with keptn.sh.>', () => {
    // given
    const subscription = setSubscription(2, 0);
    getTaskPrefix().setValue('sh.keptn');
    getTaskSuffix().setValue('>');
    getIsGlobalControl().setValue(true);
    component.editMode = true;

    // when
    component.updateSubscription('sockshop', UniformRegistrationsMock[2].id, subscription, undefined);

    // then
    expect(subscription.event).toEqual('sh.keptn.>');
  });

  it('should update subscription for deployment keptn wildcard events with keptn.sh.event.approval.>', () => {
    // given
    const subscription = setSubscription(2, 0);
    getTaskPrefix().setValue('deployment');
    getTaskSuffix().setValue('>');
    getIsGlobalControl().setValue(true);
    component.editMode = true;

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
    fixture.detectChanges();
    const webhookConfig: IWebhookConfigClient = {
      method: 'POST',
      url: 'https://keptn.sh',
      payload: '{}',
      header: [{ key: 'Content-Type', value: 'application/json' }],
      type: 'sh.keptn.event.evaluation.triggered',
      sendStarted: true,
      sendFinished: true,
    };

    // when
    component.updateSubscription('sockshop', UniformRegistrationsMock[10].id, subscription, webhookConfig);
    fixture.detectChanges();

    webhookConfig.type = subscription.event;
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
    component.updateSubscription('sockshop', UniformRegistrationsMock[1].id, subscription);

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
      'settings',
      'uniform',
      'integrations',
      UniformRegistrationsMock[1].id,
    ]);
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

  it('should initially load intersected events', () => {
    const eventPayload = { data: {} };
    const dataService = TestBed.inject(DataService);
    dataService.getIntersectedEvent = jest.fn().mockReturnValue(of(eventPayload));
    setSubscription(10, 0);
    fixture.detectChanges();
    expect(component.eventPayload).toEqual(eventPayload);
  });

  it('should remove deleted service from subscription', () => {
    // given
    const subscription = setSubscription(2, 0);
    fixture.detectChanges();
    const dataService = TestBed.inject(DataService);

    // when
    jest.spyOn(dataService, 'getProject').mockReturnValue(of(ProjectsMock[2]));
    component.data$.subscribe();

    // then
    expect(subscription.filter.services).toEqual([]);
    expect(component._dataSource.data).toEqual({
      autocomplete: [
        {
          name: 'Stage',
          autocomplete: [
            {
              name: 'dev',
            },
            {
              name: 'staging',
            },
            {
              name: 'production',
            },
          ],
        },
        {
          name: 'Service',
          autocomplete: [
            {
              name: 'carts-db',
            },
          ],
        },
      ],
    });
  });

  it('should add new service to datasource', () => {
    // given
    const dataService = TestBed.inject(DataService);
    jest.spyOn(dataService, 'getProject').mockReturnValue(of(ProjectsMock[2]));

    setSubscription(2, 0);

    // when
    jest.spyOn(dataService, 'getProject').mockReturnValue(of(ProjectsMock[0]));
    component.data$.subscribe();

    // then
    expect(component._dataSource.data).toEqual({
      autocomplete: [
        {
          name: 'Stage',
          autocomplete: [
            {
              name: 'dev',
            },
            {
              name: 'staging',
            },
            {
              name: 'production',
            },
          ],
        },
        {
          name: 'Service',
          autocomplete: [
            {
              name: 'carts-db',
            },
            {
              name: 'carts',
            },
          ],
        },
      ],
    });
  });

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
