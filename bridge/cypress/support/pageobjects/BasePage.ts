/// <reference types="cypress" />

import SettingsPage from './SettingsPage';
import NewProjectCreatePage from './NewProjectCreatePage';
import EnvironmentPage from './EnvironmentPage';
import ServicesPage from './ServicesPage';

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
    cy.get('#next_button').forceClick();
    cy.get('#password_login').type(password);
    cy.get('#no_captcha_submit').forceClick();
  }

  // go to Uniform page

  goToUniformPage(): this {
    cy.get(this.NAVIGATION_MENU_LOCATOR.replace('page_pattern', 'uniform')).forceClick();
    return this;
  }

  // go to Secrets page
  goToSecretsPage(): void {
    cy.get('[aria-label="Open uniform secrets"]').forceClick();
  }

  // go to Services page
  goToServicesPage(): ServicesPage {
    cy.get(this.NAVIGATION_MENU_LOCATOR.replace('page_pattern', 'services')).forceClick();
    return new ServicesPage();
  }

  // go to Settings page
  gotoSettingsPage(): SettingsPage {
    cy.get(this.NAVIGATION_MENU_LOCATOR.replace('page_pattern', 'settings')).forceClick();
    return new SettingsPage();
  }

  selectProject(projectName: string | number | RegExp): void {
    cy.get('dt-tile-title[uitestid="keptn-project-tile-title"]')
      .should('contain.text', projectName)
      .get('#projectSelect')
      .forceClick()
      .get('dt-option')
      .contains(projectName)
      .forceClick();
  }

  clickProjectTile(projectName: string): EnvironmentPage {
    cy.get(this.PROJECT_TILE_LOCATOR.replace('proj_patten', projectName)).forceClick();
    return new EnvironmentPage();
  }

  clickCreateNewProjectButton(): NewProjectCreatePage {
    cy.get('.dt-button-primary > span.dt-button-label').contains('Create a new project').forceClick();
    return new NewProjectCreatePage();
  }

  clickMainHeaderKeptn(): void {
    cy.get('.brand > p').contains('keptn').forceClick();
  }
}

export default BasePage;
