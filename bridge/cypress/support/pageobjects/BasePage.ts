/// <reference types="cypress" />

import SettingsPage from './SettingsPage';
import NewProjectCreatePage from './NewProjectCreatePage';
import EnvironmentPage from './EnvironmentPage';
import ServicesPage from './ServicesPage';
import Chainable = Cypress.Chainable;

class BasePage {
  NAVIGATION_MENU_LOCATOR: string;
  PROJECT_TILE_LOCATOR: string;
  constructor() {
    const NAVIGATION_MENU_LOCATOR = "button[aria-label='Open page_pattern view']";
    const PROJECT_TILE_LOCATOR = "dt-tile[id='proj_patten']";
    this.NAVIGATION_MENU_LOCATOR = NAVIGATION_MENU_LOCATOR;
    this.PROJECT_TILE_LOCATOR = PROJECT_TILE_LOCATOR;
  }

  login(username: string, password: string): void {
    cy.get('#email_verify').type(username);
    cy.get('#next_button').click();
    cy.get('#password_login').type(password);
    cy.get('#no_captcha_submit').click();
  }

  // go to Uniform page
  goToUniformPage(): SettingsPage {
    return this.gotoSettingsPage().goToUniformPage();
  }

  // go to Services page
  goToServicesPage(): ServicesPage {
    cy.get(this.NAVIGATION_MENU_LOCATOR.replace('page_pattern', 'services')).click();
    return new ServicesPage();
  }

  // go to Settings page
  gotoSettingsPage(): SettingsPage {
    cy.get(this.NAVIGATION_MENU_LOCATOR.replace('page_pattern', 'settings')).click();
    return new SettingsPage();
  }

  selectProject(projectName: string | number | RegExp): void {
    cy.get('dt-tile-title[uitestid="keptn-project-tile-title"]')
      .should('contain.text', projectName)
      .get('#projectSelect')
      .click()
      .get('dt-option')
      .contains(projectName)
      .click();
  }

  selectProjectThroughHeader(projectName: string): void {
    cy.byTestId('keptn-nav-projectSelect')
      .click()
      .get('.cdk-overlay-container dt-option')
      .contains(projectName)
      .click();
  }

  clickProjectTile(projectName: string): EnvironmentPage {
    cy.wait(500).get(this.PROJECT_TILE_LOCATOR.replace('proj_patten', projectName)).click();
    return new EnvironmentPage();
  }

  clickCreateNewProjectButton(): NewProjectCreatePage {
    cy.get('.dt-button-primary > span.dt-button-label').contains('Create a new project').click();
    return new NewProjectCreatePage();
  }

  clickMainHeaderKeptn(): void {
    cy.get('.brand > p').contains('Keptn').click();
  }

  chooseProjectFromHeaderMenu(projectName: string): this {
    cy.get('dt-select[aria-label="Choose project"]').click();
    cy.get('dt-option[id^="dt-option"]').contains(projectName).click();
    return this;
  }

  notificationSuccessVisible(text?: string): Chainable<JQuery<HTMLElement>> {
    return this.checkNotification('dt-alert.dt-alert-success', text);
  }

  notificationErrorVisible(text?: string): Chainable<JQuery<HTMLElement>> {
    let element = cy.get('.dt-alert-icon-container dt-icon');
    if (text) {
      element = element.contains(text);
    }
    return element.get('.dt-alert-icon').should('be.visible');
  }

  notificationWarningVisible(text?: string): Chainable<JQuery<HTMLElement>> {
    return this.checkNotification('dt-alert.dt-alert-warning', text);
  }

  notificationInfoVisible(text?: string): Chainable<JQuery<HTMLElement>> {
    return this.checkNotification('dt-alert.dt-alert-info', text);
  }

  private checkNotification(selector: string, text?: string): Chainable<JQuery<HTMLElement>> {
    let element = cy.get(selector);
    if (text) {
      element = element.contains(text);
    }
    return element.should('be.visible');
  }
}

export default BasePage;
