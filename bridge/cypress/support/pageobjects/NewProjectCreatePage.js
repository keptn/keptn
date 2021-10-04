/// <reference types="cypress" />

class NewProjectCreatePage {
  inputProjectName(PROJECT_NAME) {
    cy.get('input[id="projectNameInput"]').type(PROJECT_NAME);
    return this;
  }

  inputGitUrl(GIT_URL) {
    cy.get('input[placeholder="https://git-repo.com/repo.git"]').type(GIT_URL);
    return this;
  }

  inputGitUsername(GIT_USERNAME) {
    cy.get('input[placeholder="Username"]').type(GIT_USERNAME);
    return this;
  }

  inputGitToken(GIT_TOKEN) {
    cy.get('input[placeholder="Token"]').type(GIT_TOKEN);
    return this;
  }

  clickCreateProject() {
    cy.get('span.dt-button-label').contains('Create project').click();
  }
}

export default NewProjectCreatePage;
