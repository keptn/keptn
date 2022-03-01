import EnvironmentPage from '../support/pageobjects/EnvironmentPage';

const environmentPage = new EnvironmentPage();

describe('Trigger a sequence', () => {
  beforeEach(() => {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata.mock' });
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfoCD.mock' });
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
    cy.intercept(
      '/api/controlPlane/v1/sequence/sockshop?pageSize=1&keptnContext=6c98fbb0-4c40-4bff-ba9f-b20556a57c8a',
      {
        body: {
          states: [
            {
              name: 'delivery',
              service: 'carts',
              project: 'sockshop',
              time: '2022-02-23T14:28:50.504Z',
              shkeptncontext: '6c98fbb0-4c40-4bff-ba9f-b20556a57c8a',
              state: 'finished',
              stages: [
                {
                  name: 'dev',
                  state: 'finished',
                  latestEvent: {
                    type: 'sh.keptn.event.dev.delivery.finished',
                    id: '1341268c-c899-4314-b87c-9f4ea6566208',
                    time: '2022-02-23T14:28:51.596Z',
                  },
                  latestFailedEvent: {
                    type: 'sh.keptn.event.dev.delivery.finished',
                    id: '1341268c-c899-4314-b87c-9f4ea6566208',
                    time: '2022-02-23T14:28:51.596Z',
                  },
                },
              ],
            },
          ],
          totalCount: 1,
        },
      }
    );

    // Sequence screen
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25', { fixture: 'sequences.sockshop' });
    cy.intercept('/api/project/sockshop/sequences/metadata', { fixture: 'sequence.metadata.mock' });
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25&fromTime=*', {
      body: {
        states: [],
      },
    });

    cy.visit('project/sockshop');
  });

  it('should navigate through all forms and close it from everywhere properly', () => {
    // Opening of triggering component
    environmentPage.openTriggerSequence();
    environmentPage.assertOpenTriggerSequenceExists(false);
    environmentPage.assertTriggerEntryH2('have.text', 'Trigger a new sequence for project sockshop');
    environmentPage.assertTriggerNextPageEnabled(false);

    // Closing of triggering component from entry
    environmentPage.clickTriggerClose();
    environmentPage.assertTriggerEntryH2('not.exist');
    environmentPage.assertOpenTriggerSequenceExists(true);

    // Delivery navigations
    environmentPage.openTriggerSequence();
    environmentPage.selectTriggerDelivery();
    testNavigationFirstPart('keptn-trigger-delivery-h2', 'Trigger a delivery for carts in dev');
    environmentPage.selectTriggerDelivery();
    testNavigationSecondPart('keptn-trigger-delivery-h2');

    // Evaluation navigations
    environmentPage.openTriggerSequence();
    environmentPage.selectTriggerEvaluation();
    testNavigationFirstPart('keptn-trigger-evaluation-h2', ' Trigger an evaluation for carts in dev ');
    environmentPage.selectTriggerEvaluation();
    testNavigationSecondPart('keptn-trigger-evaluation-h2');

    // Custom sequence navigations
    environmentPage.openTriggerSequence();
    environmentPage.selectTriggerCustomSequence();
    testNavigationFirstPart('keptn-trigger-custom-h2', ' Trigger a custom sequence for carts in dev ');
    environmentPage.selectTriggerCustomSequence();
    testNavigationSecondPart('keptn-trigger-custom-h2');
  });

  it('should trigger a delivery sequence', () => {
    environmentPage.openTriggerSequence();
    environmentPage.assertTriggerNextPageEnabled(false);
    environmentPage.selectTriggerDelivery();
    environmentPage.assertTriggerNextPageEnabled(true).click();
    environmentPage.assertTriggerSequenceEnabled(false);
    environmentPage.typeTriggerDeliveryLabels('key1=val1');
    environmentPage.typeTriggerDeliveryValues('{"key2');
    environmentPage.assertTriggerSequenceEnabled(false);
    environmentPage.assertTriggerDeliveryValuesErrorExists(true);
    environmentPage.typeTriggerDeliveryValues('": "val2"}');
    environmentPage.assertTriggerDeliveryValuesErrorExists(false);
    environmentPage.assertTriggerSequenceEnabled(false);
    environmentPage.typeTriggerDeliveryImage('docker.io/keptn');
    environmentPage.assertTriggerSequenceEnabled(false);
    environmentPage.typeTriggerDeliveryTag('v0.1.2');
    environmentPage.assertTriggerSequenceEnabled(true).click();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should trigger an evaluation sequence with a timeframe', () => {
    environmentPage.openTriggerSequence();
    environmentPage.assertTriggerNextPageEnabled(false);
    environmentPage.selectTriggerEvaluation();
    environmentPage.assertTriggerNextPageEnabled(true).click();
    environmentPage.assertTriggerSequenceEnabled(false);
    environmentPage.selectTriggerEvaluationType(0);
    environmentPage.typeTriggerEvaluationLabels('key1=val1');
    environmentPage.assertTriggerSequenceEnabled(false);
    environmentPage.clickTriggerStartTime();
    environmentPage.selectTriggerDateTime(0);
    environmentPage.assertTriggerSequenceEnabled(false);
    environmentPage.typeTriggerEvaluationTimeInput('hours', '0');
    environmentPage.assertTriggerSequenceEnabled(true);
    environmentPage.typeTriggerEvaluationTimeInput('minutes', '1');
    environmentPage.assertTriggerSequenceEnabled(true);
    environmentPage.typeTriggerEvaluationTimeInput('seconds', '15');
    environmentPage.assertTriggerSequenceEnabled(true);
    environmentPage.typeTriggerEvaluationTimeInput('millis', '0');
    environmentPage.assertTriggerSequenceEnabled(true);
    environmentPage.typeTriggerEvaluationTimeInput('micros', '0');
    environmentPage.assertTriggerSequenceEnabled(true).click();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should trigger an evaluation sequence with a start end date', () => {
    environmentPage.openTriggerSequence();
    environmentPage.assertTriggerNextPageEnabled(false);
    environmentPage.selectTriggerEvaluation();
    environmentPage.assertTriggerNextPageEnabled(true).click();
    environmentPage.assertTriggerSequenceEnabled(false);
    environmentPage.selectTriggerEvaluationType(1);
    environmentPage.typeTriggerEvaluationLabels('key1=val1');
    environmentPage.assertTriggerSequenceEnabled(false);

    // End before start date error
    environmentPage.clickTriggerStartTime();
    environmentPage.selectTriggerDateTime(1);
    environmentPage.assertTriggerSequenceEnabled(false);
    environmentPage.clickTriggerEndTime();
    environmentPage.selectTriggerDateTime(0);
    environmentPage.assertTriggerEvaluationDateErrorExists(true);
    environmentPage.assertTriggerSequenceEnabled(false);

    // Correct date order
    environmentPage.clickTriggerStartTime();
    environmentPage.selectTriggerDateTime(0);
    environmentPage.clickTriggerEndTime();
    environmentPage.selectTriggerDateTime(1);
    environmentPage.assertTriggerEvaluationDateErrorExists(false);

    environmentPage.assertTriggerSequenceEnabled(true).click();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should trigger a custom sequence', () => {
    environmentPage.openTriggerSequence();
    environmentPage.assertTriggerNextPageEnabled(false);
    environmentPage.selectTriggerCustomSequence();
    environmentPage.assertTriggerNextPageEnabled(true).click();
    environmentPage.assertTriggerSequenceEnabled(false);
    environmentPage.typeTriggerCustomLabels('key1=val1');
    environmentPage.assertTriggerSequenceEnabled(false);
    environmentPage.selectTriggerCustomSequenceType(0);
    environmentPage.assertTriggerSequenceEnabled(true).click();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should open the trigger form from the sequence screen', () => {
    cy.visit('/project/sockshop/sequence');
    environmentPage.assertOpenTriggerSequenceExists(true).click();
    cy.url().should('include', '/project/sockshop');
    environmentPage.assertTriggerEntryH2('exist');
    environmentPage.assertTriggerEntryH2('have.text', 'Trigger a new sequence for project sockshop');
  });

  it('should have the selected stage preselected', () => {
    environmentPage.assertTriggerStageSelection(0, 'dev');
    environmentPage.assertTriggerStageSelection(1, 'staging');
    environmentPage.assertTriggerStageSelection(2, 'production');
  });

  function testNavigationFirstPart(h2Selector: string, expectedText: string): void {
    environmentPage.assertTriggerNextPageEnabled(true).click();
    cy.byTestId(h2Selector).should('have.text', expectedText);
    environmentPage.assertTriggerSequenceEnabled(false);
    environmentPage.clickTriggerClose();
    cy.byTestId(h2Selector).should('not.exist');
    environmentPage.assertOpenTriggerSequenceExists(true).click();
  }

  function testNavigationSecondPart(h2Selector: string): void {
    environmentPage.clickTriggerNext();
    environmentPage.clickTriggerBack();
    environmentPage.assertTriggerEntryH2('exist');
    cy.byTestId(h2Selector).should('not.exist');
    environmentPage.clickTriggerClose();
  }
});
