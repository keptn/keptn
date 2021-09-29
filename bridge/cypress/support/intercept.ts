export function interceptIntegrations(): void {
  cy.intercept('/api/v1/metadata', {fixture: 'metadata.mock'});
  cy.intercept('/api/bridgeInfo', {fixture: 'bridgeInfo.mock'});
  cy.intercept('/api/project/sockshop?approval=true&remediation=true', {fixture: 'project.mock'});
  cy.intercept('/api/hasUnreadUniformRegistrationLogs', {body: false});
  cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {fixture: 'projects.mock'});
  cy.intercept('/api/uniform/registration', {fixture: 'registration.mock'});
  cy.intercept('/api/controlPlane/v1/log?integrationId=355311a7bec3f35bf3abc2484ab09bcba8e2b297&pageSize=100', {
    body: {
      logs: [],
    },
  });
  cy.intercept('/api/controlPlane/v1/log?integrationId=0f2d35875bbaa72b972157260a7bd4af4f2826df&pageSize=100', {
    body: {
      logs: [],
    },
  });
  cy.intercept('/api/uniform/registration/355311a7bec3f35bf3abc2484ab09bcba8e2b297/info', {
    body: {
      isControlPlane: true,
      isWebhookService: false,
    },
  });
  cy.intercept('/api/uniform/registration/0f2d35875bbaa72b972157260a7bd4af4f2826df/info', {
    body: {
      isControlPlane: true,
      isWebhookService: true,
    },
  });
  cy.intercept('/api/controlPlane/v1/uniform/registration/355311a7bec3f35bf3abc2484ab09bcba8e2b297/subscription', {
    body: {
      id: '0b77c90e-282d-4a7e-a96d-e23027265868',
    },
  });
  cy.intercept('/api/controlPlane/v1/uniform/registration/0f2d35875bbaa72b972157260a7bd4af4f2826df/subscription', {
    body: {
      id: 'b5111b1c-446a-410d-bb6c-e1dcd409c890',
    },
  });
  cy.intercept('/api/uniform/registration/webhook-service/config', {body: true});
  cy.intercept('/api/project/sockshop/tasks', {fixture: 'tasks.mock'});
  cy.intercept('/api/secrets/scope/keptn-webhook-service', {fixture: 'secrets.mock'});
}
