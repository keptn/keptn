describe('Loading Bridge Example', () => {
  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata' });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { body: null });
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
  });
  
  it('should match the title', () => {
    cy.visit('/');
    cy.title().should('eq', 'keptn');
    cy.screenshot('entry-page');
  });
});
