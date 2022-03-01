/// <reference types="cypress" />

import Chainable = Cypress.Chainable;

class EnvironmentPage {
  public openTriggerSequence(): void {
    cy.byTestId('keptn-trigger-button-open').click();
  }

  public selectTriggerDelivery(): void {
    cy.byTestId('keptn-trigger-sequence-selection').children().first().click();
    this.selectTriggerServiceAndStage();
  }

  public selectTriggerEvaluation(): void {
    cy.byTestId('keptn-trigger-sequence-selection').children().eq(1).click();
    this.selectTriggerServiceAndStage();
  }

  public selectTriggerCustomSequence(): void {
    cy.wait(500);
    cy.byTestId('keptn-trigger-sequence-selection').children().eq(2).click();
    this.selectTriggerServiceAndStage();
  }

  public selectTriggerServiceAndStage(): void {
    cy.byTestId('keptn-trigger-service-selection').click();
    cy.wait(500);
    cy.get('dt-option').eq(0).click();
    cy.byTestId('keptn-trigger-stage-selection').click();
    cy.wait(500);
    cy.get('dt-option').eq(0).click();
  }

  public selectTriggerDateTime(calElement: number): void {
    cy.get('.dt-calendar-header-button-prev-month').click();
    cy.get('.dt-calendar-table-cell').eq(calElement).click();
    cy.byTestId('keptn-datetime-picker-submit').should('be.enabled');
    cy.get('[uitestid="keptn-datetime-picker-time"] [uitestid="keptn-time-input-hours"]').type('1');
    cy.byTestId('keptn-datetime-picker-submit').should('be.disabled');
    cy.get('[uitestid="keptn-datetime-picker-time"] [uitestid="keptn-time-input-minutes"]').type('15');
    cy.byTestId('keptn-datetime-picker-submit').should('be.disabled');
    cy.get('[uitestid="keptn-datetime-picker-time"] [uitestid="keptn-time-input-seconds"]').type('0');
    cy.byTestId('keptn-datetime-picker-submit').should('be.enabled').click();
  }

  public clickTriggerClose(): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-button-close').click();
  }

  public clickTriggerNext(): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-button-next').click();
  }

  public clickTriggerBack(): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-button-back').click();
  }

  public clickTriggerStartTime(): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-button-starttime').click();
  }

  public clickTriggerEndTime(): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-button-endtime').click();
  }

  public typeTriggerDeliveryLabels(value: string): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-delivery-labels').type(value);
  }

  public typeTriggerDeliveryValues(value: string): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-delivery-values').type(value);
  }

  public typeTriggerDeliveryImage(value: string): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-delivery-image').type(value);
  }

  public typeTriggerDeliveryTag(value: string): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-delivery-tag').type(value);
  }

  public typeTriggerEvaluationLabels(value: string): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-evaluation-labels').type(value);
  }

  public typeTriggerEvaluationTimeInput(inputName: string, value: string): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId(`keptn-time-input-${inputName}`).type(value);
  }

  public typeTriggerCustomLabels(value: string): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-custom-labels').type(value);
  }

  public selectTriggerEvaluationType(elIndex: number): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-evaluation-type').children().eq(elIndex).click();
  }

  public selectTriggerCustomSequenceType(elIndex: number): void {
    cy.byTestId('keptn-trigger-custom-sequence').click();
    cy.wait(500);
    cy.get('dt-option').eq(elIndex).click();
  }

  public assertTriggerStageSelection(elem: number, expectedText: string): void {
    cy.get('.stage-list .ktb-selectable-tile').eq(elem).find('h2').click();
    cy.byTestId('keptn-trigger-button-open').click();
    cy.get('[uitestid="keptn-trigger-stage-selection"] .dt-select-value-text > span').should('have.text', expectedText);
    cy.byTestId('keptn-trigger-button-close').click();
  }

  public assertTriggerSequenceEnabled(enabled: boolean): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-button-trigger').should(enabled ? 'be.enabled' : 'be.disabled');
  }

  public assertOpenTriggerSequenceExists(exists: boolean): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-button-open').should(exists ? 'exist' : 'not.exist');
  }

  public assertTriggerNextPageEnabled(enabled: boolean): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-button-next').should(enabled ? 'be.enabled' : 'be.disabled');
  }

  public assertTriggerEntryH2(chainer: string, value: string | undefined = undefined): Chainable<JQuery<HTMLElement>> {
    if (value) {
      return cy.byTestId('keptn-trigger-entry-h2').should(chainer, value);
    }
    return cy.byTestId('keptn-trigger-entry-h2').should(chainer);
  }

  public assertTriggerDeliveryValuesErrorExists(exists: boolean): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-delivery-values-error').should(exists ? 'exist' : 'not.exist');
  }

  public assertTriggerEvaluationDateErrorExists(exists: boolean): Chainable<JQuery<HTMLElement>> {
    return cy.byTestId('keptn-trigger-evaluation-date-error').should(exists ? 'exist' : 'not.exist');
  }
}
export default EnvironmentPage;
