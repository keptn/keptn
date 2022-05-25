describe('evaluation-heatmap', () => {
  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata.mock' });
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfoEnableD3Heatmap.mock.json' });
    cy.intercept('/api/project/sockshop?approval=true&remediation=true', { fixture: 'project.mock' });
    cy.intercept('/api/hasUnreadUniformRegistrationLogs', { body: false });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { fixture: 'projects.mock' });
    cy.intercept('GET', '/api/project/sockshop/serviceStates', {
      statusCode: 200,
      fixture: 'get.sockshop.service.states.mock.json',
    });
    cy.intercept('GET', '/api/project/sockshop/deployment/da740469-9920-4e0c-b304-0fd4b18d17c2', {
      statusCode: 200,
      fixture: 'get.sockshop.service.carts.deployment.mock.json',
    });
    cy.intercept('GET', 'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
      statusCode: 200,
      fixture: 'get.sockshop.service.carts.evaluations.heatmap.mock.json',
    });
  });
  it('should display ktb-heatmap if the feature flag is enabled', () => {
    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging');

    cy.get('ktb-heatmap').should('exist');
  });
  it('should be expandable and collapsable', () => {
    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging');

    cy.get('ktb-heatmap svg .y-axis-container .tick').should('have.length', 10);
    cy.get('ktb-heatmap button').click();
    cy.get('ktb-heatmap svg .y-axis-container .tick').should('have.length', 13);
    cy.get('ktb-heatmap button').click();
    cy.get('ktb-heatmap svg .y-axis-container .tick').should('have.length', 10);
  });
  it('should set correct color classes', () => {
    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging');

    cy.byTestId('ktb-heatmap-tile-182d10b8-b68d-49d4-86cd-5521352d7a42')
      .should('have.class', 'pass')
      .and('have.class', 'data-point');
    cy.byTestId('ktb-heatmap-tile-182d10b8-b68d-49d4-86cd-5521352d7a42')
      .should('have.class', 'warning')
      .and('have.class', 'data-point');
    cy.byTestId('ktb-heatmap-tile-52b4b2c7-fa49-41f3-9b5c-b9aea2370bb4')
      .should('have.class', 'fail')
      .and('have.class', 'data-point');
  });
  it('should have a primary and secondary highlight', () => {
    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging');

    cy.get('ktb-heatmap .highlight-primary').should('have.length', 1);
    cy.get('ktb-heatmap .highlight-secondary').should('have.length', 1);
  });
  it('should truncate long metric names', () => {
    const longName = 'A very long metric name so long it gets cut somewhere along the way';
    const shortName = 'A very long metric name ...';
    cy.contains('ktb-heatmap title', longName)
      .parent()
      .within(() => {
        cy.contains('text', shortName);
      });
  });
});
