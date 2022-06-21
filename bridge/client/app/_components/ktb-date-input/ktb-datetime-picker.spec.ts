import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbDatetimePickerComponent, KtbDatetimePickerDirective } from './ktb-datetime-picker.component';
import { ElementRef } from '@angular/core';
import { BrowserDynamicTestingModule } from '@angular/platform-browser-dynamic/testing';
import moment, { Moment } from 'moment';
import { Timeframe } from '../../_models/timeframe';
import { OverlayService } from '../../_directives/overlay-service/overlay.service';
import { KtbDateInputModule } from './ktb-date-input.module';
import { RouterTestingModule } from '@angular/router/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';

export class MockElementRef extends ElementRef {
  nativeElement = {};

  constructor() {
    super(null);
  }
}

describe('KtbDatetimePickerComponent', () => {
  let directive: KtbDatetimePickerDirective;
  let component: KtbDatetimePickerComponent;
  let fixture: ComponentFixture<KtbDatetimePickerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbDateInputModule, RouterTestingModule, HttpClientTestingModule],
      providers: [{ provide: ElementRef, useClass: MockElementRef }],
    })
      .overrideModule(BrowserDynamicTestingModule, { set: { entryComponents: [KtbDatetimePickerComponent] } })
      .compileComponents();

    fixture = TestBed.createComponent(KtbDatetimePickerComponent);
    component = fixture.componentInstance;
    directive = new KtbDatetimePickerDirective(TestBed.inject(ElementRef), TestBed.inject(OverlayService));
    directive.ngOnInit();
    fixture.detectChanges();
  });

  it('should create the component', () => {
    expect(component).toBeTruthy();
  });

  it('should set selectedDate to the given date as moment date', () => {
    // given
    const spy = jest.spyOn(component.selectedDateTime, 'emit');
    const date = moment();

    // when
    component.changeDate(date.toDate());
    component.setDateTime();

    // then
    expect(spy).toHaveBeenCalledWith(date.hours(0).minutes(0).seconds(0).milliseconds(0).toISOString());
  });

  it('should set selectedTime to given value', () => {
    // given
    const spy = jest.spyOn(component.selectedDateTime, 'emit');
    const date = moment('2022-01-01T00:00:00.000Z');
    const testDate = moment('2022-01-01T00:00:00.000Z').hours(2).minutes(15).seconds(10);
    component.secondsEnabled = true;
    const timeframe: Timeframe = {
      hours: 2,
      minutes: 15,
      seconds: 10,
      millis: undefined,
      micros: undefined,
    };

    // when
    component.changeDate(date.toDate());
    component.changeTime(timeframe);
    component.setDateTime();

    // then
    expect(spy).toHaveBeenCalledWith(testDate.toISOString());
  });

  it('should be disabled = false if hours and minutes are set and seconds are disabled ', () => {
    // given
    component.secondsEnabled = false;
    const timeframe: Timeframe = {
      hours: 1,
      minutes: 15,
      seconds: undefined,
      millis: undefined,
      micros: undefined,
    };

    // when
    component.changeTime(timeframe);

    // then
    expect(component.disabled).toEqual(false);
  });

  it('should be disabled = true if hours are not set and seconds are disabled', () => {
    // given
    component.secondsEnabled = false;
    const timeframe: Timeframe = {
      hours: undefined,
      minutes: 15,
      seconds: undefined,
      millis: undefined,
      micros: undefined,
    };

    // when
    component.changeTime(timeframe);

    // then
    expect(component.disabled).toEqual(true);
  });

  it('should be disabled = true if minutes are not set and seconds are disabled', () => {
    // given
    component.secondsEnabled = false;
    const timeframe: Timeframe = {
      hours: 1,
      minutes: undefined,
      seconds: undefined,
      millis: undefined,
      micros: undefined,
    };

    // when
    component.changeTime(timeframe);

    // then
    expect(component.disabled).toEqual(true);
  });

  it('should be disabled = false if hours and minutes are not set and seconds are disabled', () => {
    // given
    component.secondsEnabled = false;
    const timeframe: Timeframe = {
      hours: undefined,
      minutes: undefined,
      seconds: undefined,
      millis: undefined,
      micros: undefined,
    };

    // when
    component.changeTime(timeframe);

    // then
    expect(component.disabled).toEqual(false);
  });

  it('should be disabled = false if hours, minutes and seconds are set ', () => {
    // given
    component.secondsEnabled = true;
    const timeframe: Timeframe = {
      hours: 1,
      minutes: 15,
      seconds: 0,
      millis: undefined,
      micros: undefined,
    };

    // when
    component.changeTime(timeframe);

    // then
    expect(component.disabled).toEqual(false);
  });

  it('should be disabled = true if hours are not set', () => {
    // given
    component.secondsEnabled = true;
    const timeframe: Timeframe = {
      hours: undefined,
      minutes: 15,
      seconds: 0,
      millis: undefined,
      micros: undefined,
    };

    // when
    component.changeTime(timeframe);

    // then
    expect(component.disabled).toEqual(true);
  });

  it('should be disabled = true if minutes are not set', () => {
    // given
    component.secondsEnabled = true;
    const timeframe: Timeframe = {
      hours: 1,
      minutes: undefined,
      seconds: 0,
      millis: undefined,
      micros: undefined,
    };

    // when
    component.changeTime(timeframe);

    // then
    expect(component.disabled).toEqual(true);
  });

  it('should be disabled = false if hours, minutes and seconds are not set', () => {
    // given
    component.secondsEnabled = true;
    const timeframe: Timeframe = {
      hours: undefined,
      minutes: undefined,
      seconds: undefined,
      millis: undefined,
      micros: undefined,
    };

    // when
    component.changeTime(timeframe);

    // then
    expect(component.disabled).toEqual(false);
  });

  it('should emit 00:00:00.000 as time when time is not set', () => {
    // given
    const spy = jest.spyOn(component.selectedDateTime, 'emit');
    component.secondsEnabled = true;
    const date = moment();
    const expectedDate = date.hours(0).minutes(0).seconds(0).milliseconds(0);

    // when
    component.changeDate(date.toDate());
    component.setDateTime();

    // then
    expect(spy).toHaveBeenCalledWith(expectedDate.toISOString());
  });

  it('should emit the selected dateTime with seconds not set if not enabled', () => {
    // given
    const spy = jest.spyOn(component.selectedDateTime, 'emit');
    component.secondsEnabled = false;
    const momentDate = moment();
    const timeframe: Timeframe = {
      hours: 1,
      minutes: 15,
      seconds: undefined,
      millis: undefined,
      micros: undefined,
    };

    // when
    setDateAndTime(momentDate, timeframe);
    component.setDateTime();

    // then
    momentDate.hours(1).minutes(15).seconds(0).milliseconds(0);
    expect(spy).toHaveBeenCalledWith(momentDate.toISOString());
  });

  it('should emit the selected dateTime with seconds set if enabled', () => {
    // given
    const spy = jest.spyOn(component.selectedDateTime, 'emit');
    component.secondsEnabled = true;
    const momentDate = moment();
    const timeframe: Timeframe = {
      hours: 1,
      minutes: 15,
      seconds: 30,
      millis: undefined,
      micros: undefined,
    };

    // when
    setDateAndTime(momentDate, timeframe);
    component.setDateTime();

    // then
    momentDate.hours(1).minutes(15).seconds(30).milliseconds(0);
    expect(spy).toHaveBeenCalledWith(momentDate.toISOString());
  });

  function setDateAndTime(momentDate: Moment, timeframe: Timeframe): void {
    momentDate.date(1).month(1).year(2021).hours(0).minutes(0).seconds(0);
    component.changeDate(momentDate.toDate());
    component.changeTime(timeframe);
  }
});
