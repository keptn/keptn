import BasePage from '../support/pageobjects/BasePage';
import SecretsPage from '../support/pageobjects/SecretsPage';

describe('Keptn Secrets adding deleting test', () => {
  it('The test goes over the pages and does validation', () => {
    const basePage = new BasePage();
    const secretsPage = new SecretsPage();
    const SECRET_NAME = 'dynatrace-prod';
    const SECRET_KEY = 'DT_API_TOKEN';
    const SECRET_VALUE = 'secretvalue!@#$%^&*(!@#$%^&*()';
    const SECRET_SCOPE = 'dynatrace-service';
    const DYNATRACE_PROJECT = 'dynatrace';

    cy.fixture('get.project.json').as('initProjectJSON');
    cy.fixture('metadata.json').as('initmetadata');

    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
    cy.intercept('GET', 'api/v1/metadata', { fixture: 'metadata.json' }).as('metadataCmpl');
    cy.intercept('GET', 'api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', {
      fixture: 'get.project.json',
    }).as('initProjects');
    cy.intercept('GET', 'api/controlPlane/v1/sequence/dynatrace?pageSize=5', { fixture: 'project.sequences.json' });

    cy.intercept('POST', 'api/secrets/v1/secret', {
      statusCode: 200,
    }).as('postSecrets');

    cy.intercept('GET', 'api/secrets/v1/secret', {
      statusCode: 200,
      body: {
        Secrets: [
          { name: 'dynatrace', scope: 'dynatrace-service', keys: ['DT_API_TOKEN', 'DT_TENANT'] },
          { name: 'dynatrace-prod', scope: 'dynatrace-service', keys: [SECRET_KEY] },
          { name: 'api', scope: 'keptn-default', keys: ['API_TOKEN'] },
          { name: 'webhook', scope: 'keptn-webhook-service', keys: ['webhook_url', 'webhook_secret', 'webhook_proxy'] },
        ],
      },
    }).as('getSecrets');

    cy.intercept('GET', 'api/project/dynatrace?approval=true&remediation=true', {
      statusCode: 200,
    }).as('getApproval');

    cy.intercept('POST', 'api/hasUnreadUniformRegistrationLogs', {
      statusCode: 200,
    }).as('hasUnreadUniformRegistrationLogs');

    cy.intercept('POST', 'api/uniform/registration', {
      statusCode: 200,
      body: '[]',
    }).as('uniformRegPost');

    cy.intercept('DELETE', 'api/secrets/v1/secret?name=dynatrace-prod&scope=dynatrace-service', {
      statusCode: 200,
    }).as('deleteSecret');

    cy.visit('/');
    cy.wait('@metadataCmpl');
    basePage.selectProject(DYNATRACE_PROJECT);

    basePage.goToUniformPage().goToSecretsPage();

    secretsPage.addSecret(SECRET_NAME, SECRET_SCOPE, SECRET_KEY, SECRET_VALUE);

    cy.get('dt-row').eq(1).find('dt-cell').eq(0).find('p').should('have.text', SECRET_NAME);
    cy.get('dt-row').eq(1).find('dt-cell').eq(1).find('p').should('have.text', SECRET_SCOPE);
    cy.get('dt-row').eq(1).find('dt-cell').eq(2).find('p').should('contain.text', SECRET_KEY);

    secretsPage.deleteSecret(SECRET_NAME);
    cy.get('dt-row').eq(1).find('dt-cell').eq(0).find('p').should('not.have.text', SECRET_NAME);
  });
});
