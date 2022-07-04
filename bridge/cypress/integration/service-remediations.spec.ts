import ServicesPage from '../support/pageobjects/ServicesPage';

describe('Test service remediations', () => {
  const servicesPage = new ServicesPage();
  const projectName = 'sockshop';
  const serviceName = 'carts';

  beforeEach(() => {
    servicesPage.interceptAll();
  });

  it('should show remediations', () => {
    servicesPage
      .interceptRemediations()
      .visitServicePage(projectName)
      .selectService(serviceName, 'v0.1.1')
      .selectStage('staging')
      .assertRemediationSequenceCount(1);
  });

  it('should not show remediations', () => {
    servicesPage
      .visitServicePage(projectName)
      .selectService(serviceName, 'v0.1.2')
      .selectStage('staging')
      .assertRemediationSequenceCount(0);
  });
});
