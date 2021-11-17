import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'toDate',
})
export class ToDatePipe implements PipeTransform {
  transform(value: string | number): Date {
    return new Date(value);
  }
}
