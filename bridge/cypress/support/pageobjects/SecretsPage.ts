/// <reference types="cypress" />

class SecretsPage {
  clickAddSecret(): void {
    cy.get('.dt-button > .dt-button-label').contains('Add secret').click();
  }

  addSecret(SECRET_NAME: string, SECRET_KEY: string, SECRET_VALUE: string): void {
    this.clickAddSecret();
    cy.get('input[uitestid="keptn-secret-name-input"]')
      .type(SECRET_NAME)
      .get('input[uitestid="keptn-secret-key-input"]')
      .type(SECRET_KEY)
      .get('input[uitestid="keptn-secret-value-input"]')
      .type(SECRET_VALUE)
      .get('button > span')
      .contains('Add secret')
      .click();
  }

  deleteSecret(SECRET_NAME: string | number | RegExp): void {
    cy.get('dt-row.dt-row > dt-cell > p').contains(SECRET_NAME).parent().next().children('button').click();
    cy.get('span.dt-button-label').contains('Delete').click();
  }
}

export default SecretsPage;
