/// <reference types="cypress" />

class EnvironmentPage {
  STAGE_HEADER_LOC = 'div > h2';

  public visit(project: string): this {
    cy.visit(`/project/${project}`);
    return this;
  }

  public clickCreateService(stage: string): this {
    cy.get('ktb-selectable-tile h2')
      .contains(stage)
      .parentsUntil('ktb-selectable-tile')
      .find('ktb-no-service-info a')
      .click();
    return this;
  }
}
export default EnvironmentPage;
