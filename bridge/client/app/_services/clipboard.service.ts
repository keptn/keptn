import { Injectable } from '@angular/core';
import { DtToast } from '@dynatrace/barista-components/toast';
import { JsonSerializable } from '../_models/json-serializable';

@Injectable({
  providedIn: 'root',
})
export class ClipboardService {

  constructor(private readonly toast: DtToast) {
  }

  copy(serializable: JsonSerializable, label = 'value'): void {
    const value = this.stringify(serializable);
    if (navigator && 'clipboard' in navigator && typeof navigator.clipboard.writeText === 'function') {
      navigator.clipboard.writeText(value);
    } else {
      const textarea = document.createElement('textarea');
      textarea.value = value;
      textarea.setAttribute('readonly', '');
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand('copy');
      document.body.removeChild(textarea);
    }

    this.toast.create(`Copied ${label} to clipboard`);
  }

  stringify(value: JsonSerializable): string {
    if (value === undefined || value === null) {
      return '';
    } else if (typeof value === 'string') {
      return value;
    } else {
      return JSON.stringify(value, undefined, 2);
    }
  }

}
