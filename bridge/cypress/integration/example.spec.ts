describe('My First Test', () => {
  it('Does not do much!', () => {
    cy.visit('/');

    cy.get('dt-top-bar-navigation-item[uitestid="keptn-nav-projectMenu"]').click();
    cy.get('dt-tile-title[uitestid="keptn-project-tile-title"]').should('contain.text', 'dynatrace');
    cy.get('#projectSelect').click();
    cy.get('dt-option').contains('dynatrace').click();

    // go to services page and validate
    cy.get('button[aria-label="Open services view"]').click();
    cy.get('dt-info-group-title.dt-info-group-title > p').should('contain.text', 'Services');

    // go to sequences page and validate
    cy.get('button[aria-label="Open sequences view"]').click();
    cy.get('div[class="dt-quick-filter-group-headline"]');

    // go to uniform page and validate
    cy.get('button[aria-label="Open uniform view"]').click();
    cy.get('button[aria-label="Open uniform services"]').should('be.visible');

    // go to integration page and validate the date
    cy.get('button[aria-label="Open integration view"]').click();
    cy.get('#integration-container > p').should('contain.text', 'Integrate Keptn using the');
  });
});
