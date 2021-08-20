import { Root } from './root';
import { waitForAsync } from '@angular/core/testing';

describe('Root', () => {
  it('should create instances from json', waitForAsync(() => {
    const root: Root =  Root.fromJSON([]);

    expect(root).toBeInstanceOf(Root);
  }));
});
