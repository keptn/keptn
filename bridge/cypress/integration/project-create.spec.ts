/// <reference types="cypress" />
import { ProjectBoardPage } from '../support/pageobjects/ProjectBoardPage';
import NewProjectCreatePage from '../support/pageobjects/NewProjectCreatePage';

describe('Create new project test', () => {
  it('test new project create', () => {
    const basePage = new ProjectBoardPage();
    const newProjectCreatePage = new NewProjectCreatePage();
    const GIT_USERNAME = 'carpe-github-username';
    const PROJECT_NAME = 'test-project-bycypress-001';
    const GIT_REMOTE_URL = 'https://git-repo.com';
    const GIT_TOKEN = 'testtoken';
    cy.fixture('get.project.json').as('initProjectJSON');
    cy.fixture('metadata.json').as('initmetadata');

    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
    cy.intercept('GET', 'api/v1/metadata', { fixture: 'metadata' }).as('metadataCmpl');
    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      fixture: 'get.project.json',
    }).as('initProjects');

    cy.intercept('POST', 'api/controlPlane/v1/project', {
      statusCode: 200,
      body: {},
    }).as('createProjectUrl');

    // eslint-disable-next-line promise/catch-or-return,promise/always-return
    cy.fixture('create.project.request.body').then((reqBody) => {
      cy.intercept('POST', 'api/controlPlane/v1/project', (req) => {
        expect(req.body).to.deep.equal(reqBody);
        return { status: 200, body: {} };
      });
    });

    cy.intercept('GET', 'api/controlPlane/v1/sequence/dynatrace?pageSize=5', { fixture: 'project.sequences' });

    cy.intercept('GET', 'api/project/testproject?approval=true&remediation=true', {
      statusCode: 200,
    }).as('projectApproval');

    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      statusCode: 200,
    }).as('disableUpstreamSync');

    cy.visit('/');
    cy.wait('@metadataCmpl');
    basePage
      .clickCreateNewProjectButton()
      .typeProjectName(PROJECT_NAME)
      .typeGitUrl(GIT_REMOTE_URL)
      .typeGitUsername(GIT_USERNAME)
      .typeGitToken(GIT_TOKEN)
      .setShipyardFile()
      .clickCreateProject();
  });
});
