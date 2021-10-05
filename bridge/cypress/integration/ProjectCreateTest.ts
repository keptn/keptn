/// <reference types="cypress" />
import BasePage from '../support/pageobjects/BasePage';
import NewProjectCreatePage from '../support/pageobjects/NewProjectCreatePage';

describe('Create new project test', () => {
  it('test new project create', () => {
    const basePage = new BasePage();
    const newProjectCreatePage = new NewProjectCreatePage();
    const GIT_USERNAME = 'carpe-github-username';
    const PROJECT_NAME = 'test-project-bycypress-001';
    const GIT_REMOTE_URL = 'https://git-repo.com';
    const GIT_TOKEN = 'testtoken';
    cy.fixture('get.project.json').as('initProjectJSON');
    cy.fixture('metadata.json').as('initmetadata');

    cy.intercept('GET', 'api/v1/metadata', { fixture: 'metadata.json' }).as('metadataCmpl');
    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      fixture: 'get.project.json',
    }).as('initProjects');
    cy.intercept('POST', 'api/controlPlane/v1/project', {
      statusCode: 200,
      body: '',
    }).as('createProjectUrl');
    cy.intercept('GET', 'api/controlPlane/v1/sequence/dynatrace?pageSize=5', { fixture: 'project.sequences.json' });
    cy.intercept('GET', 'api/project/testproject?approval=true&remediation=true', {
      statusCode: 200,
    }).as('projectApproval');

    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      statusCode: 200,
    }).as('disableUpstreamSync');

    cy.visit('/');
    cy.wait('@metadataCmpl');
    basePage.declineAutomaticUpdate();
    basePage
      .clickCreatNewProjectButton()
      .inputProjectName(PROJECT_NAME)
      .inputGitUrl(GIT_REMOTE_URL)
      .inputGitUsername(GIT_USERNAME)
      .inputGitToken(GIT_TOKEN);

    cy.get('input[id="shipyard-file-input"]').attachFile('shipyard.yaml');

    newProjectCreatePage.clickCreateProject();

    cy.wait('@createProjectUrl', { timeout: 20000 });

    return cy.fixture('create.project.request.body.json').then((createProjjson) => {
      cy.get('@createProjectUrl').its('request.body').should('deep.equal', createProjjson);
      return;
    });

    /*cy.get('@createProjectUrl').its('request.body').should('deep.equal', {
      gitRemoteUrl: 'https://git-repo.com',
      gitToken: 'testtoken',
      gitUser: 'carpe-github-username',
      name: 'test-project-bycypress-001',
      shipyard:
        'YXBpVmVyc2lvbjogInNwZWMua2VwdG4uc2gvMC4yLjAiDQpraW5kOiAiU2hpcHlhcmQiDQptZXRhZGF0YToNCiAgbmFtZTogInNoaXB5YXJkLXF1YWxpdHktZ2F0ZXMiDQpzcGVjOg0KICBzdGFnZXM6DQogICAgLSBuYW1lOiAicXVhbGl0eS1nYXRlLXRlc3Qi',
    });*/
  });
});
