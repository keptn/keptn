describe('User Menu', () => {
  it('should open the user menu validate api token', () => {
    cy.visit('/');

    cy.xpath('//*[@uitestid="keptn-nav-userMenu"]').click();

    cy.xpath('//*[@uitestid="keptn-nav-copyKeptnAuthCommand"]/ktb-copy-to-clipboard/div/div[3]/button').click();

    cy.screenshot('user-menu-open-with-api-command-revealed');
  });
});
