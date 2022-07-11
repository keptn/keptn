import ServicesPage from '../support/pageobjects/ServicesPage';

describe('Test service remediations', () => {
  const servicesPage = new ServicesPage();
  const projectName = 'sockshop';
  const serviceName = 'carts';

  beforeEach(() => {
    servicesPage.interceptAll();
  });

  it('should show a loading indicator next to a stage if it is running', () => {
    servicesPage
      .interceptRunning()
      .visitServicePage(projectName)
      .selectService(serviceName, 'v0.1.2')
      .assertIsStageLoading('production', true)
      .assertIsStageLoading('staging', false);
  });

  it('should show remediations and loading indicator', () => {
    servicesPage
      .interceptRemediations()
      .visitServicePage(projectName)
      .selectService(serviceName, 'v0.1.1')
      .selectStage('staging')
      .assertRemediationSequenceCount(1)
      .assertIsStageLoading('staging', true)
      .assertIsStageLoading('production', false);
  });

  it('should not show remediations and loading indicator', () => {
    servicesPage
      .visitServicePage(projectName)
      .selectService(serviceName, 'v0.1.2')
      .selectStage('staging')
      .assertRemediationSequenceCount(0)
      .assertIsStageLoading('staging', false)
      .assertIsStageLoading('production', false);
  });
});
