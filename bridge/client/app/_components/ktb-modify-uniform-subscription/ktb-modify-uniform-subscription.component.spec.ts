import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbModifyUniformSubscriptionComponent } from './ktb-modify-uniform-subscription.component';
import { AppModule } from '../../app.module';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { UniformRegistrationsMock } from '../../_models/uniform-registrations.mock';
import { of } from 'rxjs';

describe('KtbModifyUniformSubscriptionComponent', () => {
  let component: KtbModifyUniformSubscriptionComponent;
  let fixture: ComponentFixture<KtbModifyUniformSubscriptionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule
      ],
      declarations: [KtbModifyUniformSubscriptionComponent],
      providers: [
        {
          provide: ActivatedRoute, useValue: {
            paramMap: of(convertToParamMap({
              projectName: 'sockshop',
              integrationId: UniformRegistrationsMock[0].id
            }))
          }
        }
      ]
    })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbModifyUniformSubscriptionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
