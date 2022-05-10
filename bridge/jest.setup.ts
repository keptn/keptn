import 'jest-preset-angular';
import 'jest-preset-angular/setup-jest';
import { TestUtils } from './client/app/_utils/test.utils';

Object.defineProperty(window, 'DragEvent', {
  value: class DragEvent {},
});

TestUtils.mockWindowMatchMedia();

jest.setTimeout(30000);
