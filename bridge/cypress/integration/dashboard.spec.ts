import DashboardPage from '../support/pageobjects/DashboardPage';

import * as projectsResponse from '../fixtures/projects.mock.json';

describe('Bridge Dashboard', () => {
  const dashboardPage = new DashboardPage();

  beforeEach(() => {
    dashboardPage.intercept();
  });

  it('should load and show projects', () => {
    dashboardPage.visit().assertProjects(projectsResponse.projects);
  });

  it('should load also if version.json is not available', () => {
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfoVersionCheck.mock' }).as('bridgeInfo');
    cy.intercept('/api/version.json', { statusCode: 500 }).as('version.json');

    // versioncheck was accepted by user
    localStorage.setItem('keptn_versioncheck', JSON.stringify({ enabled: true, time: 1647880061049 }));

    dashboardPage.visit();
    cy.wait('@projects', { timeout: 15000 });
    cy.wait('@sequences', { timeout: 15000 });
    dashboardPage.assertProjects(projectsResponse.projects);
  });
});
