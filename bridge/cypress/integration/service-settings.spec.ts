import { ServicesSettingsPage } from '../support/pageobjects/ServicesSettingsPage';

describe('Service settings', () => {
  const serviceSettingsPage = new ServicesSettingsPage();
  const projectName = 'sockshop';
  const serviceName = 'carts';

  it('should show a message when the file tree is empty', () => {
    serviceSettingsPage.intercept().visitService(projectName, serviceName).assertNoFilesMessageExists(true);
  });
});
