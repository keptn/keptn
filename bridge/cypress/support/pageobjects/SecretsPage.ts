/// <reference types="cypress" />

class SecretsPage {
  clickAddSecret(): void {
    cy.get('.dt-button > .dt-button-label').contains('Add secret').click();
    return;
  }

  addSecret(SECRET_NAME: string, SECRET_KEY: string, SECRET_VALUE: string): void {
    this.clickAddSecret();
    cy.get('input[placeholder="Secret name"]')
      .type(SECRET_NAME)
      .get('input[placeholder="Key"]')
      .type(SECRET_KEY)
      .get('input[placeholder="Value"]')
      .type(SECRET_VALUE)
      .get('button > span')
      .contains('Add secret')
      .click();
    return;
  }

  deleteSecret(SECRET_NAME: string | number | RegExp): void {
    cy.get('dt-row.dt-row > dt-cell').contains(SECRET_NAME).next().children('button').click();
    cy.get('span.dt-button-label').contains('Delete').click();
    return;
  }
}

export default SecretsPage;
