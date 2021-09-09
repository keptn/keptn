describe('Create Project', () => {
  beforeEach(() => {
    cy.visit('/create/project');
  });

  it('should create a new project', async () => {
    cy.screenshot('create-project-page.png');
  });
});
