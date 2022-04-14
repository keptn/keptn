/// <reference types="cypress" />
import { ProjectBoardPage } from '../support/pageobjects/ProjectBoardPage';
import DashboardPage from '../support/pageobjects/DashboardPage';

describe('Project delete test', () => {
  const projectBoardPage = new ProjectBoardPage();
  const dashboardPage = new DashboardPage();

  it('test', () => {
    cy.fixture('get.project.json').as('initProjectJSON');
    cy.fixture('metadata.json').as('initmetadata');

    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
    cy.intercept('GET', 'api/v1/metadata', { fixture: 'metadata.json' }).as('metadataCmpl');
    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      fixture: 'get.project.json',
    }).as('initProjects');
    cy.intercept('GET', 'api/controlPlane/v1/sequence/dynatrace?pageSize=5', { fixture: 'project.sequences.json' });

    cy.intercept('GET', 'api/project/dynatrace?approval=true&remediation=true', {
      statusCode: 200,
      fixture: 'get.approval.json',
    });

    cy.intercept('DELETE', '/api/controlPlane/v1/project/dynatrace', {
      statusCode: 200,
    });

    cy.intercept('POST', 'api/hasUnreadUniformRegistrationLogs', {
      statusCode: 200,
    });

    cy.intercept('GET', 'api/project/dynatrace', {
      statusCode: 200,
      fixture: 'get.approval.json',
    });

    cy.visit('/');
    cy.wait('@metadataCmpl');
    cy.wait('@initProjects');
    cy.wait(500);
    dashboardPage.clickProjectTile('dynatrace');
    projectBoardPage.gotoSettingsPage().clickDeleteProjectButton().typeProjectNameToDelete('dynatrace').submitDelete();
  });
});
