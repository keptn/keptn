import NewProjectCreatePage from '../support/pageobjects/NewProjectCreatePage';
import { Project } from '../../shared/models/project';

describe('Git upstream extended settings project https test', () => {
  const projectSettingsPage = new NewProjectCreatePage();

  it('should not show https or ssh form if resource service is disabled', () => {
    const project: Project = {
      projectName: 'sockshop',
      stages: [],
      gitUser: 'myGitUser',
      gitRemoteURI: 'https://myGitURL.com',
      gitProxyInsecure: true,
      gitProxyScheme: 'https',
      gitProxyUrl: 'myProxyUrl:5000',
      gitProxyUser: 'myProxyUser',
      shipyardVersion: '0.14',
    };
    cy.intercept('/api/project/sockshop', {
      body: project,
    });
    projectSettingsPage
      .interceptSettings()
      .visitSettings('sockshop')
      .assertSshFormExists(false)
      .assertHttpsFormExists(false);
  });

  it('should select HTTPS and fill out inputs', () => {
    const project: Project = {
      projectName: 'sockshop',
      stages: [],
      gitUser: 'myGitUser',
      gitRemoteURI: 'https://myGitURL.com',
      gitProxyInsecure: true,
      gitProxyScheme: 'https',
      gitProxyUrl: 'myProxyUrl:5000',
      gitProxyUser: 'myProxyUser',
      shipyardVersion: '0.14',
    };
    cy.intercept('/api/project/sockshop', {
      body: project,
    });
    projectSettingsPage
      .interceptSettings(true)
      .visitSettings('sockshop')
      .assertGitUsername('myGitUser')
      .assertGitUrl('https://myGitURL.com')
      .assertHttpsFormVisible(true)
      .assertProxyEnabled(true)
      .assertProxyFormVisible(true)
      .assertProxyScheme('HTTPS')
      .assertProxyInsecure(true)
      .assertProxyUsername('myProxyUser')
      .assertProxyUrl('myProxyUrl')
      .assertProxyPort(5000);
  });

  it('should select SSH', () => {
    const project: Project = {
      projectName: 'sockshop',
      stages: [],
      gitProxyInsecure: false,
      gitUser: 'myGitUser',
      gitRemoteURI: 'ssh://myGitURL.com',
      shipyardVersion: '0.14',
    };
    cy.intercept('/api/project/sockshop', {
      body: project,
    });
    projectSettingsPage.interceptSettings(true).visitSettings('sockshop');

    projectSettingsPage.assertSshFormVisible(true).assertGitUsernameSsh('myGitUser');
  });

  it('should show "Set Git upstream" button', () => {
    projectSettingsPage.assertUpdateButtonExists(true);
  });
});

describe('Automatic provisioning enabled test', () => {
  const projectSettingsPage = new NewProjectCreatePage();

  beforeEach(() => {
    projectSettingsPage.interceptSettings(true, true);
    projectSettingsPage.visitSettings('sockshop');
  });

  it('should select no upstream radio button as default when no upstream was configured for a project', () => {
    const project: Project = {
      projectName: 'sockshop',
      stages: [],
      gitUser: '',
      gitRemoteURI: '',
      gitProxyInsecure: false,
      gitProxyUrl: '',
      gitProxyUser: '',
      shipyardVersion: '0.14',
    };
    cy.intercept('/api/project/sockshop', {
      body: project,
    }).as('project');

    cy.wait('@metadata').wait('@project');

    projectSettingsPage.assertNoUpstreamSelected(true);
  });

  it('should select https radio button as default if filled in and disable no upstream radio button', () => {
    const project: Project = {
      projectName: 'sockshop',
      stages: [],
      gitUser: 'myGitUser',
      gitRemoteURI: 'https://myGitURL.com',
      gitProxyInsecure: false,
      gitProxyUrl: '',
      gitProxyUser: '',
      shipyardVersion: '0.14',
    };
    cy.intercept('/api/project/sockshop', {
      body: project,
    }).as('project');

    cy.wait('@metadata').wait('@project');

    projectSettingsPage.assertHttpsFormVisible(true).assertNoUpstreamSelected(false).assertNoUpstreamEnabled(false);

    projectSettingsPage
      .enterBasicHttps()
      .assertUpdateButtonEnabled(true)
      .clearGitToken()
      .assertUpdateButtonEnabled(false);
  });

  it('should select ssh radio button as default if filled in and disable no upstream radio button', () => {
    const project: Project = {
      projectName: 'sockshop',
      stages: [],
      gitProxyInsecure: false,
      gitUser: 'myGitUser',
      gitRemoteURI: 'ssh://myGitURL.com',
      shipyardVersion: '0.14',
    };
    cy.intercept('/api/project/sockshop', {
      body: project,
    }).as('project');

    cy.wait('@metadata').wait('@project');

    projectSettingsPage.assertSshFormVisible(true).assertNoUpstreamSelected(false).assertNoUpstreamEnabled(false);

    projectSettingsPage
      .enterBasicSsh()
      .assertUpdateButtonEnabled(true)
      .clearSshPrivateKey()
      .assertUpdateButtonEnabled(false);
  });
});
