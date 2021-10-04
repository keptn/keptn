/// <reference types="cypress" />

class SecretsPage {
  clickAddSecret() {
    cy.get('.dt-button > .dt-button-label').contains('Add secret').click();
  }

  addSecret(SECRET_NAME, SECRET_KEY, SECRET_VALUE) {
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
  }

  deleteSecret(SECRET_NAME) {
    cy.get('dt-row.dt-row > dt-cell').contains(SECRET_NAME).next().children('button').click();
    cy.get('span.dt-button-label').contains('Delete').click();
  }
}

export default SecretsPage;
