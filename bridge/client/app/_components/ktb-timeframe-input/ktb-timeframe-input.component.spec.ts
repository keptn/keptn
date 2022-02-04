import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbTimeframeInputComponent } from './ktb-timeframe-input.component';

describe('KtbTimeframeInputComponent', () => {
  let component: KtbTimeframeInputComponent;
  let fixture: ComponentFixture<KtbTimeframeInputComponent>;

  beforeEach(async () => {
    await TestBed.compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbTimeframeInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
