import { ServerErrors } from '../../../client/app/_models/server-error';

export class ErrorPage {
  public visit(status?: ServerErrors | string, queryParams?: Record<string, string>): this {
    let query;
    if (status !== undefined) {
      query = `?status=${status}`;
    } else {
      query = '';
    }
    cy.visit(`/error${query}`, {
      qs: queryParams,
    });
    return this;
  }

  public assertHeaderText(text: string): this {
    cy.get('h2').should('have.text', text);
    return this;
  }

  public assertMessage(text: string): this {
    cy.get('.text p').should('have.text', text);
    return this;
  }

  public locationExists(status: boolean): this {
    cy.byTestId('ktb-location-link').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public clickLocation(): this {
    cy.byTestId('ktb-location-link').click();
    return this;
  }

  public isInternalError(): this {
    return this.assertHeaderText('Internal error')
      .assertMessage('Error while handling the redirect. Please retry and check whether the problem still exists.')
      .locationExists(true);
  }

  public isInsufficientPermissionError(): this {
    return this.assertHeaderText('Permission denied')
      .assertMessage('User is not allowed to access the instance.')
      .locationExists(false);
  }

  public isTraceError(status: boolean): this {
    cy.byTestId('ktb-error-trace').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertTraceErrorKeptnContext(keptnContext: string): this {
    return this.assertHeaderText(` Traces for ${keptnContext} not found `);
  }

  public assertTraceErrorWithoutKeptnContext(): this {
    return this.assertHeaderText('No traces found');
  }
}
