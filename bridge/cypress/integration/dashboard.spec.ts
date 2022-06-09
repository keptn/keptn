import DashboardPage from '../support/pageobjects/DashboardPage';

import * as projectsResponse from '../fixtures/projects.mock.json';
import BasePage from '../support/pageobjects/BasePage';

describe('Bridge Dashboard', () => {
  const dashboardPage = new DashboardPage();

  beforeEach(() => {
    dashboardPage.intercept();
  });

  it('should load and show projects', () => {
    dashboardPage.visit().assertProjects(projectsResponse.projects);
  });

  it('should trigger loadProjects once per dashboard visit', () => {
    dashboardPage.visit().clickCreateNewProjectButton(); // 1 call
    const basePage = new BasePage();
    basePage.clickMainHeaderKeptn(); // 1 call
    cy.wait('@projects').get('@projects.all').should('have.length', 2);
  });

  it('should load also if version.json is not available', () => {
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfoVersionCheck.mock' }).as('bridgeInfo');
    cy.intercept('/api/version.json', { statusCode: 500 }).as('version.json');

    // versioncheck was accepted by user
    localStorage.setItem('keptn_versioncheck', JSON.stringify({ enabled: true, time: 1647880061049 }));

    dashboardPage.visit(false);
    // wait for all retries
    for (let i = 0; i < 4; ++i) {
      cy.wait('@version.json');
    }
    cy.wait('@projects').wait('@sequences');
    dashboardPage.assertProjects(projectsResponse.projects);
  });
});
