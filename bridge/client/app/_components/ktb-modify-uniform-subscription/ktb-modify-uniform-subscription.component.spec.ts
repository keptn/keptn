import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbModifyUniformSubscriptionComponent } from './ktb-modify-uniform-subscription.component';
import { AppModule } from '../../app.module';
import { ActivatedRoute, convertToParamMap, ParamMap, Router } from '@angular/router';
import { UniformRegistrationsMock } from '../../_models/uniform-registrations.mock';
import { BehaviorSubject, of } from 'rxjs';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { UniformRegistrationLocations } from '../../../../shared/interfaces/uniform-registration-locations';

describe('KtbModifyUniformSubscriptionComponent', () => {
  let component: KtbModifyUniformSubscriptionComponent;
  let fixture: ComponentFixture<KtbModifyUniformSubscriptionComponent>;
  let paramMap: BehaviorSubject<ParamMap>;

  beforeEach(async () => {
    paramMap = new BehaviorSubject<ParamMap>(convertToParamMap({
      projectName: 'sockshop',
      integrationId: UniformRegistrationsMock[0].id,
    }));
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
      providers: [
        {provide: DataService, useClass: DataServiceMock},
        {
          provide: ActivatedRoute, useValue: {
            paramMap: paramMap.asObservable(),
          },
        },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbModifyUniformSubscriptionComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should have disabled button', () => {
    // given
    fixture.detectChanges();

    // when first and second invalid
    // then
    expect(fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').getAttribute('disabled')).not.toBeNull();

    // when first valid and second invalid
    // @ts-ignore
    component.taskControl.setValue('deployment');
    fixture.detectChanges();
    // then
    expect(fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').getAttribute('disabled')).not.toBeNull();

    // when first invalid and second valid
    // @ts-ignore
    component.taskControl.setValue('');
    // @ts-ignore
    component.taskSuffixControl.setValue('triggered');
    fixture.detectChanges();
    // then
    expect(fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').getAttribute('disabled')).not.toBeNull();
  });

  it('should have enabled button', () => {
    // given
    fixture.detectChanges();
    // when
    // @ts-ignore
    component.taskControl.setValue('deployment');
    // @ts-ignore
    component.taskSuffixControl.setValue('triggered');
    fixture.detectChanges();

    // then
    expect(fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').getAttribute('disabled')).toBeNull();
  });

  it('should fill data', () => {
    // given
    const subscription = setSubscription(1, 0);
    fixture.detectChanges();

    // then
    // @ts-ignore
    expect(component.isGlobalControl.value).toEqual(true);
    // @ts-ignore
    expect(component.taskControl.value).toEqual('deployment');
    // @ts-ignore
    expect(component.taskSuffixControl.value).toEqual('triggered');

    const filterPairs: HTMLElement[] = Array.from(fixture.nativeElement.querySelectorAll('.dt-filter-field-tag-container'));
    expect(subscription.filter.stages?.every(stage => filterPairs.some(pair => pair.textContent === `Stage${stage}`))).toEqual(true);
    expect(subscription.filter.services?.every(service => filterPairs.some(pair => pair.textContent === `Service${service}`))).toEqual(true);
    expect(filterPairs.length).toEqual((subscription.filter.stages?.length ?? 0) + (subscription.filter.services?.length ?? 0));
    expect(fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').getAttribute('disabled')).toBeNull();
  });

  it('should update subscription', () => {
    // given
    const subscription = setSubscription(2, 0);
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').textContent.trim()).toEqual('Update subscription');
    const dataService = TestBed.inject(DataService);
    const updateSpy = jest.spyOn(dataService, 'updateUniformSubscription');
    // when
    fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').click();
    fixture.detectChanges();

    // then
    expect(updateSpy).toHaveBeenCalledWith(UniformRegistrationsMock[2].id, subscription);
    expect(subscription.filter.projects?.includes('sockshop')).toEqual(true);
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
    // @ts-ignore
    component.taskControl.setValue('deployment');
    // @ts-ignore
    component.taskSuffixControl.setValue('triggered');
    // @ts-ignore
    component.isGlobalControl.setValue(true);
    fixture.detectChanges();
    fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').click();
    fixture.detectChanges();

    // then
    expect(updateSpy).toHaveBeenCalledWith(UniformRegistrationsMock[1].id, Object.assign(subscription,
      {
        event: 'sh.keptn.event.deployment.triggered',
        _filter: [],
        filter: {
          projects: [],
          services: [],
          stages: [],
        },
      }));
    expect(routerSpy).toHaveBeenCalledWith(['/', 'project', 'sockshop', 'uniform', 'services', UniformRegistrationsMock[1].id]);
    expect(fixture.nativeElement.querySelector('button[uitestid=updateSubscriptionButton]').textContent.trim()).toEqual('Create subscription');
  });

  it('should only have triggered suffix', () => {
    // given
    setSubscription(6, 0);
    fixture.detectChanges();

    expect(component.suffixes).toEqual([{displayValue: 'triggered', value: 'triggered'}]);
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
      }]);
  });

  function setSubscription(integrationIndex: number, subscriptionIndex?: number): UniformSubscription {
    const dataService = TestBed.inject(DataService);
    const uniformRegistration = UniformRegistrationsMock[integrationIndex];
    const subscription = subscriptionIndex !== undefined ? uniformRegistration.subscriptions[subscriptionIndex] : new UniformSubscription('sockshop');
    dataService.getUniformSubscription = jest.fn().mockReturnValue(of(subscription));
    dataService.getIsUniformRegistrationControlPlane = jest.fn().mockReturnValue(of(uniformRegistration.metadata.location === UniformRegistrationLocations.CONTROL_PLANE));
    paramMap.next(convertToParamMap({
      projectName: 'sockshop',
      integrationId: uniformRegistration.id,
      subscriptionId: subscription.id,
    }));
    // set it again because of paramMap change
    fixture = TestBed.createComponent(KtbModifyUniformSubscriptionComponent);
    component = fixture.componentInstance;
    return subscription;
  }
});
