import { SequencesPage } from '../support/pageobjects/SequencesPage';
import { interceptProjectBoard } from '../support/intercept';

describe('Sequences', () => {
  const sequencePage = new SequencesPage();

  beforeEach(() => {
    interceptProjectBoard();
    sequencePage.intercept();
  });

  it('should show a loading indicator when sequences are not loaded', () => {
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25', {
      delay: 2000,
      fixture: 'sequences.sockshop',
    });

    sequencePage.visit('sockshop');

    cy.byTestId('keptn-loadingSequences').should('exist');
  });

  it('should show an empty state if no sequences are loaded', () => {
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25', {
      body: {
        states: [],
      },
    });
    sequencePage.visit('sockshop');

    cy.byTestId('keptn-noSequences').should('exist');
  });

  it('should show a list of sequences if everything is loaded', () => {
    sequencePage.visit('sockshop');

    cy.byTestId('keptn-sequence-view-roots').get('ktb-selectable-tile').should('have.length', 4);
  });

  it('should show a filtered list if filters are applied', () => {
    sequencePage.visit('sockshop');
    cy.wait('@SequencesMetadata');
    cy.wait('@Sequences');
    cy.wait(500);

    // Test single filters

    cy.get('dt-quick-filter-group').eq(0).find('dt-checkbox').eq(0).click();
    testSelectableTiles(3, 'keptn-sequence-info-serviceName', 'carts');

    clearFilter();
    cy.get('dt-quick-filter-group').eq(0).find('dt-checkbox').eq(1).click();
    testSelectableTiles(1, 'keptn-sequence-info-serviceName', 'carts-db');

    clearFilter();
    cy.get('dt-quick-filter-group').eq(1).find('dt-checkbox').eq(2).click();
    cy.byTestId('keptn-sequence-info-stageDetails').each((el) => {
      cy.wrap(el).find('ktb-stage-badge').should('have.length', 3);
    });

    clearFilter();
    cy.get('dt-quick-filter-group').eq(2).find('dt-checkbox').eq(0).click();
    testSelectableTiles(3, 'keptn-sequence-info-sequenceName', 'delivery');

    clearFilter();
    cy.get('dt-quick-filter-group').eq(2).find('dt-checkbox').eq(1).click();
    testSelectableTiles(1, 'keptn-sequence-info-sequenceName', 'delivery-direct');

    clearFilter();
    cy.get('dt-quick-filter-group').eq(3).find('dt-checkbox').eq(0).click();
    cy.byTestId('keptn-noSequencesFiltered').should('exist');

    clearFilter();
    cy.get('dt-quick-filter-group').eq(3).find('dt-checkbox').eq(1).click();
    testSelectableTiles(2, 'keptn-sequence-info-status', 'failed');

    clearFilter();
    cy.get('dt-quick-filter-group').eq(3).find('dt-checkbox').eq(2).click();
    cy.byTestId('keptn-noSequencesFiltered').should('exist');

    clearFilter();
    cy.get('dt-quick-filter-group').eq(3).find('dt-checkbox').eq(3).click();
    testSelectableTiles(2, 'keptn-sequence-info-status', 'succeeded');
    cy.byTestId('keptn-loadingOldSequences').should('not.exist');

    // Test one combined filter
    clearFilter();
    cy.get('dt-quick-filter-group').eq(0).find('dt-checkbox').eq(0).click();
    cy.get('dt-quick-filter-group').eq(1).find('dt-checkbox').eq(2).click();
    cy.get('dt-quick-filter-group').eq(2).find('dt-checkbox').eq(0).click();
    cy.get('dt-quick-filter-group').eq(3).find('dt-checkbox').eq(3).click();

    testSelectableTiles(1, 'keptn-sequence-info-serviceName', 'carts');
    testSelectableTiles(1, 'keptn-sequence-info-sequenceName', 'delivery');
    testSelectableTiles(1, 'keptn-sequence-info-status', 'succeeded');
  });

  it('should select sequence and show the right timestamps in the timeline', () => {
    sequencePage
      .visit('sockshop')
      .selectSequence('62cca6f3-dc54-4df6-a04c-6ffc894a4b5e')
      .assertTimelineTime('dev', '12:41')
      .assertTimelineTime('staging', '12:42')
      .assertTimelineTime('production', '12:43');
  });

  it('should select sequence and show loading indicators if traces are not loaded yet', () => {
    cy.intercept('/api/mongodb-datastore/event?keptnContext=62cca6f3-dc54-4df6-a04c-6ffc894a4b5e&project=sockshop', {
      fixture: 'sequence.traces.mock.json',
      delay: 20_000,
    });
    sequencePage
      .visit('sockshop')
      .selectSequence('62cca6f3-dc54-4df6-a04c-6ffc894a4b5e')
      .assertTimelineTimeLoading('dev', true)
      .assertTimelineTimeLoading('staging', true)
      .assertTimelineTimeLoading('production', true);
  });

  it('should select sequence and should not have loading indicators', () => {
    sequencePage
      .visit('sockshop')
      .selectSequence('62cca6f3-dc54-4df6-a04c-6ffc894a4b5e')
      .assertTimelineTimeLoading('dev', false)
      .assertTimelineTimeLoading('staging', false)
      .assertTimelineTimeLoading('production', false);
  });

  it('should select sequence and should have image tag in name', () => {
    sequencePage
      .visit('sockshop')
      .selectSequence('62cca6f3-dc54-4df6-a04c-6ffc894a4b5e')
      .assertServiceName('carts', 'v0.13.1');
  });

  it('should select sequence and should not have image tag in name', () => {
    const context = 'bb03865b-2bdd-43cc-9848-2a9cced86ff3';
    const project = 'sockshop';
    cy.intercept(`/api/mongodb-datastore/event?keptnContext=${context}&project=${project}`, {
      fixture: 'sequence.traces.mock.json',
    });

    sequencePage.visit(project).selectSequence(context).assertServiceName('carts');
  });

  function clearFilter(): void {
    cy.get('.dt-filter-field-clear-all-button').click();
    cy.get('.dt-filter-field-input ').type('{esc}');
  }

  function testSelectableTiles(expectedLength: number, textElementSelector: string, expectedText: string): void {
    const selectableTiles = cy.byTestId('keptn-sequence-view-roots').get('ktb-selectable-tile');
    selectableTiles.should('have.length', expectedLength);

    cy.byTestId(textElementSelector).each((el) => {
      cy.wrap(el).should('have.text', expectedText);
    });
  }
});
