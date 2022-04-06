/// <reference types="cypress" />

import { ProjectBoardPage } from '../support/pageobjects/ProjectBoardPage';

describe('Test Navigation Buttons In Evaluation Screen', () => {
  it('The test clicks on Navigation buttons and make sure the pages are open respectively ', () => {
    const basePage = new ProjectBoardPage();
    cy.fixture('metadata.json').as('initmetadata');

    cy.intercept('GET', 'api/v1/metadata', { fixture: 'metadata.json' }).as('metadataCmpl');
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });

    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      fixture: 'get.project.json',
    }).as('initProjects');

    cy.intercept('GET', 'api/controlPlane/v1/sequence/dynatrace?*', { fixture: 'project.sequences.json' }).as(
      'getSequence'
    );

    cy.intercept('PUT', 'api/controlPlane/v1/project', {
      statusCode: 200,
    }).as('changeGitCredentials');

    cy.intercept('POST', 'api/hasUnreadUniformRegistrationLogs', {
      statusCode: 200,
    }).as('hasUnreadUniformRegistrationLogs');

    cy.intercept('GET', 'api/project/dynatrace?approval=true&remediation=true', {
      statusCode: 200,
      fixture: 'get.approval.json',
    }).as('getApproval');

    cy.intercept('GET', 'api/controlPlane/v1/sequence/dynatrace?pageSize=100&name=remediation&state=triggered', {
      statusCode: 200,
      fixture: 'sequence.dynatrace.json',
    });

    cy.intercept('GET', 'api/mongodb-datastore/event?root=true&pageSize=1&project=dynatrace&*', {
      statusCode: 200,
      fixture: 'service/get.eval.data.json',
    }).as('getEventRoot');

    cy.intercept('GET', 'api/mongodb-datastore/event?keptnContext=*&project=dynatrace', {
      statusCode: 200,
      fixture: 'service/get.event2.data.json',
    }).as('getEventKeptnContextWithProject');

    cy.intercept('GET', 'api/mongodb-datastore/event?keptnContext=*', {
      statusCode: 200,
      fixture: 'service/get.event.keptn.context.json',
    }).as('getEventWithKeptnContext');

    cy.intercept('GET', 'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
      statusCode: 200,
      fixture: 'service/get.eval.data.json',
    }).as('getEventEvalFinished');

    cy.intercept(
      'GET',
      'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?filter=data.project:dynatrace*',
      {
        statusCode: 200,
        fixture: 'service/get.eval.data.json',
      }
    ).as('getEventEvalFinishedWithProject');

    cy.intercept('GET', '/api/project/dynatrace/serviceStates', {
      statusCode: 200,
      fixture: 'get.service.states.mock.json',
    });
    cy.intercept('GET', `/api/project/dynatrace/deployment/*`, {
      statusCode: 200,
      fixture: 'get.service.deployment.mock.json',
    });

    cy.intercept('GET', '/api/project/dynatrace/sequences/metadata', {
      body: { deployments: [], filter: { stages: [], services: [] } },
    });

    cy.visit('/').wait('@initProjects');
    basePage.clickProjectTile('dynatrace');

    basePage
      .goToServicesPage()
      .clickOnServicePanelByName('items')
      .clickOnServiceInnerPanelByName('items')
      .clickEvaluationBoardButton()
      .clickViewServiceDetails()
      .verifyCurrentOpenServiceNameEvaluationPanel('items')
      .clickEvaluationBoardButton()
      .clickViewSequenceDetails();

    cy.get('*[uitestid="keptn-sequence-view-roots"]');

    basePage
      .goToServicesPage()
      .clickOnServicePanelByName('items')
      .clickOnServiceInnerPanelByName('items')
      .clickEvaluationBoardButton()
      .clickGoBack()
      .verifyCurrentOpenServiceNameEvaluationPanel('items');
    basePage.clickMainHeaderKeptn();
  });
});
