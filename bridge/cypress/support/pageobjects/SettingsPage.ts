/// <reference types="cypress" />
type subSettings = 'Integrations' | 'Project' | 'Services' | 'Secrets' | 'Common use cases';

class SettingsPage {
  public goToUniformPage(): this {
    return this.goToSubSettings('Integrations');
  }

  // go to Secrets page
  public goToSecretsPage(): this {
    return this.goToSubSettings('Secrets');
  }

  private goToSubSettings(subPage: subSettings): this {
    cy.get('dt-menu-group button[role=menuitem]').contains(subPage).click();
    return this;
  }
}

export default SettingsPage;
