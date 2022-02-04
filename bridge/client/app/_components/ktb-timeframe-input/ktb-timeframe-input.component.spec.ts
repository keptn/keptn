import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbTimeframeInputComponent } from './ktb-timeframe-input.component';
import { AppModule } from '../../app.module';

describe('KtbTimeframeInputComponent', () => {
  let component: KtbTimeframeInputComponent;
  let fixture: ComponentFixture<KtbTimeframeInputComponent>;

  const formControlNames = ['hours', 'minutes', 'seconds', 'millis', 'micros'];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbTimeframeInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should validate input for formControls with a min value and set appropriate value to formControl', () => {
    for (const control of formControlNames) {
      // given
      component.timeframeForm.controls[control].setValue(-1);

      // when
      component.validateInput(control, 0, 24);

      // then
      expect(component.timeframeForm.controls[control].value).toEqual(0);
    }
  });

  it('should validate input for formControls with a max value and set appropriate value to formControl', () => {
    for (const control of formControlNames) {
      // given
      component.timeframeForm.controls[control].setValue(25);

      // when
      component.validateInput(control, 0, 24);

      // then
      expect(component.timeframeForm.controls[control].value).toEqual(24);
    }
  });

  it('should validate input for formControls, round input and set appropriate value to formControl', () => {
    for (const control of formControlNames) {
      // given
      component.timeframeForm.controls[control].setValue(1.25);

      // when
      component.validateInput(control, 0, 24);

      // then
      expect(component.timeframeForm.controls[control].value).toEqual(1);
    }
  });

  it('should emit given values', () => {
    const spy = jest.spyOn(component.timeframe, 'emit');
    for (const control of formControlNames) {
      // given
      component.timeframeForm.controls[control].setValue(1);

      // when
      component.validateInput(control, 0, 24);
    }

    // then
    expect(spy).toHaveBeenCalledWith({ hours: 1, minutes: 1, seconds: 1, millis: 1, micros: 1 });
  });
});
