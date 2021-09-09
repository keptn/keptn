describe('Bridge Navigation', () => {

  it('should click on the project dropdown', async () => {
    cy.visit('/');

    cy.get('//*[@id="projectSelect"]//div').click();
    cy.screenshot('project-menu-open.png');

    cy.get('dt-option-0').click();
    cy.screenshot('project-menu-clicked.png');
  });
});
