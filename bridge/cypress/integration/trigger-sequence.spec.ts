import EnvironmentPage from '../support/pageobjects/EnvironmentPage';
import { interceptEnvironmentScreen } from '../support/intercept';
import { TriggerSequenceSubPage } from '../support/pageobjects/TriggerSequenceSubPage';
import { SequencesPage } from '../support/pageobjects/SequencesPage';

const environmentPage = new EnvironmentPage();
const triggerSequencePage = new TriggerSequenceSubPage();
const sequencePage = new SequencesPage();
const project = 'sockshop';

describe('Trigger a sequence', () => {
  beforeEach(() => {
    interceptEnvironmentScreen();
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfoCD.mock' });

    // Sequence screen
    sequencePage.intercept();
    triggerSequencePage.visit(project);
  });

  it('should navigate through all forms and close it from everywhere properly', () => {
    // Opening of triggering component
    triggerSequencePage
      .clickOpen()
      .assertOpenTriggerSequenceExists(false)
      .assertHeadlineDefault(project)
      .assertNextPageEnabled(false)

      // Closing of triggering component from entry
      .clickClose()
      .assertTriggerSequenceFormExists(false)
      .assertOpenTriggerSequenceExists(true);

    // Delivery navigations
    triggerSequencePage
      .clickOpen()
      .selectService('carts')
      .selectStage('dev')
      .selectDelivery()
      .clickNext()
      .assertHeadlineDelivery('carts', 'dev')
      .assertTriggerSequenceEnabled(false)
      .closeAndValidate()

      .assertNextAndBackDelivery(project, 'carts', 'dev')
      .clickClose();

    // Evaluation navigations
    triggerSequencePage
      .clickOpen()
      .selectService('carts')
      .selectStage('dev')
      .selectEvaluation()
      .clickNext()
      .assertHeadlineEvaluation('carts', 'dev')
      .assertTriggerSequenceEnabled(true)
      .closeAndValidate()

      .assertNextAndBackEvaluation(project, 'carts', 'dev')
      .clickClose();

    // Custom sequence navigations
    triggerSequencePage
      .clickOpen()
      .selectService('carts')
      .selectStage('dev')
      .selectCustomSequence('delivery-direct')
      .clickNext()
      .assertHeadlineCustomSequence('delivery-direct', 'carts', 'dev')
      .assertTriggerSequenceEnabled(true)
      .closeAndValidate()

      .assertNextAndBackCustomSequence(project, 'carts', 'dev', 'delivery-direct')
      .clickClose();
  });

  it('should trigger a delivery sequence', () => {
    triggerSequencePage
      .clickOpen()
      .assertNextPageEnabled(false)
      .selectService('carts')
      .selectStage('dev')
      .selectDelivery()
      .assertNextPageEnabled(true)
      .clickNext()
      .assertTriggerSequenceEnabled(false)
      .typeDeliveryLabels('key1=val1')
      .typeDeliveryValues('{"key2')
      .assertTriggerSequenceEnabled(false)
      .assertDeliveryValuesErrorExists(true)
      .typeDeliveryValues('": "val2"}')
      .assertDeliveryValuesErrorExists(false)
      .assertTriggerSequenceEnabled(false)
      .typeDeliveryImage('docker.io/keptn')
      .assertTriggerSequenceEnabled(false)
      .typeDeliveryTag('v0.1.2')
      .assertTriggerSequenceEnabled(true)
      .clickTriggerSequence();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should trigger an evaluation sequence with a timeframe', () => {
    const currentDate = new Date();
    const currentMonthFormatted = (currentDate.getMonth() + 1).toString().padStart(2, '0');

    triggerSequencePage
      .clickOpen()
      .assertNextPageEnabled(false)
      .selectService('carts')
      .selectStage('dev')
      .selectEvaluation()
      .assertNextPageEnabled(true)
      .clickNext()
      .assertTriggerSequenceEnabled(true)
      .selectEvaluationTimeframe()
      .typeEvaluationLabels('key1=val1')
      .assertTriggerSequenceEnabled(true)
      .setStartDate(0, '1', '15', '0')
      .assertStartDateDisplayValue(`${currentDate.getFullYear()}-${currentMonthFormatted}-01 01:15:00`)
      .assertTriggerSequenceEnabled(true)

      .typeTimeframe('hours', '0')
      .assertTriggerSequenceEnabled(false)
      .assertEvaluationTimeframeErrorExists(true)
      .clearEvaluationTimeInput('hours')
      .typeTimeframe('minutes', '0')
      .assertTriggerSequenceEnabled(false)
      .assertEvaluationTimeframeErrorExists(true)
      .clearEvaluationTimeInput('minutes')
      .typeTimeframe('seconds', '59')
      .assertTriggerSequenceEnabled(false)
      .assertEvaluationTimeframeErrorExists(true)
      .clearEvaluationTimeInput('seconds')
      .typeTimeframe('millis', '59999')
      .assertTriggerSequenceEnabled(false)
      .assertEvaluationTimeframeErrorExists(true)
      .clearEvaluationTimeInput('millis')
      .typeTimeframe('micros', '59999999')
      .assertTriggerSequenceEnabled(false)
      .assertEvaluationTimeframeErrorExists(true)
      .clearEvaluationTimeInput('micros')

      .typeTimeframe('minutes', '5')
      .assertTriggerSequenceEnabled(true)
      .assertEvaluationTimeframeErrorExists(false)
      .clickTriggerSequence();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should trigger a custom sequence', () => {
    triggerSequencePage
      .clickOpen()
      .assertNextPageEnabled(false)
      .selectService('carts')
      .selectStage('dev')
      .selectCustomSequence('delivery-direct')
      .assertNextPageEnabled(true)
      .clickNext()
      .assertTriggerSequenceEnabled(true)
      .typeCustomLabels('key1=val1')
      .assertTriggerSequenceEnabled(true)
      .clickTriggerSequence();
    cy.location('pathname').should('eq', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should open the trigger form from the sequence screen', () => {
    cy.visit('/project/sockshop/sequence');
    triggerSequencePage.assertOpenTriggerSequenceExists(true).clickOpen();
    cy.location('pathname').should('eq', '/project/sockshop');
    triggerSequencePage.assertHeadlineDefault(project);
  });

  it('should have the selected stage preselected', () => {
    environmentPage.selectStage('dev');
    triggerSequencePage.assertPreSelectStage('dev');

    environmentPage.selectStage('staging');
    triggerSequencePage.assertPreSelectStage('staging');

    environmentPage.selectStage('production');
    triggerSequencePage.assertPreSelectStage('production');
  });

  it('should revert to delivery if stage is changed and does not contain any custom sequences', () => {
    triggerSequencePage
      .clickOpen()
      .selectStage('dev')
      .selectCustomSequence('delivery-direct')
      .selectStage('production')
      .assertCustomSequenceSelected(false)
      .assertDeliverySelected(true);
  });

  it('should have disabled next button if a valid custom sequence is entered and then the stage changes', () => {
    triggerSequencePage
      .clickOpen()
      .selectStage('dev')
      .selectCustomSequence('delivery-direct')
      .selectStage('staging')
      .assertCustomSequenceSelected(true)
      .assertNextPageEnabled(false);
  });

  it('should have disabled custom sequence if stage does not have custom sequences', () => {
    triggerSequencePage.clickOpen().selectStage('production').assertCustomSequenceEnabled(false);
  });

  it('should have enabled custom sequence if stage without sequences is changed to stage with sequences', () => {
    triggerSequencePage.clickOpen().selectStage('production').selectStage('dev').assertCustomSequenceEnabled(true);
  });

  it('should have disabled next button if custom sequence radio is selected but no sequence is selected', () => {
    triggerSequencePage
      .clickOpen()
      .selectStage('dev')
      .selectService('carts')
      .selectCustomSequence()
      .assertNextPageEnabled(false);
  });
});

describe('Trigger an evaluation sequence', () => {
  beforeEach(() => {
    interceptEnvironmentScreen();
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfoCD.mock' });
    sequencePage.intercept();
    triggerSequencePage
      .visit(project)
      .clickOpen()
      .selectService('carts')
      .selectStage('dev')
      .selectEvaluation()
      .clickNext();
  });

  it('should trigger an evaluation sequence with a start end date', () => {
    const currentDate = new Date();
    const currentMonthFormatted = (currentDate.getMonth() + 1).toString().padStart(2, '0');

    triggerSequencePage
      .selectEvaluationEndDate()
      .setStartDate(0, '1', '15', '0')
      .assertStartDateDisplayValue(`${currentDate.getFullYear()}-${currentMonthFormatted}-01 01:15:00`)
      .setEndDate(0, '3', '15', '0')
      .assertEndDateDisplayValue(`${currentDate.getFullYear()}-${currentMonthFormatted}-01 03:15:00`)
      .assertEvaluationDateErrorExists(false)
      .assertTriggerSequenceEnabled(true)
      .clickTriggerSequence();
    cy.url().should('include', '/project/sockshop/sequence/6c98fbb0-4c40-4bff-ba9f-b20556a57c8a/stage/dev');
  });

  it('should have disabled trigger button if end date is selected and start date is not provided', () => {
    triggerSequencePage
      .assertTriggerSequenceEnabled(true)
      .selectEvaluationEndDate()
      .typeEvaluationLabels('key1=val1')
      .assertTriggerSequenceEnabled(false);
  });

  it('should have disabled trigger button if end date is before start date', () => {
    triggerSequencePage
      .selectEvaluationEndDate()
      .setStartDate(1, '1', '15', '0')
      .assertTriggerSequenceEnabled(false)
      .setEndDate(0, '1', '15', '0')
      .assertEvaluationDateErrorExists(true)
      .assertTriggerSequenceEnabled(false);
  });

  it('should not show date error if invalid end date is changed to valid end date', () => {
    const currentDate = new Date();
    const currentMonthFormatted = (currentDate.getMonth() + 1).toString().padStart(2, '0');

    triggerSequencePage
      .selectEvaluationEndDate()
      .setStartDate(0, '2', '15', '0')
      .assertStartDateDisplayValue(`${currentDate.getFullYear()}-${currentMonthFormatted}-01 02:15:00`)
      .assertTriggerSequenceEnabled(false)
      .setEndDate(0, '1', '15', '0')
      .assertEndDateDisplayValue(`${currentDate.getFullYear()}-${currentMonthFormatted}-01 01:15:00`)
      .assertEvaluationDateErrorExists(true)
      .assertTriggerSequenceEnabled(false)

      .setEndDate(0, '3', '15', '0')
      .assertEndDateDisplayValue(`${currentDate.getFullYear()}-${currentMonthFormatted}-01 03:15:00`)
      .assertEvaluationDateErrorExists(false)
      .assertTriggerSequenceEnabled(true);
  });

  it('should not show date error if invalid start date is changed to valid start date', () => {
    triggerSequencePage
      .selectEvaluationEndDate()
      .setStartDate(0, '3', '15', '0')
      .assertTriggerSequenceEnabled(false)
      .setEndDate(0, '2', '15', '0')
      .assertEvaluationDateErrorExists(true)
      .assertTriggerSequenceEnabled(false)

      .setStartDate(0, '1', '15', '0')
      .assertEvaluationDateErrorExists(false)
      .assertTriggerSequenceEnabled(true);
  });

  it('should have disabled button if switched from valid timeframe to empty endDate', () => {
    triggerSequencePage
      .typeTimeframe('hours', '1')
      .assertTriggerSequenceEnabled(true)
      .selectEvaluationEndDate()
      .assertTriggerSequenceEnabled(false);
  });

  it('should show error if switched from valid timeframe to invalid endDate', () => {
    triggerSequencePage
      .assertTriggerSequenceEnabled(true)
      .selectEvaluationEndDate()
      .setStartDate(1, '1', '1', '1')
      .setEndDate(0, '1', '1', '1')
      .assertEvaluationDateErrorExists(true)
      .assertTriggerSequenceEnabled(false)

      .selectEvaluationTimeframe()
      .assertEvaluationDateErrorExists(false)
      .assertTriggerSequenceEnabled(true)

      .selectEvaluationEndDate()
      .assertEvaluationDateErrorExists(true)
      .assertTriggerSequenceEnabled(false);
  });

  it('should show error if switched from valid endDate to invalid timeframe', () => {
    triggerSequencePage
      .typeTimeframe('seconds', '1')
      .assertTriggerSequenceEnabled(false)
      .assertEvaluationTimeframeErrorExists(true)

      .selectEvaluationEndDate()
      .setStartDate(0, '1', '1', '1')
      .setEndDate(1, '1', '1', '1')
      .assertEvaluationDateErrorExists(false)
      .assertTriggerSequenceEnabled(true)
      .assertEvaluationTimeframeErrorExists(false)

      .selectEvaluationTimeframe()
      .assertEvaluationTimeframeErrorExists(true)
      .assertTriggerSequenceEnabled(false);
  });
});
