import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'toType',
})
export class ToType implements PipeTransform {
  // cls argument is unused, but is serving the main goal: the type gets inferred from the constructor.
  // eslint-disable-next-line @typescript-eslint/no-explicit-any,@typescript-eslint/no-unused-vars
  transform<T>(value: unknown, _cls: new (...args: any[]) => T): T {
    return value as T;
  }
}
