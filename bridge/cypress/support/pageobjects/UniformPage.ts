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
}
export default UniformPage;
