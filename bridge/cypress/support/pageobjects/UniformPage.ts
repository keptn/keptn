/// <reference types="cypress" />

class UniformPage {
  UNIFORM_NAME_LOC = 'span.ng-star-inserted';
  UNIFORM_SUBMENU_LOC = 'integration-submenu';
  UNIFORM_INTEGRATION_TABLE_LOC = 'keptn-uniform-integrations-table';
  SUBSCRIPTION_EXP_HEADER_LOC = 'ktb-expandable-tile h3';
  SUBSCRIPTION_DETAILS_LOC = 'ktb-expandable-tile dt-expandable-panel ktb-subscription-item';
  SUBSCRIPTION_EXPANDABLE_LOC = 'ktb-expandable-tile';
  EDIT_WEBHOOK_PAYLOAD_ID = 'edit-webhook-field-payload';
  EDIT_WEBHOOK_FIELD_HEADER_NAME_ID = 'edit-webhook-field-headerName';
  EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID = 'edit-webhook-field-headerValue';
  UPDATE_SUBSCRIPTION_BUTTON_ID = 'updateSubscriptionButton';
  EDIT_SUBSCRIPTION_FIELD_GLOBAL_ID = 'edit-subscription-field-isGlobal';
  EDIT_SUBSCRIPTION_FIELD_TASK_ID = 'edit-subscription-field-task';
  EDIT_SUBSCRIPTION_FIELD_SUFFIX_ID = 'edit-subscription-field-suffix';
  EDIT_WEBHOOK_FIELD_URL_ID = 'edit-webhook-field-url';
  EDIT_WEBHOOK_SECRET_SELECTOR_URL = 'secret-url';
  EDIT_WEBHOOK_EVENT_SELECTOR_URL = 'event-url';
  EDIT_WEBHOOK_SECRET_SELECTOR_PAYLOAD = 'secret-payload';
  EDIT_WEBHOOK_EVENT_SELECTOR_PAYLOAD = 'event-payload';
  EDIT_WEBHOOK_SECRET_SELECTOR_HEADER = 'secret-header';
  EDIT_WEBHOOK_EVENT_SELECTOR_HEADER = 'event-header';
  EDIT_WEBHOOK_FILTER = 'edit-subscription-field-filterStageService';
  EDIT_WEBHOOK_METHOD = 'edit-webhook-field-method';
  EDIT_WEBHOOK_PROXY = 'edit-webhook-field-proxy';
  SEND_FINISHED_BUTTONS = 'edit-webhook-field-sendFinished';
  SEND_STARTED_BUTTONS = 'edit-webhook-field-sendStarted';

  public visit(project: string): this {
    cy.visit(`/project/${project}/settings/uniform/integrations`).wait('@metadata');
    return this;
  }

  public visitAdd(integrationID: string): this {
    cy.visit(`/project/sockshop/settings/uniform/integrations/${integrationID}/subscriptions/add`).wait('@metadata');
    return this;
  }

  public visitEdit(integrationID: string, subscriptionId: string): this {
    cy.visit(
      `/project/sockshop/settings/uniform/integrations/${integrationID}/subscriptions/${subscriptionId}/edit`
    ).wait('@metadata');
    return this;
  }

  public addSubscription(): this {
    cy.byTestId('addSubscriptionButton').click();
    return this;
  }

  public editSubscription(index: number): this {
    cy.byTestId('subscriptionEditButton').eq(index).click();
    return this;
  }

  public deleteSubscription(index: number): this {
    cy.byTestId('subscriptionDeleteButton').eq(index).click();
    return this;
  }

  public assertIsUpdateButtonEnabled(isEnabled: boolean): this {
    cy.byTestId(this.UPDATE_SUBSCRIPTION_BUTTON_ID).should(isEnabled ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public assertIsSendStartedButtonsEnabled(isEnabled: boolean): this {
    cy.byTestId(this.SEND_STARTED_BUTTONS)
      .find('input.dt-radio-input')
      .eq(0)
      .should(isEnabled ? 'be.enabled' : 'be.disabled');
    cy.byTestId(this.SEND_STARTED_BUTTONS)
      .find('input.dt-radio-input')
      .eq(1)
      .should(isEnabled ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public assertIsSendStarted(status: boolean): this {
    cy.byTestId(this.SEND_STARTED_BUTTONS)
      .find('input.dt-radio-input')
      .eq(0)
      .should(status ? 'be.checked' : 'not.be.checked');
    cy.byTestId(this.SEND_STARTED_BUTTONS)
      .find('input.dt-radio-input')
      .eq(1)
      .should(status ? 'not.be.checked' : 'be.checked');
    return this;
  }

  public assertIsSendFinishedButtonsEnabled(isEnabled: boolean): this {
    cy.byTestId(this.SEND_FINISHED_BUTTONS)
      .find('input.dt-radio-input')
      .eq(0)
      .should(isEnabled ? 'be.enabled' : 'be.disabled');
    cy.byTestId(this.SEND_FINISHED_BUTTONS)
      .find('input.dt-radio-input')
      .eq(1)
      .should(isEnabled ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public assertIsSendFinished(status: boolean): this {
    cy.byTestId(this.SEND_FINISHED_BUTTONS)
      .find('input.dt-radio-input')
      .eq(0)
      .should(status ? 'be.checked' : 'not.be.checked');
    cy.byTestId(this.SEND_FINISHED_BUTTONS)
      .find('input.dt-radio-input')
      .eq(1)
      .should(status ? 'not.be.checked' : 'be.checked');
    return this;
  }

  public switchIsGlobalState(): this {
    cy.byTestId(this.EDIT_SUBSCRIPTION_FIELD_GLOBAL_ID).click();
    return this;
  }

  public setTaskPrefix(selection: string): this {
    cy.byTestId(this.EDIT_SUBSCRIPTION_FIELD_TASK_ID).dtSelect(selection);
    return this;
  }

  public setTaskSuffix(selection: string): this {
    cy.byTestId(this.EDIT_SUBSCRIPTION_FIELD_SUFFIX_ID).dtSelect(selection);
    return this;
  }

  public appendStages(...stages: string[]): this {
    cy.byTestId(this.EDIT_WEBHOOK_FILTER)
      .find('input')
      .focus()
      .type(stages.map((stage) => `Stage{enter}${stage}{enter}`).join(''))
      .clickOutside();
    return this;
  }

  public appendServices(...services: string[]): this {
    cy.byTestId(this.EDIT_WEBHOOK_FILTER)
      .find('input')
      .focus()
      .type(services.map((service) => `Service{enter}${service}{enter}`).join(''))
      .clickOutside();
    return this;
  }

  public setWebhookMethod(method: string): this {
    cy.byTestId(this.EDIT_WEBHOOK_METHOD).find('dt-select').focus().type(method);
    return this;
  }

  public appendURL(content: string): this {
    cy.byTestId(this.EDIT_WEBHOOK_FIELD_URL_ID).find('textarea').focus().type(content);
    return this;
  }

  public appendPayload(content: string): this {
    cy.byTestId(this.EDIT_WEBHOOK_PAYLOAD_ID).find('textarea').focus().type(content);
    return this;
  }

  public appendProxy(proxy: string): this {
    cy.byTestId(this.EDIT_WEBHOOK_PROXY).find('input').focus().type(proxy);
    return this;
  }

  public openSecretSelectorURL(): this {
    cy.byTestId(this.EDIT_WEBHOOK_SECRET_SELECTOR_URL).find('button').click();
    return this;
  }

  public openSecretSelectorPayload(): this {
    cy.byTestId(this.EDIT_WEBHOOK_SECRET_SELECTOR_PAYLOAD).find('button').click();
    return this;
  }

  public openSecretSelectorHeader(index = 0): this {
    cy.byTestId(this.EDIT_WEBHOOK_SECRET_SELECTOR_HEADER).eq(index).find('button').click();
    return this;
  }

  public openEventSelectorURL(): this {
    cy.byTestId(this.EDIT_WEBHOOK_EVENT_SELECTOR_URL).find('button').click();
    return this;
  }

  public openEventSelectorPayload(): this {
    cy.byTestId(this.EDIT_WEBHOOK_EVENT_SELECTOR_PAYLOAD).find('button').click();
    return this;
  }

  public openEventSelectorHeader(index = 0): this {
    cy.byTestId(this.EDIT_WEBHOOK_EVENT_SELECTOR_HEADER).eq(index).find('button').click();
    return this;
  }

  public selectFirstItemOfVariableSelector(): this {
    return this.selectItemOfSelector([0], 1);
  }

  public appendHeaderName(content: string, index = 0): this {
    cy.byTestId(this.EDIT_WEBHOOK_FIELD_HEADER_NAME_ID).eq(index).find('input').focus().type(content);
    return this;
  }

  public appendHeaderValue(content: string, index = 0): this {
    cy.byTestId(this.EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID).eq(index).find('input').focus().type(content);
    return this;
  }

  public selectItemOfSelector(shallExpand: number[], clickIndex: number): this {
    for (const expand of shallExpand) {
      cy.get('ktb-tree-list-select dt-tree-table-toggle-cell').eq(expand).click();
    }
    cy.get('ktb-tree-list-select dt-tree-table-toggle-cell').eq(clickIndex).click();
    return this;
  }

  public assertIsGlobalChecked(status: boolean): this {
    cy.byTestId(this.EDIT_SUBSCRIPTION_FIELD_GLOBAL_ID)
      .find('dt-checkbox')
      .should(status ? 'have.class' : 'not.have.class', 'dt-checkbox-checked');
    return this;
  }

  public taskPrefixEquals(text: string): this {
    cy.byTestId(this.EDIT_SUBSCRIPTION_FIELD_TASK_ID).find('dt-select').should('have.text', text);
    return this;
  }

  public taskSuffixEquals(text: string): this {
    cy.byTestId(this.EDIT_SUBSCRIPTION_FIELD_SUFFIX_ID).find('dt-select').should('have.text', text);
    return this;
  }

  public filterTagsLengthEquals(length: number): this {
    cy.get('.dt-filter-field-tag-container').should('have.length', length);
    return this;
  }

  public shouldHaveStages(stages: string[]): this {
    const tags = cy.get('.dt-filter-field-tag-container');
    for (const stage of stages) {
      tags.should('contain.text', `Stage${stage}`);
    }
    return this;
  }

  public shouldHaveServices(services: string[]): this {
    const tags = cy.get('.dt-filter-field-tag-container');
    for (const service of services) {
      tags.should('contain.text', `Service${service}`);
    }
    return this;
  }

  public shouldHaveTaskSuffixes(suffixes: string[]): this {
    cy.byTestId(this.EDIT_SUBSCRIPTION_FIELD_SUFFIX_ID).find('dt-select').click();
    const items = cy.get('.dt-select-content').find('dt-option');
    // eslint-disable-next-line promise/catch-or-return
    items
      .should('have.length', suffixes.length)
      .then(($els) => Cypress._.map(Cypress.$.makeArray($els), 'innerText'))
      .should('deep.equal', suffixes);

    return this;
  }

  public update(): this {
    cy.byTestId(this.UPDATE_SUBSCRIPTION_BUTTON_ID).click();
    return this;
  }

  public clearFilter(): this {
    cy.get('.dt-filter-field-clear-all-button').click();
    return this;
  }

  public isGlobalEnabled(status: boolean): this {
    cy.byTestId(this.EDIT_SUBSCRIPTION_FIELD_GLOBAL_ID)
      .get('input')
      .should(status ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public isCreateButton(): this {
    cy.byTestId(this.UPDATE_SUBSCRIPTION_BUTTON_ID).should('contain.text', 'Create subscription');
    return this;
  }

  public isUpdateButton(): this {
    cy.byTestId(this.UPDATE_SUBSCRIPTION_BUTTON_ID).should('contain.text', 'Update subscription');
    return this;
  }

  public selectIntegration(name: string): this {
    cy.byTestId(this.UNIFORM_INTEGRATION_TABLE_LOC).contains('dt-cell', name).click();
    return this;
  }

  public assertHasIntegrationErrorIndicator(name: string, status: boolean): this {
    cy.byTestId(this.UNIFORM_INTEGRATION_TABLE_LOC)
      .contains('dt-cell', name)
      .find('.notification-indicator-text')
      .should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertIndicatorsShowing(count: number): this {
    cy.get('.notification-indicator').should('have.length', count);
    return this;
  }

  public assertIndicatorsTextShowing(count: number): this {
    cy.get('.notification-indicator-text').should('have.length', count);
    return this;
  }

  public assertErrorEventsShowing(count: number): this {
    cy.get('ktb-uniform-registration-logs .notification-indicator-text').should('have.length', count);
    return this;
  }

  public assertIntegrationErrorCount(name: string, count: number): this {
    cy.byTestId(this.UNIFORM_INTEGRATION_TABLE_LOC)
      .contains('dt-cell', name)
      .find('.notification-indicator-text')
      .should('have.text', count);
    return this;
  }

  public addHeader(): this {
    cy.byTestId('ktb-webhook-settings-add-header-button').click();
    return this;
  }

  public assertURL(content: string): this {
    cy.byTestId(this.EDIT_WEBHOOK_FIELD_URL_ID).find('textarea').should('have.value', content);
    return this;
  }

  public assertPayload(content: string): this {
    cy.byTestId(this.EDIT_WEBHOOK_PAYLOAD_ID).find('textarea').should('have.value', content);
    return this;
  }

  public assertHeaderValue(index: number, content: string): this {
    cy.byTestId(this.EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID).find('input').eq(index).should('have.value', content);
    return this;
  }
}
export default UniformPage;
