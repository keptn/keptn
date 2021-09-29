import { interceptIntegrations } from '../support/intercept';

describe('Integrations', () => {
  beforeEach(() => {
    interceptIntegrations();
    cy.visit('/project/sockshop/uniform/services');
  });

  it('should be on page uniform', () => {
    const buttons = cy.byTestId('uniform-submenu').find('button');
    buttons.first().should('have.class', 'active');
  });

  it('should select a integration and show related subscriptions', () => {
    // given, when
    cy.byTestId('keptn-uniform-integrations-table').get('dt-row').first().click();

    // then
    cy.get('ktb-expandable-tile h3').first().should('have.text', 'Subscriptions');
    cy.get('ktb-expandable-tile h3').last().should('have.text', 'Error events');
  });

  it('should add a simple subscription', () => {
    // given, when
    cy.byTestId('keptn-uniform-integrations-table').find('dt-row').eq(1).click();
    cy.get('ktb-expandable-tile dt-expandable-panel ktb-subscription-item').should('have.length', 1);
    cy.get('ktb-expandable-tile').first().find('dt-expandable-panel').byTestId('addSubscriptionButton').click();

    // Fill in all form fields
    cy.byTestId('edit-subscription-field-isGlobal').click().find('dt-checkbox').should('have.class', 'dt-checkbox-checked');
    cy.byTestId('edit-subscription-field-task').find('dt-select').focus().type('dep');
    cy.byTestId('edit-subscription-field-suffix').find('dt-select').focus().type('fin');
    cy.byTestId('edit-subscription-field-filterStageService').find('input').focus().type('St{enter}de{enter}Ser{enter}cart{enter}');
    cy.byTestId('updateSubscriptionButton').click();

    // then
    // It should redirect after successfully sending the subscription
    cy.location('pathname').should('eq', '/project/sockshop/uniform/services/355311a7bec3f35bf3abc2484ab09bcba8e2b297');
  });

  it('should add a webhook subscription', () => {
    cy.byTestId('keptn-uniform-integrations-table').find('dt-row').eq(0).click();
    cy.get('ktb-expandable-tile').first().find('dt-expandable-panel').byTestId('addSubscriptionButton').click();

    // then
    cy.get('h2').eq(1).should('have.text', 'Webhook configuration');
    cy.byTestId('edit-subscription-field-isGlobal').should('not.exist');

    cy.byTestId('edit-subscription-field-task').find('dt-select').focus().type('dep');
    cy.byTestId('edit-subscription-field-suffix').find('dt-select').focus().type('fin');
    cy.byTestId('edit-webhook-field-method').find('dt-select').focus().type('GET');

    // URL: insert text and add secret
    cy.byTestId('edit-webhook-field-url').find('input').focus().type('https://example.com?secret=');
    cy.byTestId('edit-webhook-field-url').find('.dt-form-field-suffix button').click();
    addSecret();
    cy.byTestId('edit-webhook-field-url').find('input').should('have.value', 'https://example.com?secret={{.SecretA.key1}}');

    // Payload: insert text and add secret
    cy.byTestId('edit-webhook-field-payload').find('textarea').focus().type('{id: \'123456789\', secret: ');
    cy.byTestId('edit-webhook-field-payload').find('.dt-form-field-suffix button').click();
    addSecret();
    cy.byTestId('edit-webhook-field-payload').find('textarea').focus().type('}');
    cy.byTestId('edit-webhook-field-payload').find('textarea').should('have.value', '{id: \'123456789\', secret: {{.SecretA.key1}}}');

    cy.byTestId('edit-webhook-field-proxy').find('input').focus().type('https://proxy.com');

    cy.byTestId('edit-webhook-field-headerName').should('not.exist');
    cy.byTestId('edit-webhook-field-headerValue').should('not.exist');
    cy.byTestId('ktb-webhook-settings-add-header-button').click();
    cy.byTestId('edit-webhook-field-headerName').should('exist');
    cy.byTestId('edit-webhook-field-headerValue').should('exist');
    cy.byTestId('edit-webhook-field-headerName').find('input').focus().type('x-token');
    cy.byTestId('edit-webhook-field-headerValue').find('input').focus().type('Bearer: ');
    cy.byTestId('edit-webhook-field-headerValue').find('.dt-form-field-suffix button').click();
    addSecret();
    cy.byTestId('edit-webhook-field-headerValue').find('input').should('have.value', 'Bearer: {{.SecretA.key1}}');

    cy.byTestId('updateSubscriptionButton').click();
    // It should redirect after successfully sending the subscription
    cy.location('pathname').should('eq', '/project/sockshop/uniform/services/0f2d35875bbaa72b972157260a7bd4af4f2826df');
  });


  function addSecret(): void {
    cy.get('ktb-tree-list-select dt-tree-table-toggle-cell').first().find('button').click();
    cy.get('ktb-tree-list-select dt-tree-table-toggle-cell').eq(1).click();
  }
});
