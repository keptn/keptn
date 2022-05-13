import ServicesPage from '../support/pageobjects/ServicesPage';
import { SliResult } from '../../client/app/_models/sli-result';
import { interceptProjectBoard } from '../support/intercept';

describe('sli-breakdown', () => {
  const servicesPage = new ServicesPage();

  beforeEach(() => {
    interceptProjectBoard();
  });

  it('should load the heatmap with sli breakdown in service screen', () => {
    servicesPage
      .intercept()
      .visitServicePage('sockshop')
      .selectService('carts', 'v0.1.2')
      .verifySliBreakdown(
        {
          name: 'go_routines',
          value: 88,
          result: 'pass',
          score: 33.99,
          passTargets: [
            {
              criteria: '<=100',
              targetValue: 100,
              violated: false,
            },
          ],
          warningTargets: null,
          keySli: false,
          success: true,
          expanded: false,
          weight: 1,
          comparedValue: 8,
          calculatedChanges: {
            absolute: 80,
            relative: 1000,
          },
        } as SliResult,
        false
      )
      .verifySliBreakdown(
        {
          name: 'request_throughput',
          value: 18.42,
          result: 'fail',
          score: 0,
          passTargets: [
            {
              criteria: '<=+100%',
              targetValue: 0,
              violated: true,
            },
            {
              criteria: '>=-80%',
              targetValue: 0,
              violated: false,
            },
          ],
          warningTargets: null,
          keySli: false,
          success: true,
          expanded: false,
          weight: 1,
          comparedValue: 0,
          calculatedChanges: {
            absolute: 18.42,
            relative: 1742,
          },
        } as SliResult,
        false
      );
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

    servicesPage
      .visitServicePage('sockshop')
      .selectService('carts', 'v0.1.2')
      .expandSliBreakdown('go_routines')
      .verifySliBreakdown(
        {
          name: 'go_routines',
          value: 88,
          result: 'pass',
          score: 33.99,
          passTargets: [
            {
              criteria: '<=100',
              targetValue: 100,
              violated: false,
            },
          ],
          warningTargets: null,
          keySli: false,
          success: true,
          expanded: false,
          weight: 1,
          comparedValue: 8,
          calculatedChanges: {
            absolute: 80,
            relative: 1000,
          },
        } as SliResult,
        true
      );
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

    servicesPage.visitServicePage('sockshop').selectService('carts', 'v0.1.2');

    // sort name asc
    servicesPage
      .clickSliBreakdownHeader('Name')
      .verifySliBreakdownSorting(1, 'ascending', 'go_routines', 'http_response_time_seconds_main_page_sum');

    // sort name desc
    servicesPage
      .clickSliBreakdownHeader('Name')
      .verifySliBreakdownSorting(1, 'descending', 'request_throughput', 'http_response_time_seconds_main_page_sum');

    // sort score asc
    servicesPage.clickSliBreakdownHeader('Score').verifySliBreakdownSorting(7, 'ascending', '0', '0');

    // sort score desc
    servicesPage.clickSliBreakdownHeader('Score').verifySliBreakdownSorting(7, 'descending', '33.99', '0');
  });
});
