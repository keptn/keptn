describe('Loading Bridge Example', () => {
  it('should match the title', () => {
    cy.visit('/');
    cy.title().should('eq', 'keptn');
    cy.screenshot('entry-page.png');
  });
});
