/// <reference types="cypress" />

import SettingsPage from './SettingsPage';
import ServicesPage from './ServicesPage';
import { interceptProjectBoard } from '../intercept';

enum View {
  SERVICE_VIEW = 'service-view',
  ENVIRONMENT_VIEW = 'environment-view',
  SEQUENCE_VIEW = 'sequence-view',
  SETTINGS_VIEW = 'settings-view',
}

export class ProjectBoardPage {
  NAVIGATION_MENU_LOCATOR = 'button[aria-label="Open page_pattern view"]';

  public interceptDeepLinks(): this {
    const mockedKeptnContext = '62cca6f3-dc54-4df6-a04c-6ffc894a4b5e';
    cy.intercept(`/api/mongodb-datastore/event?keptnContext=${mockedKeptnContext}`, {
      fixture: 'sequence.traces.mock.json',
    });
    return this;
  }

  public interceptError(projectName: string): this {
    interceptProjectBoard();
    cy.intercept(`/api/project/${projectName}?approval=true&remediation=true`, { forceNetworkError: true }).as(
      'project'
    );
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
    return this.assertMenuItemSelected(View.ENVIRONMENT_VIEW);
  }

  public assertOnlyServicesViewSelected(): this {
    return this.assertMenuItemSelected(View.SERVICE_VIEW);
  }

  public assertOnlySequencesViewSelected(): this {
    return this.assertMenuItemSelected(View.SEQUENCE_VIEW);
  }

  public assertOnlySettingsViewSelected(): this {
    return this.assertMenuItemSelected(View.SETTINGS_VIEW);
  }

  private assertMenuItemSelected(view: View): this {
    return this.assertServicesViewSelected(view === View.SERVICE_VIEW)
      .assertEnvironmentViewSelected(view === View.ENVIRONMENT_VIEW)
      .assertSequencesViewSelected(view === View.SEQUENCE_VIEW)
      .assertSettingsViewSelected(view === View.SETTINGS_VIEW);
  }

  private assertMenuSelected(selector: string, status: boolean): this {
    cy.byTestId(selector).should(status ? 'have.class' : 'not.have.class', 'active');
    return this;
  }
}
