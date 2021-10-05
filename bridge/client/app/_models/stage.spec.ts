import { Stage } from './stage';

describe('Stage', () => {
  it('should create instances from json', () => {
    const stage: Stage = Stage.fromJSON({ services: [] });

    expect(stage).toBeInstanceOf(Stage);
  });
});
