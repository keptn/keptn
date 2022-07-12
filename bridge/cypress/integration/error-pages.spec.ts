import { ErrorPage } from '../support/pageobjects/ErrorPage';
import { ServerErrors } from '../../client/app/_models/server-error';

describe('Error pages', () => {
  const errorPage = new ErrorPage();

  it('should show trace error page with keptnContext', () => {
    const keptnContext = 'myContext';
    errorPage.visit(ServerErrors.TRACE, { keptnContext }).isTraceError(true).assertTraceErrorKeptnContext(keptnContext);
  });

  it('should show trace error page without keptnContext', () => {
    errorPage.visit(ServerErrors.TRACE).isTraceError(true).assertTraceErrorWithoutKeptnContext();
  });
});
