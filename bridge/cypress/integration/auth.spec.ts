import BasePage from '../support/pageobjects/BasePage';

describe('Test auth errors', () => {
  const basePage = new BasePage();

  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata' });
  });

  it('should show error for 401 response', () => {
    cy.intercept('/api/bridgeInfo', { statusCode: 401 });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      statusCode: 401,
    });

    cy.visit('/');
    basePage.notificationErrorVisible('Could not authorize.');
  });

  it('should show error for invalid token', () => {
    cy.intercept('/api/bridgeInfo', {
      statusCode: 401,
      body: 'incorrect api key auth',
    });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      statusCode: 401,
      body: 'incorrect api key auth',
    });

    cy.visit('/');
    basePage.notificationErrorVisible('Could not authorize API token. Please check the configured API token.');
  });
});

describe('Test BASIC auth', () => {
  const basePage = new BasePage();

  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata' });
  });

  it('should show an error for 401 response', () => {
    cy.intercept('/api/bridgeInfo', { statusCode: 401, headers: { 'keptn-auth-type': 'BASIC' } });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      statusCode: 401,
      headers: { 'keptn-auth-type': 'BASIC' },
    });

    cy.visit('/');
    basePage.notificationErrorVisible('Login credentials invalid. Please check your provided username and password.');
  });
});

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

    cy.location('pathname').should('eq', '/oauth/login');
  });

  it('should show a message for 403 response', () => {
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { statusCode: 403 });

    cy.visit('/');
    basePage.notificationErrorVisible('You do not have the permissions to perform this action.');
  });
});
