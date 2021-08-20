import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSubscriptionItemComponent } from './ktb-subscription-item.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../../app.module';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';

describe('KtbSubscriptionItemComponent', () => {
  let component: KtbSubscriptionItemComponent;
  let fixture: ComponentFixture<KtbSubscriptionItemComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
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
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
