/// <reference types="cypress" />

import { ProjectBoardPage } from '../support/pageobjects/ProjectBoardPage';

describe('Test project navigation', () => {
  it('should correctly load service states of different service', () => {
    const basePage = new ProjectBoardPage();
    const secondProject = 'dynatrace2';
    const serviceNameOfSecondProject = 'serviceOfDynatrace2';

    cy.fixture('get.projects.json').as('initProjectJSON');
    cy.fixture('metadata.json').as('initmetadata');

    cy.intercept('GET', 'api/v1/metadata', { fixture: 'metadata.json' }).as('metadataCmpl');
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });

    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      fixture: 'get.projects.json',
    }).as('initProjects');

    cy.intercept('POST', 'api/hasUnreadUniformRegistrationLogs', {
      statusCode: 200,
    }).as('hasUnreadUniformRegistrationLogs');

    cy.intercept('GET', 'api/project/dynatrace?approval=true&remediation=true', {
      statusCode: 200,
      fixture: 'get.approval.json',
    }).as('getApproval');

    cy.intercept('GET', `api/project/${secondProject}?approval=true&remediation=true`, {
      statusCode: 200,
      fixture: 'get.approval.json',
    }).as('getApproval');

    cy.intercept('GET', '/api/project/dynatrace/serviceStates', {
      statusCode: 200,
      fixture: 'get.service.states.mock.json',
    });

    cy.intercept('GET', '/api/project/dynatrace2/serviceStates', {
      statusCode: 200,
      body: [
        {
          deploymentInformation: [],
          name: serviceNameOfSecondProject,
        },
      ],
    });

    cy.intercept('GET', `/api/project/dynatrace/deployment/*`, {
      statusCode: 200,
      fixture: 'get.service.deployment.mock.json',
    });

    cy.visit('/project/dynatrace/service');
    cy.wait(500).get('div.dt-info-group-content').get('h2').contains('health');
    basePage.selectProjectThroughHeader(secondProject);
    cy.wait(500).get('div.dt-info-group-content').get('h2').contains(serviceNameOfSecondProject);
  });
});
