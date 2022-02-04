import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Timeframe } from '../../_models/timeframe';
import { FormControl, FormGroup } from '@angular/forms';

@Component({
  selector: 'ktb-timeframe-input',
  templateUrl: './ktb-timeframe-input.component.html',
  styleUrls: ['./ktb-timeframe-input.component.scss'],
})
export class KtbTimeframeInputComponent {
  @Input() required: boolean | undefined;
  @Output() timeframe: EventEmitter<Timeframe> = new EventEmitter<Timeframe>();

  public isFocused = false;

  timeframeForm = new FormGroup({
    hours: new FormControl(''),
    minutes: new FormControl(''),
    seconds: new FormControl(''),
    millis: new FormControl(''),
    micros: new FormControl(''),
  });

  public emitChangedValues(): void {
    this.timeframe.emit({
      hours: this.timeframeForm.controls.hours.value ? this.timeframeForm.controls.hours.value : undefined,
      minutes: this.timeframeForm.controls.minutes.value ? this.timeframeForm.controls.minutes.value : undefined,
      seconds: this.timeframeForm.controls.seconds.value ? this.timeframeForm.controls.seconds.value : undefined,
      millis: this.timeframeForm.controls.millis.value ? this.timeframeForm.controls.millis.value : undefined,
      micros: this.timeframeForm.controls.micros.value ? this.timeframeForm.controls.micros.value : undefined,
    });
  }

  public validateInput(formControlName: string, min: number, max: number): void {
    if (this.timeframeForm.controls[formControlName].value) {
      let val = this.timeframeForm.controls[formControlName].value;
      val = Math.round(val);
      if (val < min) val = min;
      if (val > max) val = max;
      this.timeframeForm.controls[formControlName].setValue(val);
    }
    this.emitChangedValues();
  }
}
