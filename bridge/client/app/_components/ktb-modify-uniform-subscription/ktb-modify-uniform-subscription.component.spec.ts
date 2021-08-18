import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbModifyUniformSubscriptionComponent } from './ktb-modify-uniform-subscription.component';

describe('KtbModifyUniformSubscriptionComponent', () => {
  let component: KtbModifyUniformSubscriptionComponent;
  let fixture: ComponentFixture<KtbModifyUniformSubscriptionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbModifyUniformSubscriptionComponent ]
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
