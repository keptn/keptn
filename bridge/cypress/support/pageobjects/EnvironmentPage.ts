/// <reference types="cypress" />

import { interceptEmptyEnvironmentScreen, interceptEnvironmentScreen } from '../intercept';

class EnvironmentPage {
  public intercept(): this {
    interceptEnvironmentScreen();
    return this;
  }

  public interceptEmpty(): this {
    interceptEmptyEnvironmentScreen();
    return this;
  }

  public visit(project: string): this {
    cy.visit(`/project/${project}`).wait('@metadata').wait('@project');
    return this;
  }

  public clickCreateService(stage: string): this {
    cy.get('ktb-selectable-tile h2')
      .contains(stage)
      .parentsUntil('ktb-selectable-tile')
      .find('ktb-no-service-info a')
      .click();
    return this;
  }

  public selectStage(stage: string): this {
    cy.get('ktb-selectable-tile h2').contains(stage).click();
    return this;
  }

  public assertEvaluationHistoryLoadingCount(service: string, count: number): this {
    this.getServiceDetailsContainer(service)
      .find('ktb-evaluation-info dt-tag-list[aria-label="evaluation-history"] dt-tag ktb-loading-spinner')
      .should('have.length', count);
    return this;
  }

  public assertEvaluationHistoryCount(service: string, count: number): this {
    const tags = this.getServiceDetailsContainer(service).find(
      'ktb-evaluation-info dt-tag-list[aria-label="evaluation-history"] dt-tag'
    );
    tags.should('have.length', count);
    if (count !== 0) {
      tags.should('have.class', 'border');
    }
    return this;
  }

  public assertEvaluationInDetails(service: string, score: number | '-', status?: 'success' | 'error'): this {
    const evaluationTag = this.getServiceDetailsContainer(service).find(
      'ktb-evaluation-info dt-tag-list[aria-label="evaluation-info"] dt-tag'
    );
    evaluationTag.should('not.have.class', 'border').should('have.text', score);
    if (status) {
      evaluationTag.should('have.class', status);
    }
    return this;
  }

  public getEvaluationHistoryURL(project: string, stage: string, service: string, limit: 5 | 6): string {
    return `/api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?filter=data.project:${project}%20AND%20data.service:${service}%20AND%20data.stage:${stage}%20AND%20source:lighthouse-service&excludeInvalidated=true&limit=${limit}`;
  }

  private getServiceDetailsContainer(service: string): Cypress.Chainable<JQuery<HTMLElement>> {
    return cy.get('ktb-stage-details ktb-expandable-tile h2').contains(service).parentsUntil('ktb-expandable-tile');
  }

  public assertIsLoaded(status: boolean): this {
    cy.byTestId('ktb-environment-is-loading').should(status ? 'not.exist' : 'exist');
    return this;
  }
}

export default EnvironmentPage;
