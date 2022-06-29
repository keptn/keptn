import { fromKeptInfo } from './root.store.selectors';
import { Features, LoadingState } from '../store';
import { initialRootState } from './root.store.reducer';

describe('Root State Selectors', () => {
  it('should select KeptnInfo', () => {
    const result = fromKeptInfo({
      [Features.ROOT]: initialRootState,
    });
    expect(result).toEqual({ call: LoadingState.INIT });
  });
});
