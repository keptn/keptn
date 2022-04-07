import BasePage from '../support/pageobjects/BasePage';

describe('Test OAuth', () => {
  const basePage = new BasePage();

  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata' });
  });

  it('should redirect to login page for 401 response', () => {
    cy.intercept('/api/bridgeInfo', { statusCode: 401, headers: { 'keptn-auth-type': 'OAUTH' } });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      statusCode: 401,
      headers: { 'keptn-auth-type': 'OAUTH' },
    });

    cy.visit('/');
    basePage.notificationInfoVisible('Login required. Redirecting to login.');

    cy.url().should('include', '/oauth/login');
  });

  it('should show message for 403 response', () => {
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { statusCode: 403 });

    cy.visit('/');
    basePage.notificationErrorVisible('You do not have the permissions to perform this action.');
  });
});
