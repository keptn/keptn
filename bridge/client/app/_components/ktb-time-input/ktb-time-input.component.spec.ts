import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbTimeInputComponent } from './ktb-time-input.component';
import { AppModule } from '../../app.module';

describe('KtbTimeInputComponent', () => {
  let component: KtbTimeInputComponent;
  let fixture: ComponentFixture<KtbTimeInputComponent>;

  const formControlNames = ['hours', 'minutes', 'seconds', 'millis', 'micros'];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbTimeInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should validate input for formControls with a min value and set appropriate value to formControl', () => {
    for (const control of formControlNames) {
      // given
      component.timeControls[control].setValue(-1);

      // when
      component.validateInput(control, undefined, 24);

      // then
      expect(component.timeControls[control].value).toEqual(0);
    }
  });

  it('should validate input for formControls with a max value and set appropriate value to formControl', () => {
    for (const control of formControlNames) {
      // given
      component.timeControls[control].setValue(25);

      // when
      component.validateInput(control, undefined, 24);

      // then
      expect(component.timeControls[control].value).toEqual(24);
    }
  });

  it('should validate input for formControls, round input and set appropriate value to formControl', () => {
    for (const control of formControlNames) {
      // given
      component.timeControls[control].setValue(1.25);

      // when
      component.validateInput(control, undefined, 24);

      // then
      expect(component.timeControls[control].value).toEqual(1);
    }
  });

  it('should validate input for formControls with min set to undefined, should be 0 for min', () => {
    for (const control of formControlNames) {
      // given
      component.timeControls[control].setValue(-1);

      // when
      component.validateInput(control, undefined, undefined);

      // then
      expect(component.timeControls[control].value).toEqual(0);
    }
  });

  it('should validate input for formControls with max set to undefined - round but values is used as given', () => {
    for (const control of formControlNames) {
      // given
      component.timeControls[control].setValue(10000.25);

      // when
      component.validateInput(control, undefined, undefined);

      // then
      expect(component.timeControls[control].value).toEqual(10000);
    }
  });

  it('should validate input for formControls with min set to different value than 0', () => {
    for (const control of formControlNames) {
      // given
      component.timeControls[control].setValue(-1);

      // when
      component.validateInput(control, 5, undefined);

      // then
      expect(component.timeControls[control].value).toEqual(5);
    }
  });

  it('should emit given values', () => {
    const spy = jest.spyOn(component.timeChanged, 'emit');
    for (const control of formControlNames) {
      // given
      component.timeControls[control].setValue(1);

      // when
      component.validateInput(control, undefined, 24);
    }

    // then
    expect(spy).toHaveBeenCalledWith({ hours: 1, minutes: 1, seconds: 1, millis: 1, micros: 1 });
  });

  it('should emit 0 as values', () => {
    const spy = jest.spyOn(component.timeChanged, 'emit');
    for (const control of formControlNames) {
      // given
      component.timeControls[control].setValue(0);

      // when
      component.validateInput(control, undefined, 24);
    }

    // then
    expect(spy).toHaveBeenCalledWith({ hours: 0, minutes: 0, seconds: 0, millis: 0, micros: 0 });
  });

  it('should emit undefined for not given values', () => {
    const spy = jest.spyOn(component.timeChanged, 'emit');
    for (const control of formControlNames) {
      // given
      component.timeControls[control].setValue(null);

      // when
      component.validateInput(control, undefined, 24);
    }

    // then
    expect(spy).toHaveBeenCalledWith({
      hours: undefined,
      minutes: undefined,
      seconds: undefined,
      millis: undefined,
      micros: undefined,
    });
  });

  it('should set all components to enabled by default', () => {
    for (const control of Object.values(component.timeControls)) {
      expect(control.enabled).toBe(true);
    }
  });

  it('should set all components to disabled and then to enabled', () => {
    // when
    component.disabled = true;

    // then
    for (const control of Object.values(component.timeControls)) {
      expect(control.enabled).toBe(false);
    }

    // when
    component.disabled = false;

    // then
    for (const control of Object.values(component.timeControls)) {
      expect(control.enabled).toBe(true);
    }
  });
});
