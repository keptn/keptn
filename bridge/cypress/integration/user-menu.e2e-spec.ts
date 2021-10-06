describe('User Menu', () => {
  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata' });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { body: null });
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
  });

  it('should open the user menu validate api token', () => {
    cy.visit('/');

    cy.xpath('//*[@uitestid="keptn-nav-userMenu"]').click();

    cy.xpath('//*[@uitestid="keptn-nav-copyKeptnAuthCommand"]/ktb-copy-to-clipboard/div/div[3]/button').click();

    cy.screenshot('user-menu-open-with-api-command-revealed');
  });
});
