import 'jest-preset-angular';
import 'jest-preset-angular/setup-jest';
import fetchMock from 'jest-fetch-mock';
import { TestUtils } from './client/app/_utils/test.utils';

fetchMock.enableMocks();

Object.defineProperty(window, 'DragEvent', {
  value: class DragEvent {},
});

TestUtils.mockWindowMatchMedia();

jest.setTimeout(30000);
