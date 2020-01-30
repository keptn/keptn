import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'atob',
})
export class AtobPipe implements PipeTransform {
  transform(text: string) {
    return atob(text);
  }
}
