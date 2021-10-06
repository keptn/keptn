describe('Create Project', () => {
  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata' });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { body: null });
    cy.visit('/create/project');
  });

  it('should create a new project', () => {
    cy.screenshot('create-project-page');
  });
});
