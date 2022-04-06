/// <reference types="cypress" />

import SettingsPage from './SettingsPage';
import NewProjectCreatePage from './NewProjectCreatePage';
import EnvironmentPage from './EnvironmentPage';
import ServicesPage from './ServicesPage';
import Chainable = Cypress.Chainable;

export class ProjectBoardPage {
  NAVIGATION_MENU_LOCATOR = 'button[aria-label="Open page_pattern view"]';
  PROJECT_TILE_LOCATOR = 'dt-tile[id="proj_patten"]';

  public interceptDeepLinks(): this {
    const mockedKeptnContext = '62cca6f3-dc54-4df6-a04c-6ffc894a4b5e';
    cy.intercept(`/api/mongodb-datastore/event?keptnContext=${mockedKeptnContext}`, {
      fixture: 'sequence.traces.mock.json',
    });
    return this;
  }

  public login(username: string, password: string): void {
    cy.get('#email_verify').type(username);
    cy.get('#next_button').click();
    cy.get('#password_login').type(password);
    cy.get('#no_captcha_submit').click();
  }

  // go to Uniform page
  public goToUniformPage(): SettingsPage {
    return this.gotoSettingsPage().goToUniformPage();
  }

  // go to Services page
  public goToServicesPage(): ServicesPage {
    cy.get(this.NAVIGATION_MENU_LOCATOR.replace('page_pattern', 'services')).click();
    return new ServicesPage();
  }

  // go to Settings page
  public gotoSettingsPage(): SettingsPage {
    cy.get(this.NAVIGATION_MENU_LOCATOR.replace('page_pattern', 'settings')).click();
    return new SettingsPage();
  }

  public selectProject(projectName: string | number | RegExp): void {
    cy.get('dt-tile-title[uitestid="keptn-project-tile-title"]')
      .should('contain.text', projectName)
      .get('#projectSelect')
      .click()
      .get('dt-option')
      .contains(projectName)
      .click();
  }

  public selectProjectThroughHeader(projectName: string): void {
    cy.byTestId('keptn-nav-projectSelect')
      .click()
      .get('.cdk-overlay-container dt-option')
      .contains(projectName)
      .click();
  }

  public clickProjectTile(projectName: string): EnvironmentPage {
    cy.wait(500).get(this.PROJECT_TILE_LOCATOR.replace('proj_patten', projectName)).click();
    return new EnvironmentPage();
  }

  public clickCreateNewProjectButton(): NewProjectCreatePage {
    cy.get('.dt-button-primary > span.dt-button-label').contains('Create a new project').click();
    return new NewProjectCreatePage();
  }

  public clickMainHeaderKeptn(): void {
    cy.byTestId('ktb-header-title').click();
  }

  public chooseProjectFromHeaderMenu(projectName: string): this {
    cy.get('dt-select[aria-label="Choose project"]').click();
    cy.get('dt-option[id^="dt-option"]').contains(projectName).click();
    return this;
  }

  public notificationSuccessVisible(text?: string): Chainable<JQuery<HTMLElement>> {
    return this.checkNotification('dt-alert.dt-alert-success', text);
  }

  public notificationErrorVisible(text?: string): Chainable<JQuery<HTMLElement>> {
    let element = cy.get('.dt-alert-icon-container dt-icon');
    if (text) {
      element = element.contains(text);
    }
    return element.get('.dt-alert-icon').should('be.visible');
  }

  public notificationWarningVisible(text?: string): Chainable<JQuery<HTMLElement>> {
    return this.checkNotification('dt-alert.dt-alert-warning', text);
  }

  public notificationInfoVisible(text?: string): Chainable<JQuery<HTMLElement>> {
    return this.checkNotification('dt-alert.dt-alert-info', text);
  }

  private assertEnvironmentViewSelected(status: boolean): this {
    return this.assertMenuSelected('ktb-environment-menu-button', status);
  }

  private assertServicesViewSelected(status: boolean): this {
    return this.assertMenuSelected('ktb-services-menu-button', status);
  }

  private assertSequencesViewSelected(status: boolean): this {
    return this.assertMenuSelected('ktb-sequences-menu-button', status);
  }

  private assertSettingsViewSelected(status: boolean): this {
    return this.assertMenuSelected('ktb-settings-menu-button', status);
  }

  public assertOnlyEnvironmentViewSelected(): this {
    return this.assertMenuItemsSelected(true, false, false, false);
  }

  public assertOnlyServicesViewSelected(): this {
    return this.assertMenuItemsSelected(false, true, false, false);
  }

  public assertOnlySequencesViewSelected(): this {
    return this.assertMenuItemsSelected(false, false, true, false);
  }

  public assertOnlySettingsViewSelected(): this {
    return this.assertMenuItemsSelected(false, false, false, true);
  }

  private assertMenuItemsSelected(
    environmentView: boolean,
    servicesView: boolean,
    sequencesView: boolean,
    settingsView: boolean
  ): this {
    return this.assertServicesViewSelected(servicesView)
      .assertEnvironmentViewSelected(environmentView)
      .assertSequencesViewSelected(sequencesView)
      .assertSettingsViewSelected(settingsView);
  }

  private assertMenuSelected(selector: string, status: boolean): this {
    cy.byTestId(selector).should(status ? 'have.class' : 'not.have.class', 'active');
    return this;
  }

  private checkNotification(selector: string, text?: string): Chainable<JQuery<HTMLElement>> {
    let element = cy.get(selector);
    if (text) {
      element = element.contains(text);
    }
    return element.should('be.visible');
  }
}
