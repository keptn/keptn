describe('Bridge Navigation', () => {
  beforeEach(() => {
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {fixture: 'projects.mock.json'});
  });


  it('should click on the project dropdown', () => {
    cy.visit('/');
    cy.wait(1000);

    cy.get('[uitestid="keptn-nav-projectMenu"]').click();
    cy.screenshot('project-menu-open');

    cy.fixture('projects.mock.json').then((projectFixture) => {
      cy.get('dt-option').contains(projectFixture.projects[0].projectName).click();
    });
    cy.screenshot('project-menu-clicked');
  });
});
