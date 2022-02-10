import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbTriggerSequenceComponent } from './ktb-trigger-sequence.component';

describe('KtbTriggerSequenceComponent', () => {
  let component: KtbTriggerSequenceComponent;
  let fixture: ComponentFixture<KtbTriggerSequenceComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [KtbTriggerSequenceComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbTriggerSequenceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
