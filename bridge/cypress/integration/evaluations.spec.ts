import { HeatmapComponentPage } from '../support/pageobjects/HeatmapComponentPage';

const heatmap = new HeatmapComponentPage();

describe('evaluations', () => {
  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata.mock' });
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfoEnableD3Heatmap.mock.json' });
    cy.intercept('/api/project/sockshop?approval=true&remediation=true', { fixture: 'project.mock' }).as('project');
    cy.intercept('/api/hasUnreadUniformRegistrationLogs', { body: false });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { fixture: 'projects.mock' });
  });

  xit('should load the heatmap with sli breakdown in service screen', () => {
    cy.intercept('GET', '/api/project/sockshop/serviceStates', {
      statusCode: 200,
      fixture: 'get.sockshop.service.states.mock.json',
    }).as('serviceStates');
    cy.intercept(
      'GET',
      '/api/project/sockshop/deployment/da740469-9920-4e0c-b304-0fd4b18d17c2?includeRemediations=false',
      {
        statusCode: 200,
        fixture: 'get.sockshop.service.carts.deployment.mock.json',
      }
    ).as('ServiceDeployment');
    cy.intercept('GET', 'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
      statusCode: 200,
      fixture: 'get.sockshop.service.carts.evaluations.mock.json',
    });

    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging');
    cy.byTestId('keptn-service-view-service-carts').should('exist');
    cy.get('ktb-heatmap').should('exist');
  });

  xit('should truncate score to 2 decimals', () => {
    cy.intercept('GET', '/api/project/sockshop/serviceStates', {
      statusCode: 200,
      fixture: 'get.sockshop.service.states.mock.json',
    }).as('serviceStates');
    cy.intercept(
      'GET',
      '/api/project/sockshop/deployment/da740469-9920-4e0c-b304-0fd4b18d17c2?includeRemediations=false',
      {
        statusCode: 200,
        fixture: 'get.sockshop.service.carts.deployment.mock.json',
      }
    ).as('ServiceDeployment');
    cy.intercept('GET', 'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
      statusCode: 200,
      fixture: 'get.sockshop.service.carts.evaluations.mock.json',
    });

    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging');

    cy.byTestId('keptn-evaluation-details-scoreInfo').should('have.text', '33.99');
  });

  it('should show key sli info', () => {
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
      fixture: 'get.sockshop.service.carts.evaluations.keysli.mock.json',
    });

    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging');
    cy.byTestId('keptn-service-view-service-carts').should('exist');
    cy.get('ktb-heatmap').should('exist');

    cy.byTestId('keptn-evaluation-details-scoreInfo').should('have.text', '50 < 75');
    cy.byTestId('keptn-evaluation-details-resultInfo').should('have.text', 'Result: fail');
    cy.byTestId('keptn-evaluation-details-keySliInfo').should('have.text', 'Key SLI: passed');

    heatmap.clickScore('52b4b2c7-fa49-41f3-9b5c-b9aea2370bb4');
    cy.byTestId('keptn-evaluation-details-scoreInfo').should('have.text', '75 >= 75');
    cy.byTestId('keptn-evaluation-details-resultInfo').should('have.text', 'Result: fail');
    cy.byTestId('keptn-evaluation-details-keySliInfo').should('have.text', 'Key SLI: failed');

    heatmap.clickScore('182d10b8-b68d-49d4-86cd-5521352d7a42');
    cy.byTestId('keptn-evaluation-details-scoreInfo').should('have.text', '100 >= 90');
    cy.byTestId('keptn-evaluation-details-resultInfo').should('have.text', 'Result: pass');
    cy.byTestId('keptn-evaluation-details-keySliInfo').should('have.text', 'Key SLI: passed');
  });
});
