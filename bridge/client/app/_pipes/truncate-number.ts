import { Pipe, PipeTransform } from '@angular/core';
import { AppUtils } from '../_utils/app.utils';

@Pipe({
  name: 'truncateNumber',
})
export class TruncateNumberPipe implements PipeTransform {
  transform(value: number | undefined, decimals: number): number | undefined {
    if (value) {
      return AppUtils.truncateNumber(value, decimals);
    }
    return value;
  }
}
