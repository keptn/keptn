/// <reference types="cypress" />

import SettingsPage from './SettingsPage';
import NewProjectCreatePage from './NewProjectCreatePage';
import EnvironmentPage from './EnvironmentPage';

class BasePage {
  constructor() {
    const NAVIGATION_MENU_LOCATOR = "button[aria-label='Open page_pattern view']";
    const PROJECT_TILE_LOCATOR = "dt-tile[id='proj_patten']";
    //     let NAVIGATION_MENU_LOCATOR = "//button[@aria-label='Open page_pattern view']"
    this.NAVIGATION_MENU_LOCATOR = NAVIGATION_MENU_LOCATOR;
    this.PROJECT_TILE_LOCATOR = PROJECT_TILE_LOCATOR;
  }

  login(username, password) {
    cy.get('#email_verify').type(username);
    cy.get('#next_button').click();
    cy.get('#password_login').type(password);
    cy.get('#no_captcha_submit').click();
  }

  // go to Uniform page

  goToUniformPage() {
    cy.get(this.NAVIGATION_MENU_LOCATOR.replace('page_pattern', 'uniform')).click();
    return this;
  }

  // go to Secrets page
  goToSecretsPage() {
    cy.get('[aria-label="Open uniform secrets"]').click();
  }

  // go to Settings page

  gotoSettingsPage() {
    cy.get(this.NAVIGATION_MENU_LOCATOR.replace('page_pattern', 'settings')).click();
    return new SettingsPage();
  }

  selectProject(projectName) {
    cy.get('dt-tile-title[uitestid="keptn-project-tile-title"]')
      .should('contain.text', projectName)
      .get('#projectSelect')
      .click()
      .get('dt-option')
      .contains(projectName)
      .click();
  }

  clickProjectTile(projectName) {
    cy.get(this.PROJECT_TILE_LOCATOR.replace('proj_patten', projectName)).click();
    return new EnvironmentPage();
  }

  declineAutomaticUpdate() {
    cy.get('.dt-button-secondary > span.dt-button-label').contains('Decline').click();
    return this;
  }

  clickCreatNewProjectButton() {
    cy.get('.dt-button-primary > span.dt-button-label').contains('Create a new project').click();
    return new NewProjectCreatePage();
  }

  clickMainHeaderKeptn() {
    cy.get('.brand > p').contains('keptn').click();
  }
}

export default BasePage;
