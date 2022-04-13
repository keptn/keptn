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
    }).as('Sequences');

    sequencePage.visit('sockshop').assertIsLoadingSequences(true);
  });

  it('should show an empty state if no sequences are loaded', () => {
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25', {
      body: {
        states: [],
      },
    }).as('Sequences');
    sequencePage.visit('sockshop');

    cy.byTestId('keptn-noSequences').should('exist');
  });

  it('should show a list of sequences if everything is loaded', () => {
    sequencePage.visit('sockshop').assertSequenceCount(5);
  });

  it('should show a filtered list if filters are applied', () => {
    sequencePage.visit('sockshop');
    cy.wait('@SequencesMetadata');
    cy.wait('@Sequences');
    cy.wait(500);

    // Test single filters
    sequencePage
      .checkServiceFilter('carts')
      .assertSequenceCount(3)
      .assertServiceNameOfSequences('carts')

      .clearFilter()
      .checkServiceFilter('carts-db')
      .assertSequenceCount(2)
      .assertServiceNameOfSequences('carts-db')

      .clearFilter()
      .checkStageFilter('production')
      .assertSequenceCount(2)
      .assertStageNamesOfSequences(['dev', 'staging', 'production'])

      .clearFilter()
      .checkSequenceFilter('delivery')
      .assertSequenceCount(3)
      .assertSequenceNameOfSequences('delivery')

      .clearFilter()
      .checkSequenceFilter('delivery-direct')
      .assertSequenceCount(2)
      .assertSequenceNameOfSequences('delivery-direct')

      .clearFilter()
      .checkStatusFilter('Active')
      .assertNoSequencesMessageExists(true)

      .clearFilter()
      .checkStatusFilter('Failed')
      .assertSequenceCount(2)
      .assertStatusOfSequences('failed')

      .clearFilter()
      .checkStatusFilter('Aborted')
      .assertNoSequencesFilteredMessageExists(true)

      .clearFilter()
      .checkStatusFilter('Succeeded')
      .assertSequenceCount(2)
      .assertStatusOfSequences('succeeded')
      .assertLoadingOldSequencesButtonExists(false)

      // Test one combined filter
      .clearFilter()
      .checkServiceFilter('carts')
      .checkStageFilter('production')
      .checkSequenceFilter('delivery')
      .checkStatusFilter('Succeeded')
      .assertStageNameOfSequences('production')
      .assertServiceNameOfSequences('carts')
      .assertSequenceNameOfSequences('delivery')
      .assertStatusOfSequences('succeeded');
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

  it('should show waiting sequence', () => {
    const context = 'f78c2fc7-d272-4bcd-9845-3f3041080ae1';
    const project = 'sockshop';
    cy.intercept(`/api/mongodb-datastore/event?keptnContext=${context}&project=${project}`, {
      body: {
        events: [],
      },
    });

    sequencePage
      .visit(project)
      .assertIsWaitingSequence(context, true)
      .selectSequence(context)
      .assertIsSelectedSequenceWaiting(true);
  });

  it('should filter waiting sequences', () => {
    sequencePage.visit('sockshop');
    cy.wait('@SequencesMetadata');
    cy.wait('@Sequences');
    cy.wait(500);

    sequencePage.checkStatusFilter('Waiting').assertSequenceCount(1).assertStatusOfSequences('waiting');
  });
});
