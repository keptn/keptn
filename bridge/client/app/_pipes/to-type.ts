import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'toType'
})
export class ToType implements PipeTransform {
  // clss argument is unused, but is serving the main goal: the type gets inferred from the constructor.
  // tslint:disable-next-line:no-any
  transform<T>(value: unknown, _clss: new (...args: any[]) => T): T {
    return value as T;
  }
}
