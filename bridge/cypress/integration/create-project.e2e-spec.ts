describe('Create Project', () => {
  beforeEach(() => {
    cy.visit('/create/project');
  });

  it('should create a new project', () => {
    cy.screenshot('create-project-page');
  });
});
