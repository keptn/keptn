import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSubscriptionItemComponent } from './ktb-subscription-item.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../../app.module';
import { Subscription } from '../../_models/subscription';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';

describe('KtbSubscriptionItemComponent', () => {
  let component: KtbSubscriptionItemComponent;
  let fixture: ComponentFixture<KtbSubscriptionItemComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
      providers: [
        {
          provide: ActivatedRoute, useValue: {
            paramMap: of(convertToParamMap({projectName: 'sockshop'}))
          }
        }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSubscriptionItemComponent);
    component = fixture.componentInstance;
    component.subscription = new Subscription();
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
