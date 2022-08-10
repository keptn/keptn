import { handleDeletionError } from './ktb-danger-zone.utils';
import { DeleteResult } from '../../_interfaces/delete';

describe('KtbDangerZoneUtils', () => {
  it('handleDeletionError should return DeletionProgressEvent', (done) => {
    handleDeletionError('Project')({ name: 'my error', message: 'some message' }).subscribe((deletionProgressEvent) => {
      expect(deletionProgressEvent).toStrictEqual({
        error: 'Project could not be deleted: some message',
        isInProgress: false,
        result: DeleteResult.ERROR,
      });
      done();
    });
  });
});
