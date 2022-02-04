import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Timeframe } from '../../_models/timeframe';
import { FormControl, FormGroup } from '@angular/forms';

@Component({
  selector: 'ktb-time-input',
  templateUrl: './ktb-time-input.component.html',
  styleUrls: ['./ktb-time-input.component.scss'],
})
export class KtbTimeInputComponent {
  @Input() required: boolean | undefined;
  @Input() label = '';
  @Input() secondsEnabled = true;
  @Input() millisEnabled = true;
  @Input() microsEnabled = true;
  @Output()
  timeChanged: EventEmitter<Timeframe> = new EventEmitter<Timeframe>();

  public isFocused = false;

  public console = console;

  timeForm = new FormGroup({
    hours: new FormControl(''),
    minutes: new FormControl(''),
    seconds: new FormControl(''),
    millis: new FormControl(''),
    micros: new FormControl(''),
  });

  public validateInput(formControlName: string, min: number, max: number): void {
    if (this.timeForm.controls[formControlName].value) {
      let val = this.timeForm.controls[formControlName].value;
      val = Math.round(val);
      if (val < min) val = min;
      if (val > max) val = max;
      this.timeForm.controls[formControlName].setValue(val);
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
      hours: this.timeForm.controls.hours.value ? this.timeForm.controls.hours.value : undefined,
      minutes: this.timeForm.controls.minutes.value ? this.timeForm.controls.minutes.value : undefined,
      seconds: this.timeForm.controls.seconds.value ? this.timeForm.controls.seconds.value : undefined,
      millis: this.timeForm.controls.millis.value ? this.timeForm.controls.millis.value : undefined,
      micros: this.timeForm.controls.micros.value ? this.timeForm.controls.micros.value : undefined,
    });
  }
}
