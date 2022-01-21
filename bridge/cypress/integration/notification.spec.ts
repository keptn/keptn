/// <reference types="cypress" />

import BasePage from '../support/pageobjects/BasePage';
import { ServicesSettingsPage } from '../support/pageobjects/ServicesSettingsPage';

describe('Test notifications', () => {
  beforeEach(() => {
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
    cy.visit('/project/dynatrace/settings/services/create');
  });

  it('should test notification fade out', () => {
    const basePage = new BasePage();
    showSuccess();

    basePage
      .notificationSuccessVisible()
      .wait(2000)
      .should('not.have.css', 'opacity', '1')
      .trigger('mouseover')
      .should('not.have.css', 'opacity', '0')
      .wait(5000)
      .should('be.visible')
      .trigger('mouseleave')
      .wait(5000)
      .should('not.exist');
  });

  it('should test notification close', () => {
    const basePage = new BasePage();
    showSuccess();

    const notification = basePage.notificationSuccessVisible();
    cy.get('button[title="close"]').click();
    notification.should('not.exist');
  });

  function showSuccess(): void {
    const serviceSettings = new ServicesSettingsPage();

    serviceSettings.inputService('my-new-service');
    serviceSettings.createService();
  }
});
