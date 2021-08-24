import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbUniformSubscriptionsComponent } from './ktb-uniform-subscriptions.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';
import { UniformRegistrationsMock } from '../../_models/uniform-registrations.mock';

describe('KtbUniformSubscriptionsComponent', () => {
  let component: KtbUniformSubscriptionsComponent;
  let fixture: ComponentFixture<KtbUniformSubscriptionsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
      providers: [
        {provide: DataService, useClass: DataServiceMock},
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({}),
            params: of({}),
            paramMap: of(convertToParamMap({
              projectName: 'sockshop'
            })),
            queryParams: of({})
          }
        }
      ]
    }).compileComponents();
    fixture = TestBed.createComponent(KtbUniformSubscriptionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should disable "Add subscription"', () => {
    // given
    component.uniformRegistration = UniformRegistrationsMock[0];
    fixture.detectChanges();

    // then
    expect(fixture.nativeElement.querySelector('*[uitestid=addSubscriptionButton]').getAttribute('disabled')).not.toBeNull();
  });

  it('should enable "Add subscription"', () => {
    // given
    component.uniformRegistration = UniformRegistrationsMock[1];
    fixture.detectChanges();

    // then
    expect(fixture.nativeElement.querySelector('*[uitestid=addSubscriptionButton]').getAttribute('disabled')).toBeNull();
  });

  it('should delete subscription', () => {
    // given
    component.uniformRegistration = UniformRegistrationsMock[1];
    const subscription = component.uniformRegistration.subscriptions[0];

    // when
    component.deleteSubscription(subscription);

    // then
    expect(component.uniformRegistration.subscriptions.length).toEqual(0);
  });
});
