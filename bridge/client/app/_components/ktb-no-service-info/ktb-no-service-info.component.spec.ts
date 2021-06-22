import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbNoServiceInfoComponent } from './ktb-no-service-info.component';

describe('KtbNoServiceInfoComponent', () => {
  let component: KtbNoServiceInfoComponent;
  let fixture: ComponentFixture<KtbNoServiceInfoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbNoServiceInfoComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbNoServiceInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
