/// <reference types="cypress" />
import EnvironmentPage from 'cypress/support/pageobjects/EnvironmentPage';
import ServicesPage from 'cypress/support/pageobjects/ServicesPage';
import SettingsPage from 'cypress/support/pageobjects/SettingsPage';
import UnifromPage from 'cypress/support/pageobjects/UniformPage';
import BasePage from '../support/pageobjects/BasePage';

describe('verify Page data update when project is switched', () => {
  it('testing the data change/update ', () => {
    const basePage = new BasePage();
    const envPage = new EnvironmentPage();
    const servicePage = new ServicesPage();
    const uniformPage = new UnifromPage();
    const settingsPage = new SettingsPage();
    const PROJECT_DYNATRACE = 'dynatrace';
    const PROJECT_TEST = 'testproject';
    const QG_DYNATRACE = 'quality-gate';
    const QG_TEST_PROJ = 'newproj-quality-gate';
    const SERVICE_DYNATRACE = 'items';
    const SERVICE_TEST_PROJ = 'newService';
    const UNIFORM_REMEDIATION = 'remediation-service';
    const UNIFORM_APPROVAL_SERVICE = 'approval-service';
    const UNIFORM_DYNATRACE_SERVICE = 'dynatrace-service';

    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
    cy.intercept('GET', 'api/v1/metadata', { fixture: 'metadata.json' }).as('metadataCmpl');
    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      statusCode: 200,
      fixture: 'multi.projects.mock.json',
    }).as('initProjects');

    cy.intercept('GET', 'api/controlPlane/v1/sequence/*', { statusCode: 200 });

    cy.intercept('POST', 'api/hasUnreadUniformRegistrationLogs', {
      statusCode: 200,
    }).as('hasUnreadUniformRegistrationLogs');

    cy.intercept('GET', 'api/project/dynatrace?approval=true&remediation=true', {
      statusCode: 200,
      fixture: 'get.approval.json',
    }).as('getApproval');

    cy.intercept('GET', 'api/project/testproject?approval=true&remediation=true', {
      statusCode: 200,
      fixture: 'get.approval.testproj.json',
    }).as('getApproval');

    cy.intercept('POST', 'api/uniform/registration', {
      statusCode: 200,
      fixture: 'get.uniform.json',
    }).as('postUniform');

    cy.visit('/');
    cy.wait('@metadataCmpl');
    basePage.clickProjectTile(PROJECT_DYNATRACE);
    cy.get(envPage.STAGE_HEADER_LOC).contains(QG_DYNATRACE);

    basePage.chooseProjectFromHeaderMenu(PROJECT_TEST);
    cy.get(envPage.STAGE_HEADER_LOC).contains(QG_TEST_PROJ);

    basePage.goToServicesPage();
    cy.get(servicePage.SERVICE_PANEL_TEXT_LOC).contains(SERVICE_TEST_PROJ);

    basePage.chooseProjectFromHeaderMenu(PROJECT_DYNATRACE);
    cy.get(servicePage.SERVICE_PANEL_TEXT_LOC).contains(SERVICE_DYNATRACE);

    basePage.goToUniformPage();
    cy.get(uniformPage.UNIFORM_NAME_LOC).contains(UNIFORM_APPROVAL_SERVICE);
    cy.get(uniformPage.UNIFORM_NAME_LOC).contains(UNIFORM_REMEDIATION);
    cy.get(uniformPage.UNIFORM_NAME_LOC).contains(UNIFORM_DYNATRACE_SERVICE);

    basePage.chooseProjectFromHeaderMenu(PROJECT_TEST);
    cy.get(uniformPage.UNIFORM_NAME_LOC).contains(UNIFORM_APPROVAL_SERVICE);
    cy.get(uniformPage.UNIFORM_NAME_LOC).contains(UNIFORM_REMEDIATION);
    cy.get(uniformPage.UNIFORM_NAME_LOC).contains(UNIFORM_DYNATRACE_SERVICE);

    basePage.gotoSettingsPage();
    cy.get(settingsPage.GIT_USER_LOC).should('have.value', 'carpe-github-testproj');
    cy.get(settingsPage.GIT_URL_INPUT_LOC).should(
      'have.value',
      'https://github.com/carpe-github/claus-tenant-dev-testproj.git'
    );

    basePage.chooseProjectFromHeaderMenu(PROJECT_DYNATRACE);
    cy.get(settingsPage.GIT_USER_LOC).should('have.value', 'carpe-github');
    cy.get(settingsPage.GIT_URL_INPUT_LOC).should('have.value', 'https://github.com/carpe-github/claus-tenant-dev.git');
  });
});
