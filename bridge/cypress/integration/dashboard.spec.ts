import DashboardPage from '../support/pageobjects/DashboardPage';

import * as projectsResponse from '../fixtures/projects.mock.json';
import BasePage from '../support/pageobjects/BasePage';
import ProjectSettingsPage from '../support/pageobjects/ProjectSettingsPage';

describe('Bridge Dashboard', () => {
  const dashboardPage = new DashboardPage();

  beforeEach(() => {
    dashboardPage.intercept();
  });

  it('should load and show projects', () => {
    dashboardPage.visit().assertProjects(projectsResponse.projects);
  });

  it('should trigger loadProjects once per dashboard visit', () => {
    const basePage = new BasePage();
    const projectSettings = new ProjectSettingsPage();

    dashboardPage.visit().clickCreateNewProjectButton(); // 1 call
    projectSettings.waitForSettingsToBeVisible();
    basePage.clickMainHeaderKeptn(); // 1 call
    cy.wait('@projects').get('@projects.all').should('have.length', 2);
  });

  it('should use the AUTH_MSG as auth command when provided', () => {
    cy.intercept('/api/bridgeInfo*', { fixture: 'bridgeInfoAuthMsg.mock' }).as('bridgeInfo');
    dashboardPage.visit();
    const basePage = new BasePage();
    basePage.clickOpenUserMenu().assertAuthCommandCopyToClipboardValue('Hello handsome');
  });

  it('should show pause icon if sequence is paused', () => {
    dashboardPage.visit().assertPauseIconShown();
  });

  it('should show "Set the Git upstream of your project" message if Github remote URL is empty string', () => {
    dashboardPage.visit().assertEmptyGitRemoteUrl('my-error-project');
  });
});
