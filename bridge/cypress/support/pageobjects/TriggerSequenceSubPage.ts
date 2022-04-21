import EnvironmentPage from './EnvironmentPage';

enum TriggerSequenceView {
  DELIVERY,
  EVALUATION,
  CUSTOM,
}

type TimeUnits = 'hours' | 'minutes' | 'seconds' | 'micros' | 'millis';

export class TriggerSequenceSubPage {
  private environmentPage = new EnvironmentPage();

  public visit(project: string): this {
    this.environmentPage.visit(project);
    return this;
  }

  public selectDelivery(): this {
    cy.byTestId('ktb-trigger-sequence-delivery-radio').click();
    return this;
  }

  public assertDeliverySelected(status: boolean): this {
    cy.byTestId('ktb-trigger-sequence-delivery-radio')
      .find('input')
      .should(status ? 'be.checked' : 'not.be.checked');
    return this;
  }

  public selectEvaluation(): this {
    cy.byTestId('ktb-trigger-sequence-evaluation-radio').click();
    return this;
  }

  public assertEvaluationSelected(status: boolean): this {
    cy.byTestId('ktb-trigger-sequence-evaluation-radio')
      .find('input')
      .should(status ? 'be.checked' : 'not.be.checked');
    return this;
  }

  public selectCustomSequence(sequence?: string): this {
    if (sequence) {
      cy.byTestId('keptn-trigger-custom-sequence-select').dtSelect(sequence);
      return this;
    }
    cy.byTestId('ktb-trigger-sequence-custom-radio').dtCheck(true);
    return this;
  }

  public assertCustomSequenceSelected(status: boolean, sequence = ' Select ... '): this {
    cy.byTestId('ktb-trigger-sequence-custom-radio')
      .find('input')
      .should(status ? 'be.checked' : 'not.be.checked');
    cy.byTestId('keptn-trigger-custom-sequence-select').should('have.text', sequence);
    return this;
  }

  public selectService(service: string): this {
    cy.byTestId('keptn-trigger-service-selection').dtSelect(service);
    return this;
  }

  public selectStage(stage: string): this {
    cy.byTestId('keptn-trigger-stage-selection').dtSelect(stage);
    return this;
  }

  public setStartDate(calElement: number, hours: string, minutes: string, seconds: string): this {
    return this.clickStartTime().selectDateTime(calElement, hours, minutes, seconds);
  }

  public setEndDate(calElement: number, hours: string, minutes: string, seconds: string): this {
    return this.clickEndTime().selectDateTime(calElement, hours, minutes, seconds);
  }

  public assertStartDateDisplayValue(value: string): this {
    cy.byTestId('keptn-trigger-evaluation-start-date').should('have.value', value);
    return this;
  }

  public assertEndDateDisplayValue(value: string): this {
    cy.byTestId('keptn-trigger-evaluation-end-date').should('have.value', value);
    return this;
  }

  public selectDateTime(calElement: number, hours: string, minutes: string, seconds: string): this {
    cy.get('.dt-calendar-table-cell').eq(calElement).click();
    cy.byTestId('keptn-datetime-picker-submit').should('be.enabled');
    cy.byTestId('keptn-datetime-picker-time').byTestId('keptn-time-input-hours').type(hours);
    cy.byTestId('keptn-datetime-picker-submit').should('be.disabled');
    cy.byTestId('keptn-datetime-picker-time').byTestId('keptn-time-input-minutes').type(minutes);
    cy.byTestId('keptn-datetime-picker-submit').should('be.disabled');
    cy.byTestId('keptn-datetime-picker-time').byTestId('keptn-time-input-seconds').type(seconds);
    cy.byTestId('keptn-datetime-picker-submit').should('be.enabled').click();
    return this;
  }

  public clickOpen(): this {
    cy.byTestId('keptn-trigger-button-open').click().wait('@customSequences');
    return this;
  }

  public assertOpenTriggerSequenceExists(exists: boolean): this {
    cy.byTestId('keptn-trigger-button-open').should(exists ? 'exist' : 'not.exist');
    return this;
  }

  public clickClose(): this {
    cy.byTestId('keptn-trigger-button-close').click();
    return this;
  }

  public clickNext(): this {
    cy.byTestId('keptn-trigger-button-next').click();
    return this;
  }

  public clickBack(): this {
    cy.byTestId('keptn-trigger-button-back').click();
    return this;
  }

  private clickStartTime(): this {
    cy.byTestId('keptn-trigger-button-starttime').click();
    return this;
  }

  private clickEndTime(): this {
    cy.byTestId('keptn-trigger-button-end-date').click();
    return this;
  }

  public clickTriggerSequence(): this {
    cy.byTestId('keptn-trigger-button-trigger').click();
    return this;
  }

  public typeDeliveryLabels(value: string): this {
    cy.byTestId('keptn-trigger-delivery-labels').type(value);
    return this;
  }

  public typeDeliveryValues(value: string): this {
    cy.byTestId('keptn-trigger-delivery-values').type(value);
    return this;
  }

  public typeDeliveryImage(value: string): this {
    cy.byTestId('keptn-trigger-delivery-image').type(value);
    return this;
  }

  public typeDeliveryTag(value: string): this {
    cy.byTestId('keptn-trigger-delivery-tag').type(value);
    return this;
  }

  public typeEvaluationLabels(value: string): this {
    cy.byTestId('keptn-trigger-evaluation-labels').type(value);
    return this;
  }

  public typeTimeframe(inputName: TimeUnits, value: string): this {
    cy.byTestId(`keptn-time-input-${inputName}`).type(value);
    return this;
  }

  public clearEvaluationTimeInput(inputName: TimeUnits): this {
    cy.byTestId(`keptn-time-input-${inputName}`).clear();
    return this;
  }

  public typeCustomLabels(value: string): this {
    cy.byTestId('keptn-trigger-custom-labels').type(value);
    return this;
  }

  public selectEvaluationEndDate(): this {
    cy.byTestId('ktb-trigger-sequence-radio-end-date').click();
    return this;
  }

  public selectEvaluationTimeframe(): this {
    cy.byTestId('ktb-trigger-sequence-radio-timeframe').click();
    return this;
  }

  public assertCustomSequenceEnabled(status: boolean): this {
    cy.byTestId('ktb-trigger-sequence-custom-radio')
      .find('input')
      .should(status ? 'be.enabled' : 'be.disabled');
    cy.byTestId('keptn-trigger-custom-sequence-select').should(
      status ? 'not.have.class' : 'have.class',
      'dt-select-disabled'
    );
    return this;
  }

  public assertStageSelected(stage: string): this {
    cy.byTestId('keptn-trigger-stage-selection').find('.dt-select-value-text > span').should('have.text', stage);
    return this;
  }

  public assertPreSelectStage(stage: string): this {
    return this.clickOpen().assertStageSelected(stage).clickClose();
  }

  public assertTriggerSequenceEnabled(enabled: boolean): this {
    cy.byTestId('keptn-trigger-button-trigger').should(enabled ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public assertNextPageEnabled(enabled: boolean): this {
    cy.byTestId('keptn-trigger-button-next').should(enabled ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public assertDeliveryValuesErrorExists(exists: boolean): this {
    cy.byTestId('keptn-trigger-delivery-values-error').should(exists ? 'exist' : 'not.exist');
    return this;
  }

  public assertEvaluationDateErrorExists(exists: boolean): this {
    cy.byTestId('keptn-trigger-evaluation-date-error').should(exists ? 'exist' : 'not.exist');
    return this;
  }

  public assertEvaluationTimeframeErrorExists(exists: boolean): this {
    cy.byTestId('keptn-trigger-evaluation-timeframe-error').should(exists ? 'exist' : 'not.exist');
    return this;
  }

  private assertHeadline(text: string): this {
    cy.byTestId('ktb-trigger-headline').should('have.length', 1).should('have.text', text);
    return this;
  }

  public assertHeadlineEvaluation(service: string, stage: string): this {
    return this.assertHeadline(`Trigger an evaluation for ${service} in ${stage}`);
  }

  public assertHeadlineDelivery(service: string, stage: string): this {
    return this.assertHeadline(`Trigger a delivery for ${service} in ${stage}`);
  }

  public assertHeadlineDefault(projectName: string): this {
    return this.assertHeadline(`Trigger a new sequence for project ${projectName}`);
  }

  public assertHeadlineCustomSequence(sequence: string, service: string, stage: string): this {
    return this.assertHeadline(` Trigger a ${sequence} sequence for ${service} in ${stage} `);
  }

  public assertTriggerSequenceFormExists(status: boolean): this {
    cy.get('ktb-trigger-sequence-component').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertNextAndBackDelivery(project: string, service: string, stage: string): this {
    return this.assertNextAndBack(project, service, stage, TriggerSequenceView.DELIVERY);
  }

  public assertNextAndBackEvaluation(project: string, service: string, stage: string): this {
    return this.assertNextAndBack(project, service, stage, TriggerSequenceView.EVALUATION);
  }

  public assertNextAndBackCustomSequence(project: string, service: string, stage: string, sequence: string): this {
    return this.assertNextAndBack(project, service, stage, TriggerSequenceView.CUSTOM, sequence);
  }

  private assertNextAndBack(
    project: string,
    service: string,
    stage: string,
    view: TriggerSequenceView,
    sequence?: string
  ): this {
    this.clickOpen().selectService(service).selectStage(stage);
    switch (view) {
      case TriggerSequenceView.DELIVERY:
        return this.selectDelivery()
          .clickNext()
          .clickBack()
          .assertDeliverySelected(true)
          .assertHeadlineDefault(project);
      case TriggerSequenceView.EVALUATION:
        return this.selectEvaluation()
          .clickNext()
          .clickBack()
          .assertEvaluationSelected(true)
          .assertHeadlineDefault(project);
      default:
        return this.selectCustomSequence(sequence ?? '')
          .clickNext()
          .clickBack()
          .assertCustomSequenceSelected(true, sequence)
          .assertHeadlineDefault(project);
    }
  }

  public closeAndValidate(): this {
    return this.clickClose().assertTriggerSequenceFormExists(false).assertOpenTriggerSequenceExists(true);
  }
}
