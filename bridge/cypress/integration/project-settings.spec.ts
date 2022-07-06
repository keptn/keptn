import ProjectSettingsPage from '../support/pageobjects/ProjectSettingsPage';
import { IProject } from '../../shared/interfaces/project';
import BasePage from '../support/pageobjects/BasePage';
import { interceptFailedMetadata } from '../support/intercept';
import { ProjectBoardPage } from '../support/pageobjects/ProjectBoardPage';

describe('Git upstream extended settings project https test', () => {
  const projectSettingsPage = new ProjectSettingsPage();

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

    projectSettingsPage
      .interceptSettings()
      .interceptProject(project)
      .visitSettings('sockshop')
      .assertSshFormExists(false)
      .assertHttpsFormExists(false);
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

    projectSettingsPage
      .interceptSettings(true)
      .interceptProject(project)
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

    projectSettingsPage
      .interceptSettings(true)
      .interceptProject(project)
      .visitSettings('sockshop')
      .typeGitToken('myToken')
      .updateProject();
    basePage.notificationSuccessVisible('The Git upstream was changed successfully.');
  });

  it('should prevent data loss if git credentials are not saved before navigation', () => {
    const basePage = new BasePage();
    const projectBoardPage = new ProjectBoardPage();
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

    projectSettingsPage
      .interceptSettings(true)
      .interceptProject(project)
      .visitSettings('sockshop')
      .typeGitToken('myToken');
    projectBoardPage.goToServicesPage();
    projectSettingsPage.clickSaveChangesPopup();
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

    projectSettingsPage
      .interceptSettings(true)
      .interceptProject(project)
      .visitSettings('sockshop')
      .typeValidSshPrivateKey()
      .updateProject();
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

    projectSettingsPage
      .interceptSettings(true)
      .interceptProject(project)
      .visitSettings('sockshop')
      .assertSshFormVisible(true)
      .assertGitUsernameSsh('myGitUser');
  });

  it('should show "Set Git upstream" button', () => {
    projectSettingsPage.assertUpdateButtonExists(true);
  });
});

describe('Project settings with resource service disabled', () => {
  const projectSettingsPage = new ProjectSettingsPage();
  it('should show an error if the resource service is not enabled', () => {
    const project: IProject = {
      projectName: 'sockshop',
      stages: [],
      shipyardVersion: '0.14',
      creationDate: '',
      shipyard: '',
    };

    projectSettingsPage
      .interceptSettings(false)
      .interceptProject(project)
      .visitSettings('sockshop')
      .assertConfigurationServiceErrorExists(true);
  });
});

describe('Project settings with invalid metadata', () => {
  const projectSettingsPage = new ProjectSettingsPage();
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

    projectSettingsPage.interceptSettings(true).interceptProject(project);
    interceptFailedMetadata();
    projectSettingsPage.visitSettings('sockshop').assertErrorVisible(true);
  });
});
