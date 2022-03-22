import DashboardPage from '../support/pageobjects/DashboardPage';

import * as projectsResponse from '../fixtures/projects.mock.json';

describe('Bridge Dashboard', () => {
  const dashboardPage = new DashboardPage();

  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata' });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { fixture: 'projects.mock' }).as(
      'projects'
    );
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=5', { fixture: 'sequences.sockshop' }).as(
      'sequences'
    );
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
  });

  it('should load and show projects', () => {
    dashboardPage.visitDashboard();
    dashboardPage.assertProjects(projectsResponse.projects);
  });

  it('should load also if version.json is not available', () => {
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfoVersionCheck.mock' }).as('bridgeInfo');
    cy.intercept('/api/v1/metadata', { fixture: 'metadata' });
    cy.intercept('/api/version.json', { statusCode: 500 }).as('version.json');

    // versioncheck was accepted by user
    localStorage.setItem('keptn_versioncheck', JSON.stringify({ enabled: true, time: 1647880061049 }));

    dashboardPage.visitDashboard();
    cy.visit('/');
    cy.wait('@projects', { timeout: 15000 });
    cy.wait('@sequences', { timeout: 15000 });
    dashboardPage.assertProjects(projectsResponse.projects);
  });
});
