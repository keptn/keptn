/// <reference types="cypress" />
import BasePage from '../support/pageobjects/BasePage';
import EnvironmentPage from '../support/pageobjects/EnvironmentPage';

describe('Project delete test', () => {
  it('test', () => {
    const basePage = new BasePage();
    const envPage = new EnvironmentPage();

    cy.fixture('get.project.json').as('initProjectJSON');
    cy.fixture('metadata.json').as('initmetadata');

    cy.intercept('GET', 'api/v1/metadata', { fixture: 'metadata.json' }).as('metadataCmpl');
    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      fixture: 'get.project.json',
    }).as('initProjects');
    cy.intercept('GET', 'api/controlPlane/v1/sequence/dynatrace?pageSize=5', { fixture: 'project.sequences.json' });

    cy.intercept('GET', 'api/project/dynatrace?approval=true&remediation=true', {
      statusCode: 200,
    });

    cy.intercept('DELETE', '/api/controlPlane/v1/project/dynatrace', {
      statusCode: 200,
    });

    cy.intercept('POST', 'api/hasUnreadUniformRegistrationLogs', {
      statusCode: 200,
    });

    cy.visit('/');
    cy.wait('@metadataCmpl');
    basePage.declineAutomaticUpdate().clickProjectTile('dynatrace');
    basePage.gotoSettingsPage().clickDeleteProjectButton().typeProjectNameToDelete('dynatrace').submitDelete();
  });
});
