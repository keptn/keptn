import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing';
import { KtbSubscriptionItemComponent } from './ktb-subscription-item.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../../app.module';
import { ActivatedRoute, convertToParamMap, Router } from '@angular/router';
import { of } from 'rxjs';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';
import { TestUtils } from '../../_utils/test.utils';

describe('KtbSubscriptionItemComponent', () => {
  let component: KtbSubscriptionItemComponent;
  let fixture: ComponentFixture<KtbSubscriptionItemComponent>;
  let subscription: UniformSubscription;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
      providers: [
        {provide: DataService, useClass: DataServiceMock},
        {
          provide: ActivatedRoute, useValue: {
            paramMap: of(convertToParamMap({projectName: 'sockshop'}))
          }
        }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSubscriptionItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    subscription = new UniformSubscription('sockshop');
    subscription.id = 'mySubscriptionId';
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should have enabled buttons', () => {
    // given
    component.subscription = subscription;
    component.integrationId = 'myIntegrationId';
    fixture.detectChanges();

    // then
    expect(fixture.nativeElement.querySelector('button.disabled')).toBeNull();
  });

  it('should have disabled buttons and functionality', () => {
    // given
    subscription.id = undefined;
    component.subscription = subscription;
    component.integrationId = 'myIntegrationId';
    fixture.detectChanges();
    const router = TestBed.inject(Router);
    const routeChange = jest.spyOn(router, 'navigate');
    const subscriptionDeleted = jest.spyOn(component.subscriptionDeleted, 'emit');

    // when
    fixture.nativeElement.querySelector('button[uitestid=subscriptionDeleteButton]').click();
    fixture.nativeElement.querySelector('button[uitestid=subscriptionEditButton]').click();

    // then
    expect(fixture.nativeElement.querySelectorAll('button.disabled').length).toEqual(2);
    expect(routeChange).not.toHaveBeenCalled();
    expect(subscriptionDeleted).not.toHaveBeenCalled();
  });

  it('should delete subscription', fakeAsync(() => {
    // given
    component.subscription = subscription;
    component.integrationId = 'myIntegrationId';
    fixture.detectChanges();
    const subscriptionDeleted = jest.spyOn(component.subscriptionDeleted, 'emit');

    // when
    fixture.nativeElement.querySelector('button[title=Delete]').click();
    TestUtils.updateDialog(fixture);

    expect(document.querySelector('dt-confirmation-dialog-state[name=confirm]')).toBeTruthy();
    expect(document.querySelector('dt-confirmation-dialog-state[name=confirm] *[uitestid=dialogWarningMessage]')).toBeFalsy();

    component.deleteSubscription();
    fixture.detectChanges();

    // then
    expect(subscriptionDeleted).toHaveBeenCalledWith(subscription);
  }));

  it('should show warning on delete subscription', fakeAsync(() => {
    // given
    subscription.filter.projects = [];
    component.subscription = subscription;
    component.integrationId = 'myIntegrationId';
    fixture.detectChanges();
    const subscriptionDeleted = jest.spyOn(component.subscriptionDeleted, 'emit');

    // when
    fixture.nativeElement.querySelector('button[title=Delete]').click();
    TestUtils.updateDialog(fixture);

    expect(document.querySelector('dt-confirmation-dialog-state[name=confirm]')).toBeTruthy();
    expect(document.querySelector('dt-confirmation-dialog-state[name=confirm] *[uitestid=dialogWarningMessage]')).toBeTruthy();

    component.deleteSubscription();
    fixture.detectChanges();

    // then
    expect(subscriptionDeleted).toHaveBeenCalledWith(subscription);
  }));

  it('should edit subscription', () => {
    // given
    component.subscription = subscription;
    component.integrationId = 'myIntegrationId';
    fixture.detectChanges();
    const router = TestBed.inject(Router);
    const routeChange = jest.spyOn(router, 'navigate');

    // when
    fixture.nativeElement.querySelector('button[title=Edit]').click();

    component.deleteSubscription();
    fixture.detectChanges();

    expect(routeChange).toHaveBeenCalledWith(['/', 'project', 'sockshop', 'uniform', 'services', 'myIntegrationId', 'subscriptions', 'mySubscriptionId', 'edit']);
  });
});
