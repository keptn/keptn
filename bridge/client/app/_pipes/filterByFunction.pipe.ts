import { PipeTransform, Pipe } from '@angular/core';

@Pipe({
  name: 'filterByFunction',
  pure: false
})
export class FilterByFunctionPipe implements PipeTransform {
  transform(items: any[], callback: (item: any) => boolean): any {
    if (!items || !callback) {
      return items;
    }
    return items.filter(item => callback(item));
  }
}
