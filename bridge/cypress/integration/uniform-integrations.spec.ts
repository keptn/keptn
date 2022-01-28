import { interceptIntegrations } from '../support/intercept';
import UniformPage from '../support/pageobjects/UniformPage';

describe('Integrations', () => {
  const uniformPage = new UniformPage();

  beforeEach(() => {
    interceptIntegrations();
    cy.visit('/project/sockshop/uniform/services');
  });

  it('should be on page uniform', () => {
    const buttons = cy.byTestId(uniformPage.UNIFORM_SUBMENU_LOC).find('button');
    buttons.first().should('have.class', 'active');
  });

  it('should select an integration and show related subscriptions', () => {
    // given, when
    cy.byTestId(uniformPage.UNIFORM_INTEGRATION_TABLE_LOC).get('dt-row').first().click();

    // then
    cy.get(uniformPage.SUBSCRIPTION_EXP_HEADER_LOC).first().should('have.text', 'Subscriptions');
    cy.get(uniformPage.SUBSCRIPTION_EXP_HEADER_LOC).last().should('have.text', 'Error events');
  });

  it('should have disabled buttons for a subscription, when subscription id is not given', () => {
    // given, when
    cy.intercept('/api/uniform/registration', { fixture: 'registration-old-format.mock' });
    cy.byTestId(uniformPage.UNIFORM_INTEGRATION_TABLE_LOC).find('dt-row').eq(1).click();
    const editButton = cy.get(uniformPage.SUBSCRIPTION_EXPANDABLE_LOC).first().byTestId('subscriptionEditButton');
    const deleteButton = cy.get(uniformPage.SUBSCRIPTION_EXPANDABLE_LOC).first().byTestId('subscriptionDeleteButton');

    // then
    editButton.should('be.disabled');
    deleteButton.should('be.disabled');
  });

  it('should add a simple subscription', () => {
    // given, when
    cy.byTestId(uniformPage.UNIFORM_INTEGRATION_TABLE_LOC).find('dt-row').eq(1).click();
    cy.get(uniformPage.SUBSCRIPTION_DETAILS_LOC).should('have.length', 1);
    cy.get(uniformPage.SUBSCRIPTION_EXPANDABLE_LOC)
      .first()
      .find('dt-expandable-panel')
      .byTestId('addSubscriptionButton')
      .click();

    // Fill in all form fields
    cy.byTestId(uniformPage.EDIT_SUBSCRIPTION_FIELD_GLOBAL_ID)
      .click()
      .find('dt-checkbox')
      .should('have.class', 'dt-checkbox-checked');
    cy.byTestId(uniformPage.EDIT_SUBSCRIPTION_FIELD_TASK_ID).find('dt-select').focus().type('dep');
    cy.byTestId(uniformPage.EDIT_SUBSCRIPTION_FIELD_SUFFIX_ID).find('dt-select').focus().type('fin');
    cy.byTestId('edit-subscription-field-filterStageService')
      .find('input')
      .focus()
      .type('St{enter}de{enter}Ser{enter}cart{enter}');
    cy.byTestId(uniformPage.UPDATE_SUBSCRIPTION_BUTTON_ID).click();

    // then
    // It should redirect after successfully sending the subscription
    cy.location('pathname').should(
      'eq',
      '/project/sockshop/uniform/services/355311a7bec3f35bf3abc2484ab09bcba8e2b297/subscriptions/add'
    );
  });

  it('should add a webhook subscription', () => {
    cy.byTestId(uniformPage.UNIFORM_INTEGRATION_TABLE_LOC).find('dt-row').eq(0).click();
    cy.get(uniformPage.SUBSCRIPTION_EXPANDABLE_LOC)
      .first()
      .find('dt-expandable-panel')
      .byTestId('addSubscriptionButton')
      .click();

    // then
    cy.get('h2').eq(1).should('have.text', 'Webhook configuration');
    cy.byTestId(uniformPage.EDIT_SUBSCRIPTION_FIELD_GLOBAL_ID).should('not.exist');

    cy.byTestId(uniformPage.EDIT_SUBSCRIPTION_FIELD_TASK_ID).find('dt-select').focus().type('dep');
    cy.byTestId(uniformPage.EDIT_SUBSCRIPTION_FIELD_SUFFIX_ID).find('dt-select').focus().type('fin');
    cy.byTestId('edit-webhook-field-method').find('dt-select').focus().type('GET');

    // URL: insert text, add secret and add event variable
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_URL_ID).find('textarea').focus().type('https://example.com?secret=');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_SECRET_SELECTOR_URL).find('button').click();
    selectFirstItemOfVariableSelector();
    cy.byTestId(uniformPage.EDIT_WEBHOOK_EVENT_SELECTOR_URL).find('button').click();
    selectFirstItemOfVariableSelector();
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_URL_ID)
      .find('textarea')
      .should('have.value', 'https://example.com?secret={{.secret.SecretA.key1}}{{.data.project}}');

    // Payload: insert text, add secret and add event variable
    cy.byTestId(uniformPage.EDIT_WEBHOOK_PAYLOAD_ID).find('textarea').focus().type("{id: '123456789', secret: ");
    cy.byTestId(uniformPage.EDIT_WEBHOOK_SECRET_SELECTOR_PAYLOAD).find('button').click();
    selectFirstItemOfVariableSelector();
    cy.byTestId(uniformPage.EDIT_WEBHOOK_EVENT_SELECTOR_PAYLOAD).find('button').click();
    selectFirstItemOfVariableSelector();
    cy.byTestId(uniformPage.EDIT_WEBHOOK_PAYLOAD_ID).find('textarea').focus().type('}');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_PAYLOAD_ID)
      .find('textarea')
      .should('have.value', "{id: '123456789', secret: {{.secret.SecretA.key1}}{{.data.project}}}");

    cy.byTestId('edit-webhook-field-proxy').find('input').focus().type('https://proxy.com');

    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_NAME_ID).should('not.exist');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID).should('not.exist');
    cy.byTestId('ktb-webhook-settings-add-header-button').click();
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_NAME_ID).should('exist');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID).should('exist');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_NAME_ID).find('input').focus().type('x-token');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID).find('input').focus().type('Bearer: ');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_SECRET_SELECTOR_HEADER).find('button').click();
    selectFirstItemOfVariableSelector();
    cy.byTestId(uniformPage.EDIT_WEBHOOK_EVENT_SELECTOR_HEADER).find('button').click();
    selectFirstItemOfVariableSelector();
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID)
      .find('input')
      .should('have.value', 'Bearer: {{.secret.SecretA.key1}}{{.data.project}}');

    cy.byTestId(uniformPage.UPDATE_SUBSCRIPTION_BUTTON_ID).click();
    // It should redirect after successfully sending the subscription
    cy.location('pathname').should(
      'eq',
      '/project/sockshop/uniform/services/0f2d35875bbaa72b972157260a7bd4af4f2826df/subscriptions/add'
    );
  });

  it('should delete a subscription', () => {
    cy.byTestId(uniformPage.UNIFORM_INTEGRATION_TABLE_LOC).find('dt-row').eq(1).click();
    cy.get(uniformPage.SUBSCRIPTION_EXPANDABLE_LOC)
      .first()
      .find('dt-expandable-panel')
      .byTestId('subscriptionDeleteButton')
      .click();

    // Check if confirmation dialog pops up
    cy.byTestId('dialogWarningMessage').should(
      'have.text',
      'Deleting this subscription will affect all projects. Please be certain.'
    );
    cy.get('dt-confirmation-dialog-actions').should('exist');

    // Check if it was removed from the list
    cy.get('dt-confirmation-dialog-actions button').first().click();
    cy.get(uniformPage.SUBSCRIPTION_DETAILS_LOC).should('have.length', 0);
  });

  it('should edit a subscription', () => {
    // given
    cy.byTestId(uniformPage.UNIFORM_INTEGRATION_TABLE_LOC).find('dt-row').eq(1).click();
    cy.get(uniformPage.SUBSCRIPTION_EXPANDABLE_LOC)
      .first()
      .find('dt-expandable-panel')
      .byTestId('subscriptionEditButton')
      .click();

    cy.byTestId(uniformPage.EDIT_SUBSCRIPTION_FIELD_TASK_ID).find('dt-select').focus().type('eval');
    cy.byTestId(uniformPage.EDIT_SUBSCRIPTION_FIELD_SUFFIX_ID).find('dt-select').focus().type('fin');
    cy.byTestId('edit-subscription-field-filterStageService')
      .find('input')
      .focus()
      .type('St{enter}de{enter}Ser{enter}cart{enter}');
    cy.byTestId(uniformPage.UPDATE_SUBSCRIPTION_BUTTON_ID).click();

    // It should redirect to overview if edited successfully
    cy.location('pathname').should('eq', '/project/sockshop/uniform/services/355311a7bec3f35bf3abc2484ab09bcba8e2b297');
  });

  it('should show an error message if can not parse shipyard.yaml', () => {
    // given
    cy.intercept('/api/project/sockshop/tasks', {
      statusCode: 500,
      body: 'Could not parse shipyard.yaml',
    }).as('tasksResult');
    cy.byTestId(uniformPage.UNIFORM_INTEGRATION_TABLE_LOC).find('dt-row').eq(1).click();
    cy.get(uniformPage.SUBSCRIPTION_EXPANDABLE_LOC)
      .first()
      .find('dt-expandable-panel')
      .byTestId('addSubscriptionButton')
      .click();

    cy.wait('@tasksResult');
    cy.wait('@tasksResult');
    cy.wait('@tasksResult');

    // It should show an error message and reload button
    cy.byTestId('keptn-notification-bar-message').should('have.text', 'Could not parse shipyard.yaml');
    cy.byTestId('ktb-modify-subscription-reload-button').should('exist');

    // eslint-disable-next-line promise/always-return,promise/catch-or-return
    cy.window().then((window) => {
      window.errorCount = 1;
    });
  });

  it('should reload page correctly if shipyard.yaml was not parsed initially', () => {
    // given
    cy.intercept('/api/project/sockshop/tasks', {
      statusCode: 500,
      body: 'Could not parse shipyard.yaml',
    }).as('tasksResult');
    cy.byTestId(uniformPage.UNIFORM_INTEGRATION_TABLE_LOC).find('dt-row').eq(1).click();
    cy.get(uniformPage.SUBSCRIPTION_EXPANDABLE_LOC)
      .first()
      .find('dt-expandable-panel')
      .byTestId('addSubscriptionButton')
      .click();

    cy.wait('@tasksResult');
    cy.wait('@tasksResult');
    cy.wait('@tasksResult');

    cy.byTestId('ktb-modify-subscription-reload-button').click();
    cy.intercept('/api/project/sockshop/tasks', { fixture: 'tasks.mock' });

    // It should not show an error message and reload button
    cy.byTestId('keptn-notification-bar-message').should('not.exist');
    cy.byTestId('ktb-modify-subscription-reload-button').should('not.exist');
    cy.get('h2').first().should('have.text', 'Create subscription');
  });

  function selectFirstItemOfVariableSelector(): void {
    cy.get('ktb-tree-list-select dt-tree-table-toggle-cell').first().find('button').click();
    cy.get('ktb-tree-list-select dt-tree-table-toggle-cell').eq(1).click();
  }
});
