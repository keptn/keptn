import { interceptEmptyEnvironmentScreen } from '../support/intercept';
import EnvironmentPage from '../support/pageobjects/EnvironmentPage';
import { ServicesSettingsPage } from '../support/pageobjects/ServicesSettingsPage';

describe('Environment Screen empty', () => {
  const environmentPage = new EnvironmentPage();
  beforeEach(() => {
    interceptEmptyEnvironmentScreen();
    environmentPage.visit('dynatrace');
  });

  it('should redirect to create service and redirect back after creation', () => {
    const serviceSettings = new ServicesSettingsPage();

    environmentPage.clickCreateService('dev');
    cy.location('pathname').should('eq', '/project/dynatrace/settings/services/create');
    serviceSettings.createService('my-new-service');
    cy.location('pathname').should('eq', '/project/dynatrace');
  });
});
