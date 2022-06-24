describe('evaluations', () => {
  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata.mock' });
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
    cy.intercept('/api/project/sockshop?approval=true&remediation=true', { fixture: 'project.mock' }).as('project');
    cy.intercept('/api/hasUnreadUniformRegistrationLogs', { body: false });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { fixture: 'projects.mock' });
  });

  it('should load the heatmap with sli breakdown in service screen', () => {
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
      fixture: 'get.sockshop.service.carts.evaluations.mock.json',
    });

    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging');
    cy.byTestId('keptn-service-view-service-carts').should('exist');
    cy.byTestId('keptn-evaluation-details-chartHeatmap').should('exist');
  });

  it('should truncate score to 2 decimals', () => {
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
      fixture: 'get.sockshop.service.carts.evaluations.mock.json',
    });

    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging');

    cy.byTestId('keptn-evaluation-details-scoreInfo').should('have.text', '33.99');
  });
});
