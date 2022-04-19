import { Pipe, PipeTransform } from '@angular/core';
import moment from 'moment';

@Pipe({
  name: 'formatDate',
})
export class FormatDatePipe implements PipeTransform {
  transform(value: string | undefined, format?: string): string | undefined {
    if (!value) return value;
    return moment(value).format(format || 'YYYY-MM-DD HH:mm');
  }
}
