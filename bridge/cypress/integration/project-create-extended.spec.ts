import NewProjectCreatePage from '../support/pageobjects/NewProjectCreatePage';

describe('Create extended project https test', () => {
  const createProjectPage = new NewProjectCreatePage();

  beforeEach(() => {
    createProjectPage.intercept(true).visit();
  });

  it('should not show "Set Git upstream" button', () => {
    createProjectPage.assertUpdateButtonExists(false);
  });

  it('should have enabled button if form is valid', () => {
    createProjectPage.enterBasicValidProjectHttps().assertCreateButtonEnabled(true);
  });

  it('should be valid certificate and have enabled button', () => {
    createProjectPage.enterBasicValidProjectHttps().typeValidCertificate().assertCreateButtonEnabled(true);
  });

  it('should be invalid certificate and have disabled button and error', () => {
    createProjectPage.enterBasicValidProjectHttps().typeInvalidCertificate().assertCreateButtonEnabled(false);
  });

  it('should be valid certificate file and have enabled button', () => {
    createProjectPage.enterBasicValidProjectHttps().setValidCertificateFile().assertCreateButtonEnabled(true);
  });

  it('should be invalid certificate file and have disabled button and error', () => {
    createProjectPage.enterBasicValidProjectHttps().setInvalidCertificateFile().assertCreateButtonEnabled(false);
  });

  it('should show proxy form if proxy is enabled', () => {
    createProjectPage.setEnableProxy(true).assertProxyFormVisible(true);
  });

  it('should not show proxy form if proxy is disabled', () => {
    createProjectPage
      .assertProxyFormVisible(false)
      .setEnableProxy(true)
      .assertProxyFormVisible(true)
      .setEnableProxy(false)
      .assertProxyFormVisible(false);
  });

  it('should only be able to type numbers for the port', () => {
    createProjectPage
      .setEnableProxy(true)
      .typeProxyPort('0.1')
      .assertProxyPort('01')
      .clearProxyPort()

      .typeProxyPort('0,1')
      .assertProxyPort('01')
      .clearProxyPort()

      .typeProxyPort(-5000)
      .assertProxyPort(5000)
      .clearProxyPort()

      .typeProxyPort('abc')
      .assertProxyPort('')
      .clearProxyPort();

    for (let i = 0; i < 10; ++i) {
      createProjectPage.typeProxyPort(i).assertProxyPort(i).clearProxyPort();
    }
  });

  it('should have disabled create button if proxy form is invalid', () => {
    createProjectPage
      .enterBasicValidProjectHttps()
      .setEnableProxy(true)
      .assertCreateButtonEnabled(false)
      .typeProxyUrl('0.0.0.0')
      .assertCreateButtonEnabled(false)
      .typeProxyPort(5000)
      .assertCreateButtonEnabled(true)
      .clearProxyUrl()
      .assertCreateButtonEnabled(false);
  });

  it('should have enabled create button if invalid proxy form is disabled', () => {
    createProjectPage.enterBasicValidProjectHttps().setEnableProxy(true).assertCreateButtonEnabled(false);
    cy.wait(200);
    createProjectPage.setEnableProxy(false).assertCreateButtonEnabled(true);
  });

  it('should not change validity to false of proxy form if username or password is entered', () => {
    createProjectPage
      .enterBasicValidProjectHttps()
      .setEnableProxy(true)
      .typeProxyUrl('0.0.0.0')
      .typeProxyPort('5000')
      .assertCreateButtonEnabled(true)

      .typeProxyUsername('myUser')
      .assertCreateButtonEnabled(true)

      .typeProxyPassword('myPassword')
      .assertCreateButtonEnabled(true)

      .clearProxyPassword()
      .assertCreateButtonEnabled(true)

      .typeProxyPassword('myPassword')
      .clearProxyUsername()
      .assertCreateButtonEnabled(true);
  });

  it('should not change validity to true of proxy form if username or password is entered', () => {
    createProjectPage
      .enterBasicValidProjectHttps()
      .setEnableProxy(true)
      .typeProxyUrl('0.0.0.0')
      .assertCreateButtonEnabled(false)

      .typeProxyUsername('myUser')
      .assertCreateButtonEnabled(false)

      .typeProxyPassword('myPassword')
      .assertCreateButtonEnabled(false)

      .clearProxyPassword()
      .assertCreateButtonEnabled(false)

      .typeProxyPassword('myPassword')
      .clearProxyUsername()
      .assertCreateButtonEnabled(false);
  });

  it('should not delete proxy information if it is disabled and enabled again', () => {
    createProjectPage
      .enterFullValidProjectHttps()
      .setEnableProxy(false)
      .setEnableProxy(true)
      .validateFullValidProjectHttps();
  });
});

describe('Create extended project ssh test', () => {
  const createProjectPage = new NewProjectCreatePage();

  beforeEach(() => {
    createProjectPage.intercept(true).visit().selectSshForm();
  });

  it('should have enabled button if form is valid', () => {
    createProjectPage.enterBasicValidProjectSsh().assertCreateButtonEnabled(true);
  });

  it('should not change validity to true if git username or private key passphrase is entered', () => {
    createProjectPage
      .assertCreateButtonEnabled(false)

      .typeGitUsernameSsh('myUserName')
      .assertCreateButtonEnabled(false)

      .typeSshPrivateKeyPassphrase('myPassphrase')
      .assertCreateButtonEnabled(false)

      .clearSshPrivateKeyPassphrase()
      .assertCreateButtonEnabled(false)

      .typeSshPrivateKeyPassphrase('myPassphrase')
      .clearGitUsernameSsh()
      .assertCreateButtonEnabled(false);
  });

  it('should not change validity to false if git username or private key passphrase is entered', () => {
    createProjectPage
      .enterBasicValidProjectSsh()
      .assertCreateButtonEnabled(true)

      .typeGitUsernameSsh('myUserName')
      .assertCreateButtonEnabled(true)

      .typeSshPrivateKeyPassphrase('myPassphrase')
      .assertCreateButtonEnabled(true)

      .clearSshPrivateKeyPassphrase()
      .assertCreateButtonEnabled(true)

      .typeSshPrivateKeyPassphrase('myPassphrase')
      .clearGitUsernameSsh()
      .assertCreateButtonEnabled(true);
  });

  it('should be valid private key and have enabled button', () => {
    createProjectPage.enterBasicValidProjectSsh(false).typeValidSshPrivateKey().assertCreateButtonEnabled(true);
  });

  it('should be invalid private key and have disabled button and error', () => {
    createProjectPage.enterBasicValidProjectSsh(false).typeInvalidSshPrivateKey().assertCreateButtonEnabled(false);
  });

  it('should be valid private key file and have enabled button', () => {
    createProjectPage.enterBasicValidProjectSsh(false).setValidSshPrivateKeyFile().assertCreateButtonEnabled(true);
  });

  it('should be invalid private key file and have disabled button and error', () => {
    createProjectPage.enterBasicValidProjectSsh(false).setInvalidSshPrivateKeyFile().assertCreateButtonEnabled(false);
  });
});

describe('Create extended project test ssh and https', () => {
  const createProjectPage = new NewProjectCreatePage();

  beforeEach(() => {
    createProjectPage.intercept(true).visit();
  });

  it('should only show https or ssh', () => {
    createProjectPage
      .assertSshFormVisible(false)
      .assertHttpsFormVisible(true)
      .selectSshForm()
      .assertSshFormVisible(true)
      .assertHttpsFormVisible(false)
      .selectHttpsForm()
      .assertSshFormVisible(false)
      .assertHttpsFormVisible(true);
  });

  it('should keep data if switched from https to ssh to https form', () => {
    createProjectPage
      .enterFullValidProjectHttps()
      .assertCreateButtonEnabled(true)
      .selectSshForm()
      .assertCreateButtonEnabled(false)
      .selectHttpsForm()
      .validateFullValidProjectHttps();
  });

  it('should keep data if switched from ssh to https to ssh form', () => {
    createProjectPage
      .selectSshForm()
      .enterFullValidProjectSsh()
      .assertCreateButtonEnabled(true)
      .selectHttpsForm()
      .assertCreateButtonEnabled(false)
      .selectSshForm()
      .validateFullValidProjectSsh();
  });

  it('should have disabled button if switched from valid https form to invalid ssh form', () => {
    createProjectPage
      .enterBasicValidProjectHttps()
      .assertCreateButtonEnabled(true)
      .selectSshForm()
      .assertCreateButtonEnabled(false);
  });

  it('should have enabled button if switched from valid https form to invalid ssh form', () => {
    createProjectPage.enterBasicValidProjectHttps().selectSshForm().selectHttpsForm().assertCreateButtonEnabled(true);
  });

  it('should have disabled button if switched from valid ssh form to invalid https form', () => {
    createProjectPage
      .selectSshForm()
      .enterBasicValidProjectSsh()
      .assertCreateButtonEnabled(true)
      .selectHttpsForm()
      .assertCreateButtonEnabled(false);
  });

  it('should have enabled button if switched from valid ssh form to invalid https form', () => {
    createProjectPage
      .selectSshForm()
      .enterBasicValidProjectSsh()
      .selectHttpsForm()
      .selectSshForm()
      .assertCreateButtonEnabled(true);
  });
});

describe('Create extended project with and with automatic provisioned git upstream', () => {
  const createProjectPage = new NewProjectCreatePage();

  beforeEach(() => {
    createProjectPage.intercept(true, true).visit();
  });

  it('should show the no upstream option as default', () => {
    createProjectPage.assertNoUpstreamSelected(true);
  });

  it('should enable the create button, if everything is filled and no upstream is selected', () => {
    createProjectPage.enterBasicValidProjectWithoutGitUpstream().assertCreateButtonEnabled(true);
  });

  it('should disable the create button, if invalid https or ssh data is entered, and enable it again after no upstream is selected', () => {
    createProjectPage
      .assertNoUpstreamSelected(true)
      .selectHttpsForm()
      .enterBasicValidProjectHttps()
      .assertCreateButtonEnabled(true)
      .clearGitToken()
      .assertCreateButtonEnabled(false);

    createProjectPage
      .selectSshForm()
      .enterBasicValidProjectSsh()
      .assertCreateButtonEnabled(true)
      .clearSshPrivateKey()
      .assertCreateButtonEnabled(false);

    createProjectPage.selectNoUpstreamForm().assertCreateButtonEnabled(true);
  });
});
