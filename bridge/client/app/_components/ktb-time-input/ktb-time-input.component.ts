import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Timeframe } from '../../_models/timeframe';
import { FormControl } from '@angular/forms';

@Component({
  selector: 'ktb-time-input',
  templateUrl: './ktb-time-input.component.html',
  styleUrls: ['./ktb-time-input.component.scss'],
})
export class KtbTimeInputComponent implements OnInit {
  @Input() required = false;
  @Input() label = '';
  @Input() hint = '';
  @Input() error = '';
  @Input() secondsEnabled = true;
  @Input() millisEnabled = true;
  @Input() microsEnabled = true;
  @Input() min?: Timeframe;
  @Input() max?: Timeframe;
  @Input() timeframe: Timeframe | undefined;

  @Output()
  timeChanged: EventEmitter<Timeframe> = new EventEmitter<Timeframe>();

  public isFocused = false;

  public hoursControl: FormControl = new FormControl();
  public minutesControl: FormControl = new FormControl();
  public secondsControl: FormControl = new FormControl();
  public millisControl: FormControl = new FormControl();
  public microsControl: FormControl = new FormControl();

  public timeControls: Record<string, FormControl> = {
    hours: this.hoursControl,
    minutes: this.minutesControl,
    seconds: this.secondsControl,
    millis: this.millisControl,
    micros: this.microsControl,
  };

  public ngOnInit(): void {
    if (this.timeframe) {
      this.hoursControl.setValue(this.timeframe.hours ?? null);
      this.minutesControl.setValue(this.timeframe.minutes ?? null);
      this.secondsControl.setValue(this.timeframe.seconds ?? null);
      this.millisControl.setValue(this.timeframe.millis ?? null);
      this.microsControl.setValue(this.timeframe.micros ?? null);
    }
  }

  public validateInput(formControlName: string, min = 0, max?: number): void {
    if (this.timeControls[formControlName].value) {
      let val = this.timeControls[formControlName].value;
      val = Math.round(val);
      if (val < min) val = min;
      if (max && val > max) val = max;
      this.timeControls[formControlName].setValue(val);
    } else {
      // To prevent from infinite 0 and dots, we have to set it manually.
      // In the control the value is already set to the right value, the input is not updated.
      this.timeControls[formControlName].setValue(this.timeControls[formControlName].value);
    }
    this.emitChangedValues();
  }

  public focusClick(event: Event, hourInput: HTMLInputElement): void {
    if (event.target instanceof HTMLFormElement) {
      hourInput.focus();
    }
  }

  private emitChangedValues(): void {
    this.timeChanged.emit({
      hours: this.hoursControl.value ?? undefined,
      minutes: this.minutesControl.value ?? undefined,
      seconds: this.secondsControl.value ?? undefined,
      millis: this.millisControl.value ?? undefined,
      micros: this.microsControl.value ?? undefined,
    });
  }
}
