/// <reference types="cypress" />

import { Project } from '../../../client/app/_models/project';

class DashboardPage {
  visitDashboard(): this {
    cy.visit(`/`);
    return this;
  }

  assertProjects(projects: Project[]): this {
    cy.get('ktb-project-tile').should('have.length', projects.length);
    projects.forEach((project, index) => {
      cy.get('ktb-project-tile').eq(index).find('dt-tile-title').should('contain.text', project.projectName);
      cy.get('ktb-project-tile')
        .eq(index)
        .byTestId('keptn-project-tile-numStagesServices')
        .should('contain.text', `${project.stages.length} Stages, ${project.stages[0].services.length} Services `);
    });
    return this;
  }
}

export default DashboardPage;
