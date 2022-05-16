import 'jest-preset-angular/setup-jest';
import '@angular/localize/init';
import { Injectable } from '@angular/core';
import { WindowConfig } from './client/environments/environment.dynamic';
import { TestUtils } from './client/app/_utils/test.utils';

@Injectable()
class AppInitServiceMock {
  public init(): Promise<null | WindowConfig> {
    return Promise.resolve(null);
  }
}

Object.defineProperty(window, 'DragEvent', {
  value: class DragEvent {},
});

TestUtils.mockWindowMatchMedia();

jest.setTimeout(30000);
jest.mock('./client/app/_services/app.init', () => ({ AppInitService: AppInitServiceMock }));
