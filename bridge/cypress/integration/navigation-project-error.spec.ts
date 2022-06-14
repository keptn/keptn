import { ProjectBoardPage } from '../support/pageobjects/ProjectBoardPage';
import EnvironmentPage from '../support/pageobjects/EnvironmentPage';
import BasePage from '../support/pageobjects/BasePage';

describe('Navigation on project error', () => {
  it('should still navigate and load data if one project request failed', () => {
    const projectBoardPage = new ProjectBoardPage();
    const environmentPage = new EnvironmentPage();
    const basePage = new BasePage();
    const project = 'my-error-project';
    environmentPage.intercept();
    projectBoardPage.interceptError(project);

    environmentPage.visit(project).assertIsLoaded(false);
    basePage.selectProjectThroughHeader('sockshop');
    environmentPage.assertIsLoaded(true);
  });
});
