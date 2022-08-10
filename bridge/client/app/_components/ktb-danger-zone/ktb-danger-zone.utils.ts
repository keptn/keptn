import { Observable, of } from 'rxjs';
import { DeleteResult, DeletionProgressEvent } from '../../_interfaces/delete';

export function handleDeletionError(errorPrefix: string): (err: Error) => Observable<DeletionProgressEvent> {
  return (err: Error): Observable<DeletionProgressEvent> => {
    const deletionError = `${errorPrefix} could not be deleted: ${err.message}`;
    return of({
      error: deletionError,
      isInProgress: false,
      result: DeleteResult.ERROR,
    });
  };
}
