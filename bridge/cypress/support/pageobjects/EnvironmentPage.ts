/// <reference types="cypress" />

import { interceptEmptyEnvironmentScreen, interceptEnvironmentScreen } from '../intercept';

class EnvironmentPage {
  public intercept(): this {
    interceptEnvironmentScreen();
    return this;
  }

  public interceptEmpty(): this {
    interceptEmptyEnvironmentScreen();
    return this;
  }

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

  public clearTriggerEvaluationTimeInput(inputName: string): this {
    cy.byTestId(`keptn-time-input-${inputName}`).clear();
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

  public assertTriggerEvaluationTimeframeErrorExists(exists: boolean): this {
    cy.byTestId('keptn-trigger-evaluation-timeframe-error').should(exists ? 'exist' : 'not.exist');
    return this;
  }

  public visit(project: string): this {
    cy.visit(`/project/${project}`).wait('@metadata');
    return this;
  }

  public clickCreateService(stage: string): this {
    cy.get('ktb-selectable-tile h2')
      .contains(stage)
      .parentsUntil('ktb-selectable-tile')
      .find('ktb-no-service-info a')
      .click();
    return this;
  }

  public selectStage(stage: string): this {
    cy.get('ktb-selectable-tile h2').contains(stage).click();
    return this;
  }

  public assertEvaluationHistoryLoadingCount(service: string, count: number): this {
    this.getServiceDetailsContainer(service)
      .find('ktb-evaluation-info dt-tag-list[aria-label="evaluation-history"] dt-tag dt-loading-spinner')
      .should('have.length', count);
    return this;
  }

  public assertEvaluationHistoryCount(service: string, count: number): this {
    const tags = this.getServiceDetailsContainer(service).find(
      'ktb-evaluation-info dt-tag-list[aria-label="evaluation-history"] dt-tag'
    );
    tags.should('have.length', count);
    if (count !== 0) {
      tags.should('have.class', 'border');
    }
    return this;
  }

  public assertEvaluationInDetails(service: string, score: number | '-', status?: 'success' | 'error'): this {
    const evaluationTag = this.getServiceDetailsContainer(service).find(
      'ktb-evaluation-info dt-tag-list[aria-label="evaluation-info"] dt-tag'
    );
    evaluationTag.should('not.have.class', 'border').should('have.text', score);
    if (status) {
      evaluationTag.should('have.class', status);
    }
    return this;
  }

  public getEvaluationHistoryURL(project: string, stage: string, service: string, limit: 5 | 6): string {
    return `/api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?filter=data.project:${project}%20AND%20data.service:${service}%20AND%20data.stage:${stage}%20AND%20source:lighthouse-service&excludeInvalidated=true&limit=${limit}`;
  }

  private getServiceDetailsContainer(service: string): Cypress.Chainable<JQuery<HTMLElement>> {
    return cy.get('ktb-stage-details ktb-expandable-tile h2').contains(service).parentsUntil('ktb-expandable-tile');
  }
}
export default EnvironmentPage;
