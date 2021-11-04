/// <reference types="cypress" />

class NewProjectCreatePage {
  inputProjectName(PROJECT_NAME: string): this {
    cy.get('input[id="projectNameInput"]').type(PROJECT_NAME);
    return this;
  }

  inputGitUrl(GIT_URL: string): this {
    cy.get('input[placeholder="https://git-repo.com/repo.git"]').type(GIT_URL);
    return this;
  }

  inputGitUsername(GIT_USERNAME: string): this {
    cy.get('input[placeholder="Username"]').type(GIT_USERNAME);
    return this;
  }

  inputGitToken(GIT_TOKEN: string): this {
    cy.get('input[placeholder="Token"]').type(GIT_TOKEN);
    return this;
  }

  clickCreateProject(): void {
    cy.get('span.dt-button-label').contains('Create project').forceClick();
  }
}

export default NewProjectCreatePage;
