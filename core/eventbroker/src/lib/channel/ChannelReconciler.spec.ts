import { expect } from 'chai';
import 'mocha';
import { ChannelReconciler } from './ChannelReconciler';

describe('ChannelReconciler', function () {
  this.timeout(0);
  let channelReconciler: ChannelReconciler;
  beforeEach(() => {
    channelReconciler = new ChannelReconciler();
  });

  it('Should return a channel URI', async () => {
    const channelName: string = 'new-artefact';
    const channelUri = await channelReconciler.resolveChannel(channelName);

    // expect(channelUri.indexOf(channelName) > -1).is.true; TODO: reactivate after PR has been merged
  });
  it('Should return an empty string if no channel can be found', async () => {
    const channelName: string = 'idontexist';
    const channelUri = await channelReconciler.resolveChannel(channelName);

    expect(channelUri === '').is.true;
  });
});
