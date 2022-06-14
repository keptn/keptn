/// <reference types="cypress" />

import { Project } from '../../../client/app/_models/project';
import { interceptDashboard } from '../intercept';
import EnvironmentPage from './EnvironmentPage';
import NewProjectCreatePage from './NewProjectCreatePage';

class DashboardPage {
  private PROJECT_TILE_LOCATOR = 'dt-tile[id="proj_pattern"]';

  public intercept(): this {
    interceptDashboard();
    return this;
  }

  public visit(waitForProjects = true): this {
    cy.visit(`/`).wait('@metadata');
    if (waitForProjects) {
      cy.wait('@projects');
    }
    return this;
  }

  public clickProjectTile(projectName: string): EnvironmentPage {
    cy.wait(500).get(this.PROJECT_TILE_LOCATOR.replace('proj_pattern', projectName)).click();
    return new EnvironmentPage();
  }

  public assertProjects(projects: Project[]): this {
    cy.get('ktb-project-tile').should('have.length', projects.length);
    projects.forEach((project, index) => {
      cy.get('ktb-project-tile').eq(index).find('dt-tile-title').should('contain.text', project.projectName);
      cy.get('ktb-project-tile')
        .eq(index)
        .byTestId('keptn-project-tile-numStagesServices')
        .should(
          'contain.text',
          `${project.stages.length} Stages, ${project.stages[0]?.services.length ?? 0} Services `
        );
    });
    return this;
  }

  public clickCreateNewProjectButton(): NewProjectCreatePage {
    cy.get('.dt-button-primary > span.dt-button-label').contains('Create a new project').click();
    return new NewProjectCreatePage();
  }
}

export default DashboardPage;
