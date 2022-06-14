/// <reference types="cypress" />
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

  public login(username: string, password: string): void {
    cy.get('#email_verify').type(username);
    cy.get('#next_button').click();
    cy.get('#password_login').type(password);
    cy.get('#no_captcha_submit').click();
  }

  public selectProjectThroughHeader(projectName: string): this {
    cy.byTestId('keptn-nav-projectSelect').dtSelect(projectName);
    return this;
  }

  public selectProject(projectName: string): this {
    cy.byTestId('keptn-project-tile-title').should('contain.text', projectName);
    return this.selectProjectThroughHeader(projectName);
  }

  public clickMainHeaderKeptn(): void {
    cy.get('.brand > p').contains('Keptn').click();
  }

  public chooseProjectFromHeaderMenu(projectName: string): this {
    cy.get('dt-select[aria-label="Choose project"]').click();
    cy.get('dt-option[id^="dt-option"]').contains(projectName).click();
    return this;
  }

  public notificationSuccessVisible(text?: string): Chainable<JQuery<HTMLElement>> {
    return this.checkNotification('ktb-notification-success', text);
  }

  public notificationErrorVisible(text?: string): Chainable<JQuery<HTMLElement>> {
    return this.checkNotification('ktb-notification-error', text);
  }

  public notificationWarningVisible(text?: string): Chainable<JQuery<HTMLElement>> {
    return this.checkNotification('ktb-notification-warning', text);
  }

  public notificationInfoVisible(text?: string): Chainable<JQuery<HTMLElement>> {
    return this.checkNotification('ktb-notification-info', text);
  }

  private checkNotification(selector: string, text?: string): Chainable<JQuery<HTMLElement>> {
    let element = cy.byTestId(selector);
    if (text) {
      element = element.contains(text);
    }
    return element.should('be.visible');
  }

  public clickOpenUserMenu(): this {
    cy.byTestId('keptn-user-menu-button').click();
    return this;
  }

  public assertAuthCommandCopyToClipboardValue(value: string): this {
    cy.byTestId('keptn-auth-command-ctc')
      .find('div.dt-copy-to-clipboard-input-wrap input')
      .invoke('val')
      .should('eq', value);
    return this;
  }
}

export default BasePage;
