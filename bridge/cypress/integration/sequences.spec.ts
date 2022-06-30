import { SequencesPage } from '../support/pageobjects/SequencesPage';
import { interceptProjectBoard } from '../support/intercept';
import EnvironmentPage from '../support/pageobjects/EnvironmentPage';

describe('Sequences', () => {
  const sequencePage = new SequencesPage();
  const environmentPage = new EnvironmentPage();

  beforeEach(() => {
    interceptProjectBoard();
    sequencePage.intercept();
  });

  it('should show a loading indicator when sequences are not loaded', () => {
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25', {
      delay: 10_000,
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
    sequencePage.visit('sockshop').assertSequenceCount(25);
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

  it('should load older sequences', () => {
    sequencePage.visit('sockshop');
    cy.wait('@Sequences');
    cy.wait(500);

    sequencePage
      .assertSequenceCount(25)
      .assertLoadOlderSequencesButtonExists(true)

      .clickLoadOlderSequences()
      .assertSequenceCount(35)
      .assertLoadOlderSequencesButtonExists(true)

      .clickLoadOlderSequences()
      .assertSequenceCount(40)
      .assertLoadOlderSequencesButtonExists(false);
  });

  it('should not show show-older-sequences button if initially all sequences are loaded', () => {
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25', {
      fixture: 'eventByContext.mock',
    }).as('Sequences');

    sequencePage.visit('sockshop');
    cy.wait('@Sequences');
    sequencePage.assertSequenceCount(1).assertLoadOlderSequencesButtonExists(false);
  });

  describe('filtering', () => {
    it('should show a filtered list if filters are applied for Service', () => {
      sequencePage.visit('sockshop');
      cy.wait('@Sequences');
      cy.wait(500);

      sequencePage
        .checkServiceFilter('carts')
        .assertSequenceCount(22)
        .assertServiceNameOfSequences('carts')

        .clearFilter()
        .checkServiceFilter('carts-db')
        .assertSequenceCount(3)
        .assertServiceNameOfSequences('carts-db');
    });

    it('should show a filtered list if filters are applied for Stage', () => {
      sequencePage.visit('sockshop');
      cy.wait('@Sequences');
      cy.wait(500);

      sequencePage
        .checkStageFilter('production')
        .assertSequenceCount(7)
        .assertStageNamesOfSequences(['production'], false)

        .checkStageFilter('staging')
        .assertSequenceCount(6)
        .assertStageNamesOfSequences(['staging', 'production'], false)

        .checkStageFilter('dev')
        .assertSequenceCount(6)
        .assertStageNamesOfSequences(['dev', 'staging', 'production'], false);
    });

    it('should show a filtered list if filters are applied for Sequence', () => {
      sequencePage.visit('sockshop');
      cy.wait('@Sequences');
      cy.wait(500);

      sequencePage
        .checkSequenceFilter('delivery')
        .assertSequenceCount(16)
        .assertSequenceNameOfSequences('delivery')

        .clearFilter()
        .checkSequenceFilter('delivery-direct')
        .assertSequenceCount(3)
        .assertSequenceNameOfSequences('delivery-direct');
    });

    it('should show a filtered list if filters are applied for Status', () => {
      sequencePage.visit('sockshop');
      cy.wait('@Sequences');
      cy.wait(500);

      sequencePage
        .checkStatusFilter('Active')
        .assertSequenceCount(2)
        .assertStatusOfSequences('started')

        .clearFilter()
        .checkStatusFilter('Failed')
        .assertSequenceCount(14)
        .assertStatusOfSequences('failed')

        .clearFilter()
        .checkStatusFilter('Aborted')
        .assertNoSequencesFilteredMessageExists(true)

        .clearFilter()
        .checkStatusFilter('Succeeded')
        .assertSequenceCount(5)
        .assertStatusOfSequences('succeeded')
        .assertLoadingOldSequencesButtonExists(false);
    });

    it('should show a filtered list if combined filters are applied', () => {
      sequencePage.visit('sockshop');
      cy.wait('@Sequences');
      cy.wait(500);

      sequencePage
        .checkServiceFilter('carts')
        .checkStageFilter('production')
        .checkSequenceFilter('delivery')
        .checkStatusFilter('Succeeded')

        .assertSequenceCount(2)
        .assertStageNameOfSequences('production')
        .assertServiceNameOfSequences('carts')
        .assertSequenceNameOfSequences('delivery')
        .assertStatusOfSequences('succeeded');
    });

    it('should filter waiting sequences', () => {
      sequencePage.visit('sockshop');
      cy.wait('@Sequences');
      cy.wait(500);

      sequencePage
        .checkStatusFilter('Waiting')

        .assertSequenceCount(1)
        .assertStatusOfSequences('waiting');
    });

    it('should save filters to query params', () => {
      sequencePage.visit('sockshop');
      cy.wait('@Sequences');
      cy.wait(500);

      sequencePage
        .checkServiceFilter('carts')
        .checkStageFilter('dev')
        .checkStageFilter('production')
        .checkSequenceFilter('delivery')
        .checkStatusFilter('Active')

        .assertQueryParams('?Service=carts&Stage=dev&Stage=production&Sequence=delivery&Status=started');
    });

    it('should load filters from query params', () => {
      sequencePage.visit('sockshop', {
        Stage: 'dev',
        Service: 'carts',
        Sequence: 'delivery',
        Status: 'started',
      });
      cy.wait('@Sequences');
      cy.wait(500);

      sequencePage
        .assertFilterIsChecked('Stage', 'dev', true)
        .assertFilterIsChecked('Stage', 'staging', false)
        .assertFilterIsChecked('Stage', 'production', false)
        .assertFilterIsChecked('Service', 'carts', true)
        .assertFilterIsChecked('Service', 'carts-db', false)
        .assertFilterIsChecked('Sequence', 'delivery', true)
        .assertFilterIsChecked('Sequence', 'delivery-direct', false)
        .assertFilterIsChecked('Status', 'Active', true)
        .assertFilterIsChecked('Status', 'Waiting', false)
        .assertFilterIsChecked('Status', 'Failed', false)
        .assertFilterIsChecked('Status', 'Aborted', false)
        .assertFilterIsChecked('Status', 'Succeeded', false)
        .assertSequenceCount(2)
        .assertStatusOfSequences('started');
    });

    it('should load filters from local storage', () => {
      sequencePage.visit('sockshop', {
        Stage: 'staging',
        Service: 'carts',
        Sequence: 'delivery',
        Status: 'started',
      });
      environmentPage.intercept().visit('sockshop');

      sequencePage.visit('sockshop', {});
      cy.wait('@Sequences');
      cy.wait(500);

      sequencePage
        .assertFilterIsChecked('Stage', 'dev', false)
        .assertFilterIsChecked('Stage', 'staging', true)
        .assertFilterIsChecked('Stage', 'production', false)
        .assertFilterIsChecked('Service', 'carts', true)
        .assertFilterIsChecked('Service', 'carts-db', false)
        .assertFilterIsChecked('Sequence', 'delivery', true)
        .assertFilterIsChecked('Sequence', 'delivery-direct', false)
        .assertFilterIsChecked('Status', 'Active', true)
        .assertFilterIsChecked('Status', 'Waiting', false)
        .assertFilterIsChecked('Status', 'Failed', false)
        .assertFilterIsChecked('Status', 'Aborted', false)
        .assertFilterIsChecked('Status', 'Succeeded', false)
        .assertSequenceCount(2)
        .assertStatusOfSequences('started');
    });

    it('should apply filters also when loading more sequences', () => {
      sequencePage.visit('sockshop', {
        Stage: 'staging',
        Sequence: 'evaluation',
        Service: 'carts',
      });
      cy.wait('@Sequences');
      cy.wait(500);

      sequencePage
        .assertSequenceCount(3)
        .assertLoadOlderSequencesButtonExists(true)

        .clickLoadOlderSequences()
        .assertSequenceCount(12)
        .assertLoadOlderSequencesButtonExists(true)

        .clickLoadOlderSequences()
        .assertSequenceCount(14)
        .assertLoadOlderSequencesButtonExists(false);
    });

    it('should load older sequences also if filter initially has empty list', () => {
      sequencePage
        .visit('sockshop', {
          Stage: 'testing',
          Sequence: 'evaluation',
          Service: 'carts',
        })

        .assertSequenceCount(0)
        .assertLoadOlderSequencesButtonExists(true)

        .clickLoadOlderSequences()
        .assertSequenceCount(0)
        .assertLoadOlderSequencesButtonExists(true)

        .clickLoadOlderSequences()
        .assertSequenceCount(1)
        .assertLoadOlderSequencesButtonExists(false);
    });

    it('should not add the filter to the URL after it was cleared and the page was reloaded', () => {
      sequencePage.visit('sockshop', {
        Stage: 'staging',
        Sequence: 'evaluation',
        Service: 'carts',
      });
      cy.wait('@Sequences');
      cy.wait(500);

      sequencePage
        .assertAmountOfQueryParameters(3)
        .reload()
        .clearFilter()
        .selectSequence('62cca6f3-dc54-4df6-a04c-6ffc894a4b5e')
        .assertAmountOfQueryParameters(0);
    });
  });
});
