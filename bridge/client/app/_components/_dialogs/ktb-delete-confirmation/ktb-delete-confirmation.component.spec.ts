import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbDeleteConfirmationComponent } from './ktb-delete-confirmation.component';

describe('KtbDeleteConfirmationComponent', () => {
  let component: KtbDeleteConfirmationComponent;
  let fixture: ComponentFixture<KtbDeleteConfirmationComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbDeleteConfirmationComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbDeleteConfirmationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
