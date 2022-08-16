import ServicesPage from '../support/pageobjects/ServicesPage';
import { HeatmapComponentPage } from '../support/pageobjects/HeatmapComponentPage';
import { interceptD3 } from '../support/intercept';
import { ResultTypes } from '../../shared/models/result-types';

describe('sli-breakdown', () => {
  const servicesPage = new ServicesPage();

  beforeEach(() => {
    servicesPage.interceptAll();
    interceptD3();
    servicesPage.visitServicePage('sockshop').selectService('carts', 'v0.1.2').waitForEvaluations();
  });

  it('should load the heatmap with sli breakdown in service screen', () => {
    servicesPage
      .verifySliBreakdown(
        {
          name: 'go_routines',
          value: 88,
          result: ResultTypes.PASSED,
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
          availableScore: 33.33,
        },
        false
      )
      .verifySliBreakdown(
        {
          name: 'request_throughput',
          value: 18.42,
          result: ResultTypes.FAILED,
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
          availableScore: 33.33,
        },
        false
      );
  });

  it('should show more details when expanding sli breakdown in service screen', () => {
    servicesPage.expandSliBreakdown('go_routines').verifySliBreakdown(
      {
        name: 'go_routines',
        value: 88,
        result: ResultTypes.PASSED,
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
        availableScore: 33.33,
      },
      true
    );
  });

  it('should sort elements correctly', () => {
    // sort name asc
    servicesPage
      .clickSliBreakdownHeader('Name')
      .verifySliBreakdownSorting(1, 'ascending', 'go_routines', 'http_response_time_seconds_main_page_sum');

    // sort name desc
    servicesPage
      .clickSliBreakdownHeader('Name')
      .verifySliBreakdownSorting(1, 'descending', 'request_throughput', 'http_response_time_seconds_main_page_sum');

    // sort score asc
    servicesPage.clickSliBreakdownHeader('Score').verifySliBreakdownSorting(7, 'ascending', ' 0/33.33 ', ' 0/33.33 ');

    // sort score desc
    servicesPage
      .clickSliBreakdownHeader('Score')
      .verifySliBreakdownSorting(7, 'descending', ' 33.99/33.33 ', ' 0/33.33 ');
  });

  describe('score overlay', () => {
    it('should show default score overlay', () => {
      servicesPage.assertSliScoreOverlayDefault('go_routines');
    });

    it('should show score overlay with warning message', () => {
      const heatmapPage = new HeatmapComponentPage();
      heatmapPage.selectEvaluation('8a549059-8dcd-43ea-adff-b7c2ea9a0d99');
      servicesPage.assertSliScoreOverlayWarning('http_response_time_seconds_main_page_sum');
    });

    it('should show score overlay with error message', () => {
      servicesPage.assertSliScoreOverlayFailed('http_response_time_seconds_main_page_sum');
    });
  });
});

describe('sli-breakdown with fallback api call', () => {
  const servicesPage = new ServicesPage();
  const heatmapPage = new HeatmapComponentPage();

  beforeEach(() => {
    servicesPage.interceptAll().interceptSliFallback('sockshop', ['91a77341-fe5e-43e1-a8a7-be9761b9cee5']);
    interceptD3();
  });

  it('should fallback to api call', () => {
    servicesPage.visitServicePage('sockshop').selectService('carts', 'v0.1.2');
    heatmapPage.selectEvaluation('c1b2761f-5b6d-4bdc-9bb7-4661a05ea3b2');
    servicesPage
      .waitForSliFallbackFetch()
      .expandSliBreakdown('go_routines')
      .expandSliBreakdown('http_response_time_seconds_main_page_sum')
      .expandSliBreakdown('request_throughput')
      .assertSliValueColumnExpanded('go_routines', 10, 2, 25, 8)
      .assertSliValueColumnExpanded('http_response_time_seconds_main_page_sum', 0, 0, 0, 0)
      .assertSliValueColumnExpanded('request_throughput', 0, 0, 0, 0)
      .assertSliBreakdownLoading(false);
  });

  it('should show loading indicator if fallback api call is triggered', () => {
    servicesPage
      .interceptSliFallback('sockshop', ['91a77341-fe5e-43e1-a8a7-be9761b9cee5'], true)
      .visitServicePage('sockshop')
      .selectService('carts', 'v0.1.2');
    heatmapPage.selectEvaluation('c1b2761f-5b6d-4bdc-9bb7-4661a05ea3b2');
    servicesPage.assertSliBreakdownLoading(true);
  });
});
