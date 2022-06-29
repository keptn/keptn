import { keptnInfoLoaded, loadRootState } from './root.store.actions';
import { KeptnInfo } from '../../_models/keptn-info';

describe('root.store.actions', () => {
  it('should return actions', () => {
    expect(loadRootState().type).toBe('[Root] Load Root State');

    const keptnInfo = <KeptnInfo>{};
    expect(keptnInfoLoaded({ keptnInfo }).type).toBe('[Root] KeptnInfo Loaded');
    expect(keptnInfoLoaded({ keptnInfo }).keptnInfo).toBe(keptnInfo);
  });
});
