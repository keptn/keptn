/// <reference types="cypress" />

class EnvironmentPage {
  public selectTriggerDelivery(): this {
    cy.byTestId('keptn-trigger-sequence-selection').children().first().click();
    this.selectTriggerServiceAndStage();
    return this;
  }

  public selectTriggerEvaluation(): this {
    cy.byTestId('keptn-trigger-sequence-selection').children().eq(1).click();
    this.selectTriggerServiceAndStage();
    return this;
  }

  public selectTriggerCustomSequence(): this {
    cy.wait(500);
    cy.byTestId('keptn-trigger-sequence-selection').children().eq(2).click();
    this.selectTriggerServiceAndStage();
    return this;
  }

  public selectTriggerServiceAndStage(): this {
    cy.byTestId('keptn-trigger-service-selection').click();
    cy.wait(500);
    cy.get('dt-option').eq(0).click();
    cy.byTestId('keptn-trigger-stage-selection').click();
    cy.wait(500);
    cy.get('dt-option').eq(0).click();
    return this;
  }

  public selectTriggerDateTime(calElement: number, hours: string, minutes: string, seconds: string): this {
    cy.get('.dt-calendar-header-button-prev-month').click();
    cy.get('.dt-calendar-table-cell').eq(calElement).click();
    cy.byTestId('keptn-datetime-picker-submit').should('be.enabled');
    cy.get('[uitestid="keptn-datetime-picker-time"] [uitestid="keptn-time-input-hours"]').type(hours);
    cy.byTestId('keptn-datetime-picker-submit').should('be.disabled');
    cy.get('[uitestid="keptn-datetime-picker-time"] [uitestid="keptn-time-input-minutes"]').type(minutes);
    cy.byTestId('keptn-datetime-picker-submit').should('be.disabled');
    cy.get('[uitestid="keptn-datetime-picker-time"] [uitestid="keptn-time-input-seconds"]').type(seconds);
    cy.byTestId('keptn-datetime-picker-submit').should('be.enabled').click();
    return this;
  }

  public clickTriggerOpen(): this {
    cy.byTestId('keptn-trigger-button-open').click();
    return this;
  }

  public clickTriggerClose(): this {
    cy.byTestId('keptn-trigger-button-close').click();
    return this;
  }

  public clickTriggerNext(): this {
    cy.byTestId('keptn-trigger-button-next').click();
    return this;
  }

  public clickTriggerBack(): this {
    cy.byTestId('keptn-trigger-button-back').click();
    return this;
  }

  public clickTriggerStartTime(): this {
    cy.byTestId('keptn-trigger-button-starttime').click();
    return this;
  }

  public clickTriggerEndTime(): this {
    cy.byTestId('keptn-trigger-button-endtime').click();
    return this;
  }

  public clickTriggerSequence(): this {
    cy.byTestId('keptn-trigger-button-trigger').click();
    return this;
  }

  public typeTriggerDeliveryLabels(value: string): this {
    cy.byTestId('keptn-trigger-delivery-labels').type(value);
    return this;
  }

  public typeTriggerDeliveryValues(value: string): this {
    cy.byTestId('keptn-trigger-delivery-values').type(value);
    return this;
  }

  public typeTriggerDeliveryImage(value: string): this {
    cy.byTestId('keptn-trigger-delivery-image').type(value);
    return this;
  }

  public typeTriggerDeliveryTag(value: string): this {
    cy.byTestId('keptn-trigger-delivery-tag').type(value);
    return this;
  }

  public typeTriggerEvaluationLabels(value: string): this {
    cy.byTestId('keptn-trigger-evaluation-labels').type(value);
    return this;
  }

  public typeTriggerEvaluationTimeInput(inputName: string, value: string): this {
    cy.byTestId(`keptn-time-input-${inputName}`).type(value);
    return this;
  }

  public typeTriggerCustomLabels(value: string): this {
    cy.byTestId('keptn-trigger-custom-labels').type(value);
    return this;
  }

  public selectTriggerEvaluationType(elIndex: number): this {
    cy.byTestId('keptn-trigger-evaluation-type').children().eq(elIndex).click();
    return this;
  }

  public selectTriggerCustomSequenceType(elIndex: number): this {
    cy.byTestId('keptn-trigger-custom-sequence').click();
    cy.wait(500);
    cy.get('dt-option').eq(elIndex).click();
    return this;
  }

  public assertTriggerStageSelection(elem: number, expectedText: string): this {
    cy.get('.stage-list .ktb-selectable-tile').eq(elem).find('h2').click();
    cy.byTestId('keptn-trigger-button-open').click();
    cy.get('[uitestid="keptn-trigger-stage-selection"] .dt-select-value-text > span').should('have.text', expectedText);
    cy.byTestId('keptn-trigger-button-close').click();
    return this;
  }

  public assertTriggerSequenceEnabled(enabled: boolean): this {
    cy.byTestId('keptn-trigger-button-trigger').should(enabled ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public assertOpenTriggerSequenceExists(exists: boolean): this {
    cy.byTestId('keptn-trigger-button-open').should(exists ? 'exist' : 'not.exist');
    return this;
  }

  public assertTriggerNextPageEnabled(enabled: boolean): this {
    cy.byTestId('keptn-trigger-button-next').should(enabled ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public assertTriggerEntryH2HasText(text: string): this {
    cy.byTestId('keptn-trigger-entry-h2').should('have.text', text);
    return this;
  }

  public assertTriggerEntryH2Exists(exists: boolean): this {
    cy.byTestId('keptn-trigger-entry-h2').should(exists ? 'exist' : 'not.exist');
    return this;
  }

  public assertTriggerDeliveryValuesErrorExists(exists: boolean): this {
    cy.byTestId('keptn-trigger-delivery-values-error').should(exists ? 'exist' : 'not.exist');
    return this;
  }

  public assertTriggerEvaluationDateErrorExists(exists: boolean): this {
    cy.byTestId('keptn-trigger-evaluation-date-error').should(exists ? 'exist' : 'not.exist');
    return this;
  }
}
export default EnvironmentPage;
