import { ProjectBoardPage } from '../support/pageobjects/ProjectBoardPage';
import EnvironmentPage from '../support/pageobjects/EnvironmentPage';
import { SequencesPage } from '../support/pageobjects/SequencesPage';
import ServicesPage from '../support/pageobjects/ServicesPage';
import ProjectSettingsPage from '../support/pageobjects/ProjectSettingsPage';

describe('Navigate through project board sub pages', () => {
  const projectBoardPage = new ProjectBoardPage();
  const environmentPage = new EnvironmentPage();
  const sequencesPage = new SequencesPage();
  const servicesPage = new ServicesPage();
  const projectSettingsPage = new ProjectSettingsPage();
  const project = 'sockshop';

  beforeEach(() => {
    environmentPage.intercept();
    sequencesPage.intercept();
    servicesPage.intercept();
    projectSettingsPage.interceptSettings();
  });

  it('should have active menu button if navigated from a sub page to sequences', () => {
    environmentPage.visit(project);
    projectBoardPage.clickSequenceMenuitem().assertOnlySequencesViewSelected();
  });

  it('should have active menu button if navigated from a sub page to services', () => {
    environmentPage.visit(project);
    projectBoardPage.clickServicesMenuitem().assertOnlyServicesViewSelected();
  });

  it('should have active menu button if navigated from a sub page to environment', () => {
    sequencesPage.visit(project);
    projectBoardPage.clickEnvironmentMenuitem().assertOnlyEnvironmentViewSelected();
  });

  it('should have active menu button if navigated from a sub page to settings', () => {
    environmentPage.visit(project);
    projectBoardPage.clickSettingsMenuitem().assertOnlySettingsViewSelected();
  });
});
