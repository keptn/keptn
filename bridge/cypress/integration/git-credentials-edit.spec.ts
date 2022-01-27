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

    // eslint-disable-next-line promise/catch-or-return,promise/always-return
    cy.fixture('change.credentials.payload.json').then((reqBody) => {
      cy.intercept('PUT', 'api/controlPlane/v1/project', (req) => {
        expect(req.body).to.deep.equal(reqBody);
        return { status: 200 };
      });
    });

    cy.intercept('PUT', 'api/controlPlane/v1/project', {
      statusCode: 200,
    }).as('changeGitCredentials');

    cy.intercept('POST', 'api/hasUnreadUniformRegistrationLogs', {
      statusCode: 200,
    }).as('hasUnreadUniformRegistrationLogs');

    cy.intercept('GET', 'api/project/dynatrace?approval=true&remediation=true', {
      statusCode: 200,
    }).as('getApproval');

    cy.intercept('GET', 'api/project/dynatrace', {
      statusCode: 200,
    });

    cy.visit('/');
    cy.wait('@metadataCmpl');
    basePage.selectProject(DYNATRACE_PROJECT);
    basePage
      .gotoSettingsPage()
      .inputGitUrl(GIT_URL)
      .inputGitUsername(GIT_USER)
      .inputGitToken(GIT_TOKEN)
      .clickSaveGitUpstream();
  });

  it('Prevent data loss if git crendetials not saved before navigation', () => {
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

    // eslint-disable-next-line promise/catch-or-return,promise/always-return
    cy.fixture('change.credentials.payload.json').then((reqBody) => {
      cy.intercept('PUT', 'api/controlPlane/v1/project', (req) => {
        expect(req.body).to.deep.equal(reqBody);
        return { status: 200 };
      });
    });

    cy.intercept('PUT', 'api/controlPlane/v1/project', {
      statusCode: 200,
    }).as('changeGitCredentials');

    cy.intercept('POST', 'api/hasUnreadUniformRegistrationLogs', {
      statusCode: 200,
    }).as('hasUnreadUniformRegistrationLogs');

    cy.intercept('GET', 'api/project/dynatrace?approval=true&remediation=true', {
      statusCode: 200,
    }).as('getApproval');

    cy.intercept('GET', 'api/project/dynatrace', {
      statusCode: 200,
    });

    cy.intercept('GET', '/api/project/dynatrace/serviceStates', { statusCode: 200 });

    cy.visit('/');
    cy.wait('@metadataCmpl');
    basePage.selectProject(DYNATRACE_PROJECT);

    const settingsPage = basePage.gotoSettingsPage();
    settingsPage.inputGitUrl(GIT_URL).inputGitUsername(GIT_USER).inputGitToken(GIT_TOKEN);
    basePage.goToServicesPage();
    settingsPage.clickSaveChanges();
  });
});
