import { Stage } from './stage';
import { waitForAsync } from '@angular/core/testing';

describe('Stage', () => {
  it('should create instances from json', waitForAsync(() => {
    const stage: Stage =  Stage.fromJSON({services: []});

    expect(stage).toBeInstanceOf(Stage);
  }));
});
