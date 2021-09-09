describe('User Menu', () => {
  it('should open the user menu validate api token', async () => {
    cy.visit('/');

    cy.get('//*[@uitestid="keptn-nav-userMenu"]').click();

    cy.get('//*[@uitestid="keptn-nav-copyKeptnAuthCommand"]/ktb-copy-to-clipboard/div/div[3]/button').click();

    cy.screenshot('user-menu-open-with-api-command-revealed.png');
  });
});
