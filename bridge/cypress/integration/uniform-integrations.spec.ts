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

    // URL: insert text and add secret
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_URL_ID).find('input').focus().type('https://example.com?secret=');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_URL_ID).find('.dt-form-field-suffix button').click();
    addSecret();
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_URL_ID)
      .find('input')
      .should('have.value', 'https://example.com?secret={{.secret.SecretA.key1}}');

    // Payload: insert text and add secret
    cy.byTestId(uniformPage.EDIT_WEBHOOK_PAYLOAD_ID).find('textarea').focus().type("{id: '123456789', secret: ");
    cy.byTestId(uniformPage.EDIT_WEBHOOK_PAYLOAD_ID).find('.dt-form-field-suffix button').click();
    addSecret();
    cy.byTestId(uniformPage.EDIT_WEBHOOK_PAYLOAD_ID).find('textarea').focus().type('}');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_PAYLOAD_ID)
      .find('textarea')
      .should('have.value', "{id: '123456789', secret: {{.secret.SecretA.key1}}}");

    cy.byTestId('edit-webhook-field-proxy').find('input').focus().type('https://proxy.com');

    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_NAME_ID).should('not.exist');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID).should('not.exist');
    cy.byTestId('ktb-webhook-settings-add-header-button').click();
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_NAME_ID).should('exist');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID).should('exist');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_NAME_ID).find('input').focus().type('x-token');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID).find('input').focus().type('Bearer: ');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID).find('.dt-form-field-suffix button').click();
    addSecret();
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID)
      .find('input')
      .should('have.value', 'Bearer: {{.secret.SecretA.key1}}');

    cy.byTestId(uniformPage.UPDATE_SUBSCRIPTION_BUTTON_ID).click();
    // It should redirect after successfully sending the subscription
    cy.location('pathname').should(
      'eq',
      '/project/sockshop/uniform/services/0f2d35875bbaa72b972157260a7bd4af4f2826df/subscriptions/add'
    );
  });

  function addSecret(): void {
    cy.get('ktb-tree-list-select dt-tree-table-toggle-cell').first().find('button').click();
    cy.get('ktb-tree-list-select dt-tree-table-toggle-cell').eq(1).click();
  }
});
