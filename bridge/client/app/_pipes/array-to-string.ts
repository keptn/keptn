import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'arrayAsString',
})
export class ArrayToStringPipe implements PipeTransform {
  transform(value: string[] | number[] | undefined): string {
    if (value) {
      return value.join(', ');
    }
    return '';
  }
}
