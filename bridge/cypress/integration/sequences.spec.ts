describe('Sequences', () => {
  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata.mock' });
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
    cy.intercept('/api/project/sockshop?approval=true&remediation=true', { fixture: 'project.mock' });
    cy.intercept('/api/hasUnreadUniformRegistrationLogs', { body: false });

    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { fixture: 'projects.mock' });
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25', { fixture: 'sequences.sockshop' }).as(
      'Sequences'
    );
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25&fromTime=*', {
      body: {
        states: [],
      },
    });

    cy.intercept('/api/project/sockshop/sequences/metadata', { fixture: 'sequence.metadata.mock' }).as(
      'SequencesMetadata'
    );
  });

  it('should show loading indicator when sequences are not loaded', () => {
    let sendResponse;
    const trigger = new Promise((resolve) => {
      sendResponse = resolve;
    });

    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25', (request) => {
      // eslint-disable-next-line promise/always-return
      return trigger.then(() => {
        request.reply({ fixture: 'sequences.sockshop' });
      });
    });

    cy.visit('/project/sockshop/sequence');

    cy.byTestId('keptn-loadingSequences').should('exist');

    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    sendResponse();
    cy.byTestId('keptn-loadingSequences').should('not.exist');
  });

  it('should show an empty state if no sequences are loaded', () => {
    cy.visit('/project/sockshop/sequence');
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25', {
      body: {
        states: [],
      },
    });

    cy.byTestId('keptn-noSequences').should('exist');
  });

  it('should a list of sequences if everything is loaded', () => {
    cy.visit('/project/sockshop/sequence');

    cy.byTestId('keptn-sequence-view-roots').get('ktb-selectable-tile').should('have.length', 4);
  });

  it('should show a filtered list if filters are applied', () => {
    cy.visit('/project/sockshop/sequence');
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
