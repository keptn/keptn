import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSubscriptionItemComponent } from './ktb-subscription-item.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, convertToParamMap, Router } from '@angular/router';
import { of } from 'rxjs';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { DataService } from '../../_services/data.service';
import { KtbSubscriptionItemModule } from './ktb-subscription-item.module';
import { RouterTestingModule } from '@angular/router/testing';
import { AppRouting } from '../../app.routing';
import { APP_BASE_HREF } from '@angular/common';
import { environment } from '../../../environments/environment';

describe('KtbSubscriptionItemComponent', () => {
  let component: KtbSubscriptionItemComponent;
  let fixture: ComponentFixture<KtbSubscriptionItemComponent>;
  let subscription: UniformSubscription;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbSubscriptionItemModule, HttpClientTestingModule, RouterTestingModule, AppRouting],
      providers: [
        { provide: ApiService, useClass: ApiServiceMock },
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of(convertToParamMap({ projectName: 'sockshop' })),
          },
        },
        {
          provide: APP_BASE_HREF,
          useValue: environment.baseUrl,
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSubscriptionItemComponent);
    component = fixture.componentInstance;
    TestBed.inject(DataService).loadProjects().subscribe();
    fixture.detectChanges();

    subscription = new UniformSubscription('sockshop');
    subscription.id = 'mySubscriptionId';
  });

  it('should create and have given project set', () => {
    expect(component).toBeTruthy();
    expect(component.project?.projectName).toEqual('sockshop');
  });

  it('should navigate to subscription to edit', () => {
    // given
    const router = TestBed.inject(Router);
    const routerSpy = jest.spyOn(router, 'navigate');
    component.integrationId = 'myIntegrationId';

    // when
    component.editSubscription(subscription);

    // then
    expect(routerSpy).toHaveBeenCalled();
    expect(routerSpy).toHaveBeenCalledWith([
      '/',
      'project',
      'sockshop',
      'settings',
      'uniform',
      'integrations',
      'myIntegrationId',
      'subscriptions',
      'mySubscriptionId',
      'edit',
    ]);
  });

  it('should trigger a deletion dialog', () => {
    // given, when
    component.triggerDeleteSubscription(subscription);

    // then
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    expect(component.currentSubscription).toEqual(subscription);
    expect(component.deleteState).toEqual('confirm');
  });

  it('should delete a subscription', () => {
    // given
    component.integrationId = 'myIntegrationId';
    component.subscription = subscription;
    component.isWebhookService = false;
    const dataService = TestBed.inject(DataService);
    const spy = jest.spyOn(dataService, 'deleteSubscription');

    // when
    component.deleteSubscription();

    // then
    expect(spy).toHaveBeenCalledWith('myIntegrationId', 'mySubscriptionId', false);
  });
});
