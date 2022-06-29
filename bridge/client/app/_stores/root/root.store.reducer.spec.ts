import { initialRootState, rootStoreReducer } from './root.store.reducer';
import { loadRootState } from './root.store.actions';
import { LoadingState } from '../store';

describe('Root Store Reducer', () => {
  it('should return the previous state', () => {
    const actual = rootStoreReducer(initialRootState, loadRootState());
    expect(actual.keptInfo).toStrictEqual({ call: LoadingState.LOADING, data: undefined });
  });
});
