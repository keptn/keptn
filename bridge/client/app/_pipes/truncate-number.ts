import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'truncateNumber',
})
export class TruncateNumberPipe implements PipeTransform {
  transform(value: number | undefined, decimals: number): number | undefined {
    if (value) {
      return Math.trunc(value * Math.pow(10, decimals)) / Math.pow(10, decimals);
    }
    return value;
  }
}
