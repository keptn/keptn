import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'toType',
})
export class ToType implements PipeTransform {
  // cls argument is unused, but is serving the main goal: the type gets inferred from the constructor.
  transform<T>(value: unknown, _cls: new (...args: unknown[]) => T): T {
    return value as T;
  }
}
