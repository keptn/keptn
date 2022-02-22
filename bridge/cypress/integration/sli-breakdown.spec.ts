describe('evaluations', () => {
  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata.mock' });
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
    cy.intercept('/api/project/sockshop?approval=true&remediation=true', { fixture: 'project.mock' });
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
    cy.byTestId('keptn-sli-breakdown').should('exist');
    cy.byTestId('keptn-sli-breakdown-row-go_routines').find('dt-cell').eq(1).should('have.text', 'go_routines');
    cy.byTestId('keptn-sli-breakdown-row-go_routines').find('dt-cell').eq(2).should('have.text', '88 (+1000%) ');
    cy.byTestId('keptn-sli-breakdown-row-go_routines').find('dt-cell').eq(3).should('have.text', '1');
    cy.byTestId('keptn-sli-breakdown-row-go_routines').find('dt-cell').eq(7).should('have.text', '33.33');

    cy.byTestId('keptn-sli-breakdown-row-request_throughput')
      .find('dt-cell')
      .eq(2)
      .find('.error')
      .should('have.length', 1);
    cy.byTestId('keptn-sli-breakdown-row-go_routines').find('dt-cell').eq(2).find('.error').should('have.length', 0);
  });

  it('should show more details when expanding sli breakdown in service screen', () => {
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

    cy.byTestId('keptn-sli-breakdown-row-go_routines').find('dt-cell').eq(2).should('have.text', '88 (+1000%) ');
    cy.byTestId('keptn-sli-breakdown-row-go_routines')
      .find('dt-cell')
      .eq(1)
      .should('not.contain.text', 'Absolute change:');
    cy.byTestId('keptn-sli-breakdown-row-go_routines')
      .find('dt-cell')
      .eq(1)
      .should('not.contain.text', 'Relative change:');
    cy.byTestId('keptn-sli-breakdown-row-go_routines')
      .find('dt-cell')
      .eq(1)
      .should('not.contain.text', 'Compared with:');

    cy.byTestId('keptn-sli-breakdown-row-go_routines').find('dt-cell').eq(0).find('button').click();

    cy.byTestId('keptn-sli-breakdown-row-go_routines').find('dt-cell').eq(2).should('have.text', '88+80+1000% 8');
    cy.byTestId('keptn-sli-breakdown-row-go_routines').find('dt-cell').eq(1).should('contain.text', 'Absolute change:');
    cy.byTestId('keptn-sli-breakdown-row-go_routines').find('dt-cell').eq(1).should('contain.text', 'Relative change:');
    cy.byTestId('keptn-sli-breakdown-row-go_routines').find('dt-cell').eq(1).should('contain.text', 'Compared with:');
  });

  it('should sort elements correctly', () => {
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

    // sort name asc
    let nameColumnHeader = cy
      .byTestId('keptn-sli-breakdown')
      .find('dt-header-row')
      .first()
      .find('dt-header-cell')
      .eq(1);
    nameColumnHeader.click();
    nameColumnHeader.should('have.class', 'dt-header-cell');

    nameColumnHeader.find('.dt-sort-header-container').first().should('have.class', 'dt-sort-header-sorted');
    nameColumnHeader.find('dt-icon').first().invoke('attr', 'ng-reflect-name').should('equal', 'sorter2-up');

    cy.byTestId('keptn-sli-breakdown').find('dt-row').eq(0).find('dt-cell').eq(1).should('have.text', 'go_routines');
    cy.byTestId('keptn-sli-breakdown')
      .find('dt-row')
      .eq(1)
      .find('dt-cell')
      .eq(1)
      .should('have.text', 'http_response_time_seconds_main_page_sum');

    // sort name desc
    nameColumnHeader = cy.byTestId('keptn-sli-breakdown').find('dt-header-row').first().find('dt-header-cell').eq(1);
    nameColumnHeader.click();

    nameColumnHeader.find('.dt-sort-header-container').first().should('have.class', 'dt-sort-header-sorted');
    nameColumnHeader.find('dt-icon').first().invoke('attr', 'ng-reflect-name').should('equal', 'sorter2-down');

    cy.byTestId('keptn-sli-breakdown')
      .find('dt-row')
      .eq(0)
      .find('dt-cell')
      .eq(1)
      .should('have.text', 'request_throughput');
    cy.byTestId('keptn-sli-breakdown')
      .find('dt-row')
      .eq(1)
      .find('dt-cell')
      .eq(1)
      .should('have.text', 'http_response_time_seconds_main_page_sum');
  });
});
