import { Injectable } from '@angular/core';
import { DtToast } from '@dynatrace/barista-components/toast';
import { JsonSerializable } from '../_models/json-serializable';
import { Clipboard } from '@angular/cdk/clipboard';

@Injectable({
  providedIn: 'root',
})
export class ClipboardService {
  constructor(private readonly toast: DtToast, private readonly clipboard: Clipboard) {}

  copy(serializable: JsonSerializable, label = 'value'): void {
    const value = this.stringify(serializable);
    this.clipboard.copy(value);
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
