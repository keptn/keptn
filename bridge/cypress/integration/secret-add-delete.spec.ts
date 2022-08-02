import { ProjectBoardPage } from '../support/pageobjects/ProjectBoardPage';
import SecretsPage from '../support/pageobjects/SecretsPage';
import DashboardPage from '../support/pageobjects/DashboardPage';
import EnvironmentPage from '../support/pageobjects/EnvironmentPage';

describe('Keptn Secrets adding deleting test', () => {
  const basePage = new ProjectBoardPage();
  const dashboardPage = new DashboardPage();
  const environmentPage = new EnvironmentPage();
  const secretsPage = new SecretsPage();
  const project = 'sockshop';

  beforeEach(() => {
    dashboardPage.intercept();
    environmentPage.intercept();
    secretsPage.intercept();
  });

  it('should navigate to add secret page', () => {
    dashboardPage.visit();
    basePage.selectProject(project);
    basePage.goToUniformPage().goToSecretsPage();
    secretsPage.clickAddSecret();
    cy.location('pathname').should('eq', `/project/${project}/settings/uniform/secrets/add`);
  });

  it('should add a secret', () => {
    secretsPage
      .visitCreate(project)
      .setSecret('dynatrace-prod', 'dynatrace-service', 'DT_API_TOKEN', 'secretvalue!@#$%^&*(!@#$%^&*()')
      .assertScopesEnabled(true)
      .createSecret();
  });

  it('should delete a secret', () => {
    const SECRET_NAME = 'dynatrace-prod';

    secretsPage.visit(project).deleteSecret(SECRET_NAME).secretExistsInList(SECRET_NAME, 1);
  });

  it('should have a specific secret in the list', () => {
    secretsPage.visit(project).assertSecretInList(1, 'dynatrace-prod', 'dynatrace-service', 'DT_API_TOKEN');
  });

  it('should have disabled "remove key-value pair" icon-button if there is only one key-value pair', () => {
    secretsPage.visitCreate(project).assertKeyValuePairLength(1).assertKeyValuePairEnabled(0, false);
  });

  it('should have enabled "remove key-value pair" icon-button if there is more than one key-value pair', () => {
    secretsPage
      .visitCreate(project)
      .addKeyValuePair()
      .assertKeyValuePairLength(2)
      .assertKeyValuePairEnabled(0, true)
      .assertKeyValuePairEnabled(1, true);
  });

  it('should have disabled scope dropdown and disabled create button if scopes are loading', () => {
    cy.intercept('GET', 'api/secrets/v1/scope', {
      statusCode: 200,
      body: {
        scopes: [],
      },
      delay: 10_000,
    });
    secretsPage
      .visitCreate(project)
      .appendSecretName('my-secret')
      .appendSecretKey(0, 'my-key')
      .appendSecretValue(0, 'my-value')
      .assertScopesEnabled(false)
      .assertCreateButtonEnabled(false);
  });
});
