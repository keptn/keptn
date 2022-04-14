import { Directive, HostBinding, HostListener, Input } from '@angular/core';
import { NgControl } from '@angular/forms';

@Directive({
  selector: '[ktbIntegerInput]',
})
export class KtbIntegerInputDirective {
  private readonly defaultInvalidNumberInputs: string[] = [',', '.'];
  private invalidNumberInputs: string[] = [...this.defaultInvalidNumberInputs, '-'];
  private _allowNegative = false;

  @HostBinding('type') type = 'number';
  @HostBinding('min') get min(): string | undefined {
    return this.allowNegative ? undefined : '0';
  }

  @Input() get allowNegative(): boolean {
    return this._allowNegative;
  }
  set allowNegative(value: boolean) {
    this._allowNegative = value;
    this.invalidNumberInputs = value ? this.defaultInvalidNumberInputs : [...this.defaultInvalidNumberInputs, '-'];
  }

  constructor(private control: NgControl) {}

  @HostListener('keydown', ['$event']) onKeyDown(event$: KeyboardEvent): void {
    if (this.invalidNumberInputs.includes(event$.key)) {
      event$.preventDefault();
    }
  }

  // on paste can't be modified.
  // It can be set directly to the control but then selectionStart and selectionEnd is missing because of type="number"
  // then the value would be replaced instead of attached

  @HostListener('input') onInput(): void {
    const value = this.control.control?.value;
    if (value && this.invalidNumberInputs.some((separator: string) => value.includes(separator))) {
      const truncated = value.split(/[,.]/)[0];
      const parsed = truncated === '' ? '' : this.getNumber(truncated).toString();
      this.control.control?.setValue(parsed === 'NaN' ? '' : parsed);
    }
  }

  private getNumber(value: string): number {
    return this.allowNegative ? +value : Math.abs(+value);
  }
}
