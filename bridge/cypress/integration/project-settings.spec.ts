import NewProjectCreatePage from '../support/pageobjects/NewProjectCreatePage';
import { Project } from '../../shared/models/project';

describe('Create extended project https test', () => {
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
    projectSettingsPage
      .interceptSettings(true)
      .visitSettings('sockshop')
      .assertSshFormVisible(true)
      .assertGitUsernameSsh('myGitUser');
  });

  it('should show "Set Git upstream" button', () => {
    projectSettingsPage.assertUpdateButtonExists(true);
  });
});
