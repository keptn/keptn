/// <reference types="cypress" />

import BasePage from '../support/pageobjects/BasePage';
import { ServicesSettingsPage } from '../support/pageobjects/ServicesSettingsPage';
import NewProjectCreatePage from '../support/pageobjects/NewProjectCreatePage';

describe('Test notifications', () => {
  beforeEach(() => {
    cy.fixture('get.projects.json').as('initProjectJSON');
    cy.fixture('metadata.json').as('initmetadata');

    cy.intercept('GET', 'api/v1/metadata', { fixture: 'metadata.json' }).as('metadataCmpl');
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });

    cy.intercept('POST', 'api/hasUnreadUniformRegistrationLogs', {
      statusCode: 200,
    }).as('hasUnreadUniformRegistrationLogs');

    cy.intercept('GET', 'api/project/dynatrace?approval=true&remediation=true', {
      statusCode: 200,
      fixture: 'get.approval.json',
    }).as('getApproval');

    cy.intercept('GET', 'api/project/dynatrace', {
      statusCode: 200,
      fixture: 'get.approval.json',
    }).as('getApproval');

    cy.intercept('GET', '/api/project/dynatrace/services', {
      statusCode: 200,
      body: ['serviceA'],
    });
    cy.intercept('POST', '/api/controlPlane/v1/project/dynatrace/service', {
      statusCode: 200,
      body: {},
    });
    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      fixture: 'get.project.json',
    }).as('initProjects');
  });

  it('should test notification fade out', () => {
    const basePage = new BasePage();
    cy.visit('/project/dynatrace/settings/services/create');
    showSuccess();

    basePage
      .notificationSuccessVisible()
      .wait(6000)
      .should('not.have.css', 'opacity', '1')
      .trigger('mouseover')
      .should('have.css', 'opacity', '1')
      .wait(8200)
      .should('be.visible')
      .trigger('mouseleave')
      .wait(8200)
      .should('not.exist');
  });

  it('should test notification close', () => {
    const basePage = new BasePage();
    cy.visit('/project/dynatrace/settings/services/create');
    showSuccess();

    const notification = basePage.notificationSuccessVisible();
    cy.get('button[title="close"]').click();
    notification.should('not.exist');
  });

  it('should not show the same notifications', () => {
    cy.visit('/project/dynatrace/settings/services').byTestId('keptn-create-service-button').click();
    showSuccess();
    cy.byTestId('keptn-create-service-button').click();
    showSuccess();
    cy.get('dt-alert.dt-alert-success').should('have.length', 1);
  });

  it('should show two notifications', () => {
    createProject();
    showSuccess();
    cy.get('dt-alert.dt-alert-success').should('have.length', 2);
  });

  function createProject(): void {
    const basePage = new BasePage();
    const newProjectCreatePage = new NewProjectCreatePage();
    const GIT_USERNAME = 'carpe-github-username';
    const PROJECT_NAME = 'dynatrace';
    const GIT_REMOTE_URL = 'https://git-repo.com';
    const GIT_TOKEN = 'testtoken';
    // don't return any projects, to make project creation easier (no duplicates)
    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      body: {
        nextPageKey: '0',
        projects: [],
      },
    }).as('initProjects');
    cy.intercept('POST', 'api/controlPlane/v1/project', {
      statusCode: 200,
      body: {},
    }).as('createProjectUrl');

    cy.visit('/');
    basePage
      .clickCreateNewProjectButton()
      .typeProjectName(PROJECT_NAME)
      .typeGitUrl(GIT_REMOTE_URL)
      .typeGitUsername(GIT_USERNAME)
      .typeGitToken(GIT_TOKEN)
      .setShipyardFile();

    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { fixture: 'get.projects' });
    newProjectCreatePage.clickCreateProject();
    cy.get('dt-alert a').click();
  }

  function showSuccess(): void {
    const serviceSettings = new ServicesSettingsPage();
    serviceSettings.createService('my-new-service');
  }
});
