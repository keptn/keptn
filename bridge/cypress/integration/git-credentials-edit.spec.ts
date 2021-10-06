/// <reference types="cypress" />

import BasePage from '../support/pageobjects/BasePage';

describe('Changing git credentials', () => {
  it('The test changes git credentials and makes sure they changed successfully', () => {
    const basePage = new BasePage();
    const DYNATRACE_PROJECT = 'dynatrace';
    const GIT_URL = 'https://git-repo.com';
    const GIT_USER = 'test-username';
    const GIT_TOKEN = 'test-token!§$%&/()=';

    cy.fixture('get.project.json').as('initProjectJSON');
    cy.fixture('metadata.json').as('initmetadata');
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });

    cy.intercept('GET', 'api/v1/metadata', { fixture: 'metadata.json' }).as('metadataCmpl');
    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      fixture: 'get.project.json',
    }).as('initProjects');
    cy.intercept('GET', 'api/controlPlane/v1/sequence/dynatrace?pageSize=5', { fixture: 'project.sequences.json' });

    cy.intercept('PUT', 'api/controlPlane/v1/project', {
      statusCode: 200,
    }).as('changeGitCredentials');

    cy.intercept('POST', 'api/hasUnreadUniformRegistrationLogs', {
      statusCode: 200,
    }).as('hasUnreadUniformRegistrationLogs');

    cy.intercept('GET', 'api/project/dynatrace?approval=true&remediation=true', {
      statusCode: 200,
    }).as('getApproval');

    cy.visit('/');
    cy.wait('@metadataCmpl');
    basePage.declineAutomaticUpdate();
    basePage.selectProject(DYNATRACE_PROJECT);
    basePage
      .gotoSettingsPage()
      .inputGitUrl(GIT_URL)
      .inputGitUsername(GIT_USER)
      .inputGitToken(GIT_TOKEN)
      .clickSaveChanges();

    return cy.fixture('change.credentials.payload.json').then((json) => {
      cy.get('@changeGitCredentials').its('request.body').should('deep.equal', json);
      return null;
    });
  });
});
