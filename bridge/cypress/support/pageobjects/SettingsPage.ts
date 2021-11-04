/// <reference types="cypress" />

class SettingsPage {
  inputGitUrl(GIT_URL: string): this {
    cy.get('input[formcontrolname="gitUrl"]').type(GIT_URL);
    return this;
  }

  inputGitUsername(GIT_USERNAME: string): this {
    cy.get('input[formcontrolname="gitUser"]').type(GIT_USERNAME);
    return this;
  }

  inputGitToken(GIT_TOKEN: string): this {
    cy.get('input[formcontrolname="gitToken"]').type(GIT_TOKEN);
    return this;
  }

  clickSaveChanges(): this {
    cy.get('.dt-button-primary > span.dt-button-label').contains('Save changes').forceClick();
    return this;
  }

  getErrorMessageText(): unknown {
    return cy.get('.small');
  }

  clickDeleteProjectButton(): this {
    cy.get('span.dt-button-label').contains('Delete this project').forceClick();
    return this;
  }

  typeProjectNameToDelete(projectName: string): this {
    const projectInputLoc = 'input[placeholder=proj_pattern]';
    cy.get(projectInputLoc.replace('proj_pattern', projectName)).forceClick().type(projectName);
    return this;
  }

  submitDelete(): void {
    cy.get('span.dt-button-label').contains('I understand the consequences, delete this project').forceClick();
  }
}

export default SettingsPage;
