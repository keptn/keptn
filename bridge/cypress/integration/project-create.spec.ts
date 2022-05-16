/// <reference types="cypress" />
import DashboardPage from '../support/pageobjects/DashboardPage';
import NewProjectCreatePage from '../support/pageobjects/NewProjectCreatePage';

describe('Create new project test', () => {
  it('test new project create', () => {
    const dashboardPage = new DashboardPage();
    const createProjectPage = new NewProjectCreatePage();
    const GIT_USERNAME = 'carpe-github-username';
    const PROJECT_NAME = 'sockshop';
    const GIT_REMOTE_URL = 'https://git-repo.com';
    const GIT_TOKEN = 'testtoken';

    dashboardPage.intercept();
    createProjectPage.intercept().interceptSettings();

    cy.intercept('/api/project/sockshop', { fixture: 'project.mock' });
    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      body: {
        projects: [],
      },
    });

    // eslint-disable-next-line promise/catch-or-return,promise/always-return
    cy.fixture('create.project.request.body').then((reqBody) => {
      cy.intercept('POST', 'api/controlPlane/v1/project', (req) => {
        expect(req.body).to.deep.equal(reqBody);
        return { status: 200, body: {} };
      });
    });

    dashboardPage
      .visit()
      .clickCreateNewProjectButton()
      .typeProjectName(PROJECT_NAME)
      .typeGitUrl(GIT_REMOTE_URL)
      .typeGitUsername(GIT_USERNAME)
      .typeGitToken(GIT_TOKEN)
      .setShipyardFile()
      .clickCreateProject();
  });
});
