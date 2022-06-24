import NewProjectCreatePage from '../support/pageobjects/NewProjectCreatePage';
import { IProject } from '../../shared/models/IProject';
import BasePage from '../support/pageobjects/BasePage';
import { interceptFailedMetadata } from '../support/intercept';

describe('Git upstream extended settings project https test', () => {
  const projectSettingsPage = new NewProjectCreatePage();

  it('should not show https or ssh form if resource service is disabled', () => {
    const project: IProject = {
      projectName: 'sockshop',
      stages: [],
      gitCredentials: {
        user: 'myGitUser',
        remoteURL: 'https://myGitURL.com',
        https: {
          insecureSkipTLS: true,
          token: '',
          proxy: {
            scheme: 'https',
            url: 'myProxyUrl:5000',
            user: 'myProxyUser',
          },
        },
      },
      shipyardVersion: '0.14',
      creationDate: '',
      shipyard: '',
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

  it('should show "Git upstream repository" headline only once', () => {
    const project: IProject = {
      projectName: 'sockshop',
      stages: [],
      gitCredentials: {
        user: 'myGitUser',
        remoteURL: 'https://myGitURL.com',
        https: {
          token: '',
          insecureSkipTLS: true,
          proxy: {
            scheme: 'https',
            url: 'myProxyUrl:5000',
            user: 'myProxyUser',
          },
        },
      },
      shipyardVersion: '0.14',
      creationDate: '',
      shipyard: '',
    };
    cy.intercept('/api/project/sockshop', {
      body: project,
    });
    projectSettingsPage.interceptSettings(true).visitSettings('sockshop').assertGitUpstreamHeadlineExistsOnce();
  });

  it('should select HTTPS and fill out inputs', () => {
    const project: IProject = {
      projectName: 'sockshop',
      stages: [],
      gitCredentials: {
        user: 'myGitUser',
        remoteURL: 'https://myGitURL.com',
        https: {
          insecureSkipTLS: true,
          token: '',
          proxy: {
            scheme: 'https',
            url: 'myProxyUrl:5000',
            user: 'myProxyUser',
          },
        },
      },
      shipyardVersion: '0.14',
      creationDate: '',
      shipyard: '',
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

  it('should submit https form and show notification', () => {
    const basePage = new BasePage();
    const project: IProject = {
      projectName: 'sockshop',
      stages: [],
      gitCredentials: {
        user: 'myGitUser',
        remoteURL: 'https://myGitURL.com',
      },
      shipyardVersion: '0.14',
      creationDate: '',
      shipyard: '',
    };
    cy.intercept('/api/project/sockshop', {
      body: project,
    });
    projectSettingsPage.interceptSettings(true).visitSettings('sockshop').typeGitToken('myToken').updateProject();
    basePage.notificationSuccessVisible('The Git upstream was changed successfully.');
  });

  it('should submit ssh form and show notification', () => {
    const basePage = new BasePage();
    const project: IProject = {
      projectName: 'sockshop',
      stages: [],
      gitCredentials: {
        user: 'myGitUser',
        remoteURL: 'ssh://myGitURL.com',
      },
      shipyardVersion: '0.14',
      creationDate: '',
      shipyard: '',
    };
    cy.intercept('/api/project/sockshop', {
      body: project,
    });
    projectSettingsPage.interceptSettings(true).visitSettings('sockshop').typeValidSshPrivateKey().updateProject();
    basePage.notificationSuccessVisible('The Git upstream was changed successfully.');
  });

  it('should select SSH', () => {
    const project: IProject = {
      projectName: 'sockshop',
      stages: [],
      gitCredentials: {
        user: 'myGitUser',
        remoteURL: 'ssh://myGitURL.com',
      },
      shipyardVersion: '0.14',
      creationDate: '',
      shipyard: '',
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

describe('Project settings with resource service disabled', () => {
  const projectSettingsPage = new NewProjectCreatePage();
  it('should show an error if the resource service is not enabled', () => {
    const project: IProject = {
      projectName: 'sockshop',
      stages: [],
      shipyardVersion: '0.14',
      creationDate: '',
      shipyard: '',
    };
    cy.intercept('/api/project/sockshop', {
      body: project,
    });
    projectSettingsPage.interceptSettings(false).visitSettings('sockshop').assertConfigurationServiceErrorExists(true);
  });
});

describe('Project settings with invalid metadata', () => {
  const projectSettingsPage = new NewProjectCreatePage();
  it('should show error if metadata endpoint does not return data', () => {
    const project: IProject = {
      projectName: 'sockshop',
      stages: [],
      gitCredentials: {
        user: 'myGitUser',
        remoteURL: 'ssh://myGitURL.com',
      },
      shipyardVersion: '0.14',
      creationDate: '',
      shipyard: '',
    };
    cy.intercept('/api/project/sockshop', {
      body: project,
    });
    projectSettingsPage.interceptSettings(true);
    interceptFailedMetadata();
    projectSettingsPage.visitSettings('sockshop').assertErrorVisible(true);
  });
});
