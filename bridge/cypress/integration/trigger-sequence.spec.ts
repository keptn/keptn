import EnvironmentPage from '../support/pageobjects/EnvironmentPage';
import { interceptEnvironmentScreen } from '../support/intercept';

const environmentPage = new EnvironmentPage();

describe('Trigger a sequence', () => {
  beforeEach(() => {
    interceptEnvironmentScreen();
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfoCD.mock' });

    // Sequence screen
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25', { fixture: 'sequences.sockshop' });
    cy.intercept('/api/project/sockshop/sequences/metadata', { fixture: 'sequence.metadata.mock' });
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25&fromTime=*', {
      body: {
        states: [],
      },
    });

    environmentPage.visit('sockshop');
  });

  it('should navigate through all forms and close it from everywhere properly', () => {
    // Opening of triggering component
    environmentPage
      .clickTriggerOpen()
      .assertOpenTriggerSequenceExists(false)
      .assertTriggerEntryH2HasText('Trigger a new sequence for project sockshop')
      .assertTriggerNextPageEnabled(false);

    // Closing of triggering component from entry
    environmentPage.clickTriggerClose().assertTriggerEntryH2Exists(false).assertOpenTriggerSequenceExists(true);

    // Delivery navigations
    environmentPage.clickTriggerOpen().selectTriggerDelivery();
    testNavigationFirstPart('keptn-trigger-delivery-h2', 'Trigger a delivery for carts in dev', false);
    environmentPage.selectTriggerDelivery();
    testNavigationSecondPart('keptn-trigger-delivery-h2');

    // Evaluation navigations
    environmentPage.clickTriggerOpen().selectTriggerEvaluation();
    testNavigationFirstPart('keptn-trigger-evaluation-h2', ' Trigger an evaluation for carts in dev ', true);
    environmentPage.selectTriggerEvaluation();
    testNavigationSecondPart('keptn-trigger-evaluation-h2');

    // Custom sequence navigations
    environmentPage.clickTriggerOpen().selectTriggerCustomSequence();
    testNavigationFirstPart('keptn-trigger-custom-h2', ' Trigger a custom sequence for carts in dev ', false);
    environmentPage.selectTriggerCustomSequence();
    testNavigationSecondPart('keptn-trigger-custom-h2');
  });

  it('should trigger a delivery sequence', () => {
    environmentPage
      .clickTriggerOpen()
      .assertTriggerNextPageEnabled(false)
      .selectTriggerDelivery()
      .assertTriggerNextPageEnabled(true)
      .clickTriggerNext()
      .assertTriggerSequenceEnabled(false)
      .typeTriggerDeliveryLabels('key1=val1')
      .typeTriggerDeliveryValues('{"key2')
      .assertTriggerSequenceEnabled(false)
      .assertTriggerDeliveryValuesErrorExists(true)
      .typeTriggerDeliveryValues('": "val2"}')
      .assertTriggerDeliveryValuesErrorExists(false)
      .assertTriggerSequenceEnabled(false)
      .typeTriggerDeliveryImage('docker.io/keptn')
      .assertTriggerSequenceEnabled(false)
      .typeTriggerDeliveryTag('v0.1.2')
      .assertTriggerSequenceEnabled(true)
      .clickTriggerSequence();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should trigger an evaluation sequence with a timeframe', () => {
    environmentPage
      .clickTriggerOpen()
      .assertTriggerNextPageEnabled(false)
      .selectTriggerEvaluation()
      .assertTriggerNextPageEnabled(true)
      .clickTriggerNext()
      .assertTriggerSequenceEnabled(true)
      .selectTriggerEvaluationType(0)
      .typeTriggerEvaluationLabels('key1=val1')
      .assertTriggerSequenceEnabled(true)
      .clickTriggerStartTime()
      .selectTriggerDateTime(0, '1', '15', '0')
      .assertTriggerSequenceEnabled(true);

    environmentPage
      .typeTriggerEvaluationTimeInput('hours', '0')
      .assertTriggerSequenceEnabled(false)
      .assertTriggerEvaluationTimeframeErrorExists(true)
      .clearTriggerEvaluationTimeInput('hours')
      .typeTriggerEvaluationTimeInput('minutes', '0')
      .assertTriggerSequenceEnabled(false)
      .assertTriggerEvaluationTimeframeErrorExists(true)
      .clearTriggerEvaluationTimeInput('minutes')
      .typeTriggerEvaluationTimeInput('seconds', '59')
      .assertTriggerSequenceEnabled(false)
      .assertTriggerEvaluationTimeframeErrorExists(true)
      .clearTriggerEvaluationTimeInput('seconds')
      .typeTriggerEvaluationTimeInput('millis', '59999')
      .assertTriggerSequenceEnabled(false)
      .assertTriggerEvaluationTimeframeErrorExists(true)
      .clearTriggerEvaluationTimeInput('millis')
      .typeTriggerEvaluationTimeInput('micros', '59999999')
      .assertTriggerSequenceEnabled(false)
      .assertTriggerEvaluationTimeframeErrorExists(true)
      .clearTriggerEvaluationTimeInput('micros');

    environmentPage
      .typeTriggerEvaluationTimeInput('minutes', '5')
      .assertTriggerSequenceEnabled(true)
      .assertTriggerEvaluationTimeframeErrorExists(false)
      .clickTriggerSequence();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should trigger an evaluation sequence with a start end date', () => {
    environmentPage
      .clickTriggerOpen()
      .assertTriggerNextPageEnabled(false)
      .selectTriggerEvaluation()
      .assertTriggerNextPageEnabled(true)
      .clickTriggerNext()
      .assertTriggerSequenceEnabled(true)
      .selectTriggerEvaluationType(1)
      .typeTriggerEvaluationLabels('key1=val1')
      .assertTriggerSequenceEnabled(false);

    // End before start date error
    environmentPage
      .clickTriggerStartTime()
      .selectTriggerDateTime(1, '1', '15', '0')
      .assertTriggerSequenceEnabled(false)
      .clickTriggerEndTime()
      .selectTriggerDateTime(0, '1', '15', '0')
      .assertTriggerEvaluationDateErrorExists(true)
      .assertTriggerSequenceEnabled(false);

    // Correct date order
    environmentPage
      .clickTriggerStartTime()
      .selectTriggerDateTime(0, '1', '15', '0')
      .clickTriggerEndTime()
      .selectTriggerDateTime(1, '1', '15', '0')
      .assertTriggerEvaluationDateErrorExists(false)
      .assertTriggerSequenceEnabled(true)
      .clickTriggerSequence();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should trigger a custom sequence', () => {
    environmentPage
      .clickTriggerOpen()
      .assertTriggerNextPageEnabled(false)
      .selectTriggerCustomSequence()
      .assertTriggerNextPageEnabled(true)
      .clickTriggerNext()
      .assertTriggerSequenceEnabled(false)
      .typeTriggerCustomLabels('key1=val1')
      .assertTriggerSequenceEnabled(false)
      .selectTriggerCustomSequenceType(0)
      .assertTriggerSequenceEnabled(true)
      .clickTriggerSequence();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should open the trigger form from the sequence screen', () => {
    cy.intercept('/api/mongodb-datastore/event?keptnContext=6c98fbb0-4c40-4bff-ba9f-b20556a57c8a&project=sockshop');
    cy.visit('/project/sockshop/sequence');
    environmentPage.assertOpenTriggerSequenceExists(true).clickTriggerOpen();
    cy.url().should('include', '/project/sockshop');
    environmentPage
      .assertTriggerEntryH2Exists(true)
      .assertTriggerEntryH2HasText('Trigger a new sequence for project sockshop');
  });

  it('should have the selected stage preselected', () => {
    environmentPage.assertTriggerStageSelection(0, 'dev');
    environmentPage.assertTriggerStageSelection(1, 'staging');
    environmentPage.assertTriggerStageSelection(2, 'production');
  });

  function testNavigationFirstPart(h2Selector: string, expectedText: string, triggerSequenceEnabled: boolean): void {
    environmentPage.assertTriggerNextPageEnabled(true).clickTriggerNext();
    cy.byTestId(h2Selector).should('have.text', expectedText);
    environmentPage.assertTriggerSequenceEnabled(triggerSequenceEnabled).clickTriggerClose();
    cy.byTestId(h2Selector).should('not.exist');
    environmentPage.assertOpenTriggerSequenceExists(true).clickTriggerOpen();
  }

  function testNavigationSecondPart(h2Selector: string): void {
    environmentPage.clickTriggerNext().clickTriggerBack().assertTriggerEntryH2Exists(true);
    cy.byTestId(h2Selector).should('not.exist');
    environmentPage.clickTriggerClose();
  }
});
