describe('Trigger a sequence', () => {
  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata.mock' });
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
    cy.intercept('/api/project/sockshop?approval=true&remediation=true', { fixture: 'project.mock' });
    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { fixture: 'projects.mock' });
    cy.intercept('/api/hasUnreadUniformRegistrationLogs', { body: false });
    cy.intercept('/api/project/sockshop/services', { body: ['carts', 'carts-db'] });
    cy.intercept('/api/project/sockshop/stages', { body: ['dev', 'staging', 'production'] });
    cy.intercept('/api/project/sockshop/customSequences', { body: ['delivery-direct', 'rollback', 'remediation'] });
    cy.intercept('/api/project/sockshop/serviceStates', { body: [] });
    cy.intercept('POST', '/api/v1/event', { body: { keptnContext: '6c98fbb0-4c40-4bff-ba9f-b20556a57c8a' } });
    cy.intercept('POST', '/api/controlPlane/v1/project/sockshop/stage/dev/service/carts/evaluation', {
      body: { keptnContext: '6c98fbb0-4c40-4bff-ba9f-b20556a57c8a' },
    });
    cy.intercept('/api/mongodb-datastore/event?keptnContext=6c98fbb0-4c40-4bff-ba9f-b20556a57c8a', {
      body: {
        events: [
          {
            data: {
              deployment: {
                deploymentNames: null,
              },
              evaluation: {
                end: '2022-02-23T10:36:47.662Z',
                start: '2022-02-23T09:36:47.662Z',
                timeframe: '',
              },
              project: 'podtato-head',
              service: 'helloservice',
              stage: 'hardening',
              test: {
                end: '',
                start: '',
              },
            },
            id: 'd792079c-1627-48f1-b66e-1cd6b2002b3c',
            source: 'https://github.com/keptn/keptn/api',
            specversion: '1.0',
            time: '2022-02-23T09:36:46.688Z',
            type: 'sh.keptn.event.hardening.evaluation.triggered',
            shkeptncontext: '8e548d83-bb84-4a62-9ea9-4609d0882f97',
            shkeptnspecversion: '0.2.3',
          },
        ],
        pageSize: 20,
        totalCount: 1,
      },
    });

    // Sequence screen
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25', { fixture: 'sequences.sockshop' });
    cy.intercept('/api/project/sockshop/sequences/metadata', { fixture: 'sequence.metadata.mock' });
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25&fromTime=*', {
      body: {
        states: [],
      },
    });
    cy.intercept(
      '/api/controlPlane/v1/sequence/sockshop?pageSize=1&keptnContext=6c98fbb0-4c40-4bff-ba9f-b20556a57c8a',
      { body: [] }
    );

    cy.visit('project/sockshop');
  });

  it('should navigate through all forms and close it from everywhere properly', () => {
    // Opening of triggering component
    cy.byTestId('keptn-trigger-button-open').click();
    cy.byTestId('keptn-trigger-button-open').should('not.exist');
    cy.byTestId('keptn-trigger-entry-h2').should('have.text', 'Trigger a new sequence for project sockshop');
    cy.byTestId('keptn-trigger-button-next').should('be.disabled');

    // Closing of triggering component from entry
    cy.byTestId('keptn-trigger-button-close').click();
    cy.byTestId('keptn-trigger-entry-h2').should('not.exist');
    cy.byTestId('keptn-trigger-button-open').should('exist');

    // Delivery navigations
    cy.byTestId('keptn-trigger-button-open').click();
    selectDelivery();
    testNavigationFirstPart('keptn-trigger-delivery-h2', 'Trigger a delivery for carts in dev');
    selectDelivery();
    testNavigationSecondPart('keptn-trigger-delivery-h2');

    // Evaluation navigations
    cy.byTestId('keptn-trigger-button-open').click();
    selectEvaluation();
    testNavigationFirstPart('keptn-trigger-evaluation-h2', ' Trigger an evaluation for carts in dev ');
    selectEvaluation();
    testNavigationSecondPart('keptn-trigger-evaluation-h2');

    // Custom sequence navigations
    cy.byTestId('keptn-trigger-button-open').click();
    selectCustomSequence();
    testNavigationFirstPart('keptn-trigger-custom-h2', ' Trigger a custom sequence for carts in dev ');
    selectCustomSequence();
    testNavigationSecondPart('keptn-trigger-custom-h2');
  });

  it('should trigger a delivery sequence', () => {
    cy.byTestId('keptn-trigger-button-open').click();
    cy.byTestId('keptn-trigger-button-next').should('be.disabled');
    selectDelivery();
    cy.byTestId('keptn-trigger-button-next').should('be.enabled').click();
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');
    cy.byTestId('keptn-trigger-delivery-labels').type('key1=val1');
    cy.byTestId('keptn-trigger-delivery-values').type('{"key2');
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');
    cy.byTestId('keptn-trigger-delivery-values-error').should('exist');
    cy.byTestId('keptn-trigger-delivery-values').type('": "val2"}');
    cy.byTestId('keptn-trigger-delivery-values-error').should('not.exist');
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');
    cy.byTestId('keptn-trigger-delivery-image').type('docker.io/keptn');
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');
    cy.byTestId('keptn-trigger-delivery-tag').type('v0.1.2');
    cy.byTestId('keptn-trigger-button-trigger').should('be.enabled').click();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should trigger an evaluation sequence with a timeframe', () => {
    cy.byTestId('keptn-trigger-button-open').click();
    cy.byTestId('keptn-trigger-button-next').should('be.disabled');
    selectEvaluation();
    cy.byTestId('keptn-trigger-button-next').should('be.enabled').click();
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');
    cy.byTestId('keptn-trigger-evaluation-type').children().first().click();
    cy.byTestId('keptn-trigger-evaluation-labels').type('key1=val1');
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');
    cy.byTestId('keptn-trigger-button-starttime').click();
    selectDateTime(0);
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');
    cy.byTestId('keptn-time-input-hours').type('0');
    cy.byTestId('keptn-trigger-button-trigger').should('be.enabled');
    cy.byTestId('keptn-time-input-minutes').type('1');
    cy.byTestId('keptn-trigger-button-trigger').should('be.enabled');
    cy.byTestId('keptn-time-input-seconds').type('15');
    cy.byTestId('keptn-trigger-button-trigger').should('be.enabled');
    cy.byTestId('keptn-time-input-millis').type('0');
    cy.byTestId('keptn-trigger-button-trigger').should('be.enabled');
    cy.byTestId('keptn-time-input-micros').type('0');
    cy.byTestId('keptn-trigger-button-trigger').should('be.enabled').click();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should trigger an evaluation sequence with a start end date', () => {
    cy.byTestId('keptn-trigger-button-open').click();
    cy.byTestId('keptn-trigger-button-next').should('be.disabled');
    selectEvaluation();
    cy.byTestId('keptn-trigger-button-next').should('be.enabled').click();
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');
    cy.byTestId('keptn-trigger-evaluation-type').children().eq(1).click();
    cy.byTestId('keptn-trigger-evaluation-labels').type('key1=val1');
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');

    // End before start date error
    cy.byTestId('keptn-trigger-button-starttime').click();
    selectDateTime(1);
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');
    cy.byTestId('keptn-trigger-button-endtime').click();
    selectDateTime(0);
    cy.byTestId('keptn-trigger-evaluation-date-error').should('exist');
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');

    // Correct date order
    cy.byTestId('keptn-trigger-button-starttime').click();
    selectDateTime(0);
    cy.byTestId('keptn-trigger-button-endtime').click();
    selectDateTime(1);
    cy.byTestId('keptn-trigger-evaluation-date-error').should('not.exist');

    cy.byTestId('keptn-trigger-button-trigger').should('be.enabled').click();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should trigger a custom sequence', () => {
    cy.byTestId('keptn-trigger-button-open').click();
    cy.byTestId('keptn-trigger-button-next').should('be.disabled');
    selectCustomSequence();
    cy.byTestId('keptn-trigger-button-next').should('be.enabled').click();
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');
    cy.byTestId('keptn-trigger-custom-labels').type('key1=val1');
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');
    cy.byTestId('keptn-trigger-custom-sequence').click();
    cy.wait(500);
    cy.get('dt-option').eq(0).click();
    cy.byTestId('keptn-trigger-button-trigger').should('be.enabled').click();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should open the trigger form from the sequence screen', () => {
    cy.visit('/project/sockshop/sequence');
    cy.byTestId('keptn-trigger-button-open').should('exist').click();
    cy.url().should('include', '/project/sockshop');
    cy.byTestId('keptn-trigger-entry-h2')
      .should('exist')
      .should('have.text', 'Trigger a new sequence for project sockshop');
  });

  it('should have the selected stage preselected', () => {
    testStageSelection(0, 'dev');
    testStageSelection(1, 'staging');
    testStageSelection(2, 'production');
  });

  function selectDelivery(): void {
    cy.byTestId('keptn-trigger-sequence-selection').children().first().click();
    selectServiceAndStage();
  }

  function selectEvaluation(): void {
    cy.byTestId('keptn-trigger-sequence-selection').children().eq(1).click();
    selectServiceAndStage();
  }

  function selectCustomSequence(): void {
    cy.wait(500);
    cy.byTestId('keptn-trigger-sequence-selection').children().eq(2).click();
    selectServiceAndStage();
  }

  function selectServiceAndStage(): void {
    cy.byTestId('keptn-trigger-service-selection').click();
    cy.wait(500);
    cy.get('dt-option').eq(0).click();
    cy.byTestId('keptn-trigger-stage-selection').click();
    cy.wait(500);
    cy.get('dt-option').eq(0).click();
  }

  function testNavigationFirstPart(h2Selector: string, expectedText: string): void {
    cy.byTestId('keptn-trigger-button-next').should('be.enabled').click();
    cy.byTestId(h2Selector).should('have.text', expectedText);
    cy.byTestId('keptn-trigger-button-trigger').should('be.disabled');
    cy.byTestId('keptn-trigger-button-close').click();
    cy.byTestId(h2Selector).should('not.exist');
    cy.byTestId('keptn-trigger-button-open').should('exist');
    cy.byTestId('keptn-trigger-button-open').click();
  }

  function testNavigationSecondPart(h2Selector: string): void {
    cy.byTestId('keptn-trigger-button-next').click();
    cy.byTestId('keptn-trigger-button-back').click();
    cy.byTestId('keptn-trigger-entry-h2').should('exist');
    cy.byTestId(h2Selector).should('not.exist');
    cy.byTestId('keptn-trigger-button-close').click();
  }

  function selectDateTime(calElement: number): void {
    cy.get('.dt-calendar-table-cell').eq(calElement).click();
    cy.byTestId('keptn-datetime-picker-submit').should('be.enabled');
    cy.get('[uitestid="keptn-datetime-picker-time"] [uitestid="keptn-time-input-hours"]').type('1');
    cy.byTestId('keptn-datetime-picker-submit').should('be.disabled');
    cy.get('[uitestid="keptn-datetime-picker-time"] [uitestid="keptn-time-input-minutes"]').type('15');
    cy.byTestId('keptn-datetime-picker-submit').should('be.disabled');
    cy.get('[uitestid="keptn-datetime-picker-time"] [uitestid="keptn-time-input-seconds"]').type('0');
    cy.byTestId('keptn-datetime-picker-submit').should('be.enabled').click();
  }

  function testStageSelection(elem: number, expectedText: string): void {
    cy.get('.stage-list .ktb-selectable-tile').eq(elem).find('h2').click();
    cy.byTestId('keptn-trigger-button-open').click();
    cy.get('[uitestid="keptn-trigger-stage-selection"] .dt-select-value-text > span').should('have.text', expectedText);
    cy.byTestId('keptn-trigger-button-close').click();
  }
});
