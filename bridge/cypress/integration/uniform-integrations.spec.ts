import { interceptIntegrations } from '../support/intercept';
import UniformPage from '../support/pageobjects/UniformPage';

const uniformPage = new UniformPage();
const webhookID = '0f2d35875bbaa72b972157260a7bd4af4f2826df';
const integrationID = '355311a7bec3f35bf3abc2484ab09bcba8e2b297'; // not webhook, is in control plane

describe('Integrations default requests', () => {
  beforeEach(() => {
    interceptIntegrations();
    uniformPage.visit('sockshop');
  });

  it('should be on page uniform', () => {
    const buttons = cy.byTestId(uniformPage.UNIFORM_SUBMENU_LOC);
    buttons.first().should('have.class', 'active');
  });

  it('should show 8 registrations', () => {
    cy.byTestId(uniformPage.UNIFORM_INTEGRATION_TABLE_LOC).find('dt-row').should('have.length', 8);
  });

  it('should show error events list', () => {
    uniformPage.selectIntegration('webhook-service');
    cy.get('ktb-uniform-registration-logs').should('exist');
    cy.contains('h3', 'webhook-service');
  });

  it('should not show logs', () => {
    uniformPage.selectIntegration('approval-service');
    cy.get('.uniform-registration-error-log').should('not.exist');
    cy.get('ktb-uniform-registration-logs').contains('No events for this integration available');
  });

  it('should have 1 log', () => {
    uniformPage.selectIntegration('webhook-service');
    cy.get('.uniform-registration-error-log>div').should('have.length', 1);
  });

  it('should show first 2 rows as unread', () => {
    cy.intercept('/api/hasUnreadUniformRegistrationLogs', { body: true });

    uniformPage
      .assertHasIntegrationErrorIndicator('webhook-service', true)
      .selectIntegration('webhook-service')
      .assertIndicatorsShowing(0)
      .assertIndicatorsTextShowing(1) // 1 unread error log
      .assertHasIntegrationErrorIndicator('webhook-service', false)
      .selectIntegration('jmeter-service');

    cy.intercept('/api/controlPlane/v1/log?integrationId=0f2d35875bbaa72b972157260a7bd4af4f2826df&pageSize=100', {
      body: {
        logs: [
          {
            integrationid: '0f2d35875bbaa72b972157260a7bd4af4f2826df',
            message: 'my error2',
            shkeptncontext: ' 7394b5b3-2fb3-4cb7-b435-d0e9d6f0cb87',
            task: 'my task2',
            time: new Date(Date.now() + 60_000),
            triggeredid: 'bd3bc477-6d0f-4d71-b15d-c33e953a74ba',
          },
          {
            integrationid: '0f2d35875bbaa72b972157260a7bd4af4f2826df',
            message: 'my error3',
            shkeptncontext: ' 7394b5b3-2fb3-4cb7-b435-d0e9d6f0cb87',
            task: 'my task3',
            time: new Date(Date.now() + 60_000),
            triggeredid: 'bd3bc477-6d0f-4d71-b15d-c33e953a74ba',
          },
        ],
      },
    });
    uniformPage
      .selectIntegration('webhook-service')
      .assertIndicatorsTextShowing(2) // 1 unread error log
      .assertHasIntegrationErrorIndicator('webhook-service', false);
  });

  it('should select an integration and show related subscriptions', () => {
    // given, when
    uniformPage.selectIntegration('webhook-service');

    // then
    cy.get(uniformPage.SUBSCRIPTION_EXP_HEADER_LOC).first().should('have.text', 'Subscriptions');
    cy.get(uniformPage.SUBSCRIPTION_EXP_HEADER_LOC).last().should('have.text', 'Error events');
  });

  it('should add a simple subscription', () => {
    // given, when
    uniformPage.selectIntegration('jmeter-service');
    cy.get(uniformPage.SUBSCRIPTION_DETAILS_LOC).should('have.length', 1);
    uniformPage
      .addSubscription()
      .switchIsGlobalState()
      .isCreateButton()
      .assertIsGlobalChecked(true)
      .setTaskPrefix('deployment')
      .setTaskSuffix('finished')
      .appendStages('dev')
      .appendServices('carts')
      .update();

    // then
    // It should redirect after successfully sending the subscription
    cy.location('pathname').should('eq', `/project/sockshop/settings/uniform/integrations/${integrationID}`);
  });

  it('should add a webhook subscription', () => {
    uniformPage.selectIntegration('webhook-service').addSubscription();

    // then
    cy.get('h2').eq(1).should('have.text', 'Webhook configuration');
    cy.byTestId(uniformPage.EDIT_SUBSCRIPTION_FIELD_GLOBAL_ID).should('not.exist');

    uniformPage
      .setTaskPrefix('deployment')
      .setTaskSuffix('finished')
      .setWebhookMethod('GET')
      .appendURL('https://example.com?secret=')
      .openSecretSelectorURL()
      .selectFirstItemOfVariableSelector()
      .openEventSelectorURL()
      .selectFirstItemOfVariableSelector()
      .appendPayload(`{id: '123456789', secret: `)
      .openSecretSelectorPayload()
      .selectFirstItemOfVariableSelector()
      .openEventSelectorPayload()
      .selectFirstItemOfVariableSelector()
      .appendPayload('}')
      .appendProxy('https://proxy.com')
      .assertURL('https://example.com?secret={{.secret.SecretA.key1}}{{.data.project}}')
      .assertPayload(`{id: '123456789', secret: {{.secret.SecretA.key1}}{{.data.project}}}`);

    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_NAME_ID).should('not.exist');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID).should('not.exist');
    uniformPage.addHeader();
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_NAME_ID).should('exist');
    cy.byTestId(uniformPage.EDIT_WEBHOOK_FIELD_HEADER_VALUE_ID).should('exist');

    uniformPage
      .appendHeaderName('x-token')
      .appendHeaderValue('Bearer: ')
      .openSecretSelectorHeader()
      .selectFirstItemOfVariableSelector()
      .openEventSelectorHeader()
      .selectFirstItemOfVariableSelector()
      .assertHeaderValue(0, 'Bearer: {{.secret.SecretA.key1}}{{.data.project}}')
      .update();

    // It should redirect after successfully sending the subscription
    cy.location('pathname').should('eq', `/project/sockshop/settings/uniform/integrations/${webhookID}`);
  });

  it('should delete a subscription', () => {
    uniformPage.selectIntegration('jmeter-service').deleteSubscription(0);

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
    uniformPage
      .selectIntegration('jmeter-service')
      .editSubscription(0)
      .isUpdateButton()
      .setTaskPrefix('evaluation')
      .setTaskSuffix('finished')
      .appendStages('dev')
      .appendServices('cart')
      .update();

    // It should redirect to overview if edited successfully
    cy.location('pathname').should('eq', `/project/sockshop/settings/uniform/integrations/${integrationID}`);
  });
});

describe('Integrations changed requests', () => {
  beforeEach(() => {
    interceptIntegrations();
  });

  it('should show error event indicators', () => {
    cy.intercept('/api/hasUnreadUniformRegistrationLogs', { body: true });
    uniformPage
      .visit('sockshop')
      .assertIntegrationErrorCount('webhook-service', 10)
      .assertIndicatorsShowing(2)
      .assertIndicatorsTextShowing(1)
      .selectIntegration('webhook-service')
      .assertErrorEventsShowing(1);
  });

  it('should remove error event indicator on selection', () => {
    cy.intercept('/api/hasUnreadUniformRegistrationLogs', { body: true });
    uniformPage
      .visit('sockshop')
      .assertHasIntegrationErrorIndicator('webhook-service', true)
      .selectIntegration('webhook-service')
      .assertIndicatorsShowing(0)
      .assertIndicatorsTextShowing(1) // 1 unread error log
      .assertHasIntegrationErrorIndicator('webhook-service', false)
      .selectIntegration('jmeter-service')
      .assertIndicatorsShowing(0)
      .assertIndicatorsTextShowing(0)
      .selectIntegration('webhook-service')
      .assertIndicatorsShowing(0)
      .assertIndicatorsTextShowing(0); // now error log is read
  });

  it('should not remove error event indicator if integration without logs is selected', () => {
    cy.intercept('/api/hasUnreadUniformRegistrationLogs', { body: true });
    uniformPage
      .visit('sockshop')
      .selectIntegration('jmeter-service')
      .assertHasIntegrationErrorIndicator('webhook-service', true)
      .assertIndicatorsShowing(2)
      .assertIndicatorsTextShowing(1);
  });

  it('should have disabled buttons for a subscription, when subscription id is not given', () => {
    // given, when
    cy.intercept('/api/uniform/registration', { fixture: 'registration-old-format.mock' });
    uniformPage.visit('sockshop').selectIntegration('jmeter-service');
    const editButton = cy.get(uniformPage.SUBSCRIPTION_EXPANDABLE_LOC).first().byTestId('subscriptionEditButton');
    const deleteButton = cy.get(uniformPage.SUBSCRIPTION_EXPANDABLE_LOC).first().byTestId('subscriptionDeleteButton');

    // then
    editButton.should('be.disabled');
    deleteButton.should('be.disabled');
  });
});

describe('Add webhook subscriptions', () => {
  beforeEach(() => {
    interceptIntegrations();
    uniformPage.visitAdd(webhookID);
  });
  it('should have disabled button if first and second control is invalid', () => {
    uniformPage.assertIsUpdateButtonEnabled(false).assertIsUpdateButtonEnabled(false);
  });

  it('should have disabled button if first control is valid and second control is invalid', () => {
    uniformPage.setTaskPrefix('deployment').assertIsUpdateButtonEnabled(false);
  });

  it('should have disabled button if first control is invalid and second control is valid', () => {
    uniformPage.setTaskSuffix('triggered').assertIsUpdateButtonEnabled(false);
  });

  it('should have disabled button if filter contains service but not a stage', () => {
    uniformPage
      .setTaskPrefix('deployment')
      .appendServices('carts')
      .setWebhookMethod('GET')
      .appendURL('https://example.com')
      .assertIsUpdateButtonEnabled(false);
  });

  it('should have a disabled button if the webhook form is empty', () => {
    uniformPage.setTaskPrefix('deployment').setTaskSuffix('triggered').assertIsUpdateButtonEnabled(false);
  });

  it('should have an enabled button if the webhook form is valid', () => {
    uniformPage
      .setTaskPrefix('deployment')
      .setTaskSuffix('triggered')
      .setWebhookMethod('GET')
      .appendURL('https://example.com')
      .assertIsUpdateButtonEnabled(true);
  });

  it('should show all suffixes', () => {
    uniformPage.shouldHaveTaskSuffixes(['*', 'triggered', 'started', 'finished']);
  });

  it('should show webhook form', () => {
    cy.get('ktb-webhook-settings').should('exist');
  });

  it('should not show project checkbox', () => {
    cy.byTestId(uniformPage.EDIT_SUBSCRIPTION_FIELD_GLOBAL_ID).should('not.exist');
  });

  it('should have sendFinished and sendStarted checkbox disabled', () => {
    uniformPage
      .setTaskPrefix('deployment')
      .setTaskSuffix('finished')
      .assertIsSendStartedButtonsEnabled(false)
      .assertIsSendFinishedButtonsEnabled(false);
  });

  it('should have sendFinished and sendStarted checkbox enabled and true by default', () => {
    uniformPage
      .setTaskPrefix('deployment')
      .setTaskSuffix('triggered')
      .assertIsSendStartedButtonsEnabled(true)
      .assertIsSendFinishedButtonsEnabled(true)
      .assertIsSendStarted(true)
      .assertIsSendFinished(true);
  });
});

describe('Add control plane subscription default requests', () => {
  beforeEach(() => {
    interceptIntegrations();
    uniformPage.visitAdd(integrationID);
  });

  it('should have enabled button if task is valid', () => {
    uniformPage.setTaskPrefix('deployment').setTaskSuffix('triggered').assertIsUpdateButtonEnabled(true);
  });

  it('should have enabled button if filter contains a stage and a service', () => {
    uniformPage
      .setTaskPrefix('deployment')
      .setTaskSuffix('triggered')
      .appendServices('carts')
      .appendStages('dev')
      .assertIsUpdateButtonEnabled(true);
  });

  it('should have enabled button if filter contains just a stage', () => {
    uniformPage
      .setTaskPrefix('deployment')
      .setTaskSuffix('triggered')
      .appendStages('dev')
      .assertIsUpdateButtonEnabled(true);
  });

  it('it should enable "use for all projects" checkbox if filter is cleared', () => {
    uniformPage.appendStages('dev').isGlobalEnabled(false).clearFilter().isGlobalEnabled(true);
  });

  it('it should disable "use for all projects" checkbox and set to false if filter is set', () => {
    uniformPage.appendStages('dev').isGlobalEnabled(false);
  });

  it('should show project checkbox', () => {
    cy.byTestId(uniformPage.EDIT_SUBSCRIPTION_FIELD_GLOBAL_ID).should('exist');
  });

  it('should not show webhook form', () => {
    cy.get('ktb-webhook-settings').should('not.exist');
  });
});

describe('Add control plane subscription dynamic request', () => {
  beforeEach(() => {
    interceptIntegrations();
  });

  it('should have disabled button if updating', () => {
    cy.intercept(`/api/uniform/registration/${integrationID}/subscription`, {
      body: {
        id: '0b77c90e-282d-4a7e-a96d-e23027265868',
      },
      delay: 5000,
    });
    uniformPage
      .visitAdd(integrationID)
      .setTaskPrefix('deployment')
      .setTaskSuffix('triggered')
      .update()
      .assertIsUpdateButtonEnabled(false);
  });

  xit('should show an error message if can not parse shipyard.yaml', () => {
    // given
    cy.intercept('/api/project/sockshop/tasks', {
      statusCode: 500,
      body: 'Could not parse shipyard.yaml',
    }).as('tasksResult');
    uniformPage.visitAdd(integrationID).selectIntegration('jmeter-service').addSubscription();

    cy.wait('@tasksResult');
    cy.wait('@tasksResult');
    cy.wait('@tasksResult');

    // It should show an error message and reload button
    cy.byTestId('keptn-notification-bar-message').should('have.text', 'Could not parse shipyard.yaml');
    cy.byTestId('ktb-modify-subscription-reload-button').should('exist');
  });

  xit('should reload page correctly if shipyard.yaml was not parsed initially', () => {
    // given
    cy.intercept('/api/project/sockshop/tasks', {
      statusCode: 500,
      body: 'Could not parse shipyard.yaml',
    }).as('tasksResult');
    uniformPage.visitAdd(integrationID).selectIntegration('jmeter-service').addSubscription();

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
});

describe('Add execution plane subscription', () => {
  const executionPlaneIntegrationID = 'myIntegrationID';

  beforeEach(() => {
    interceptIntegrations();
  });

  it('should only have triggered suffix', () => {
    cy.intercept(`/api/uniform/registration/${executionPlaneIntegrationID}/info`, {
      body: {
        isControlPlane: false,
        isWebhookService: false,
      },
    });
    uniformPage.visitAdd(executionPlaneIntegrationID).shouldHaveTaskSuffixes(['triggered']);
  });
});

describe('Edit subscriptions', () => {
  const subscriptionID = 'mySubscriptionID';
  beforeEach(() => {
    interceptIntegrations();
  });

  it('should set the right properties and enable the button when a global subscription is set', () => {
    cy.intercept(`/api/controlPlane/v1/uniform/registration/${integrationID}/subscription/${subscriptionID}`, {
      body: {
        event: 'sh.keptn.event.deployment.triggered',
        filter: {
          projects: [],
          services: [],
          stages: [],
          id: subscriptionID,
        },
      },
    });
    uniformPage
      .visitEdit(integrationID, subscriptionID)
      .assertIsGlobalChecked(true)
      .taskPrefixEquals('deployment')
      .taskSuffixEquals('triggered')
      .filterTagsLengthEquals(0)
      .assertIsUpdateButtonEnabled(true);
  });

  it('should set the right properties and enable the button when a subscription is set', () => {
    const service = 'carts';
    const stage = 'dev';
    const projectName = 'sockshop';

    cy.intercept(`/api/controlPlane/v1/uniform/registration/${integrationID}/subscription/${subscriptionID}`, {
      body: {
        event: 'sh.keptn.event.test.finished',
        filter: {
          projects: [projectName],
          services: [service],
          stages: [stage],
          id: subscriptionID,
        },
      },
    });

    uniformPage
      .visitEdit(integrationID, subscriptionID)
      .assertIsGlobalChecked(false)
      .taskPrefixEquals('test')
      .taskSuffixEquals('finished')
      .filterTagsLengthEquals(2)
      .shouldHaveStages([stage])
      .shouldHaveServices([service])
      .assertIsUpdateButtonEnabled(true);
  });
});
