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

  public visit(project: string, stage = '', filterType = ''): this {
    const query = filterType ? `?filterType=${filterType}` : '';
    cy.visit(stage ? `/project/${project}/environment/stage/${stage}${query}` : `/project/${project}`)
      .wait('@metadata')
      .wait('@project');
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

  public clickFilterType(stage: string, filterType: 'problem' | 'evaluation' | 'approval'): this {
    cy.byTestId(`filter-type-${stage}-${filterType}`).click();
    return this;
  }

  public clickServiceFromStageOverview(stage: string, service: string): this {
    cy.get('ktb-selectable-tile h2')
      .contains(stage)
      .parentsUntil('ktb-selectable-tile')
      .find('ktb-services-list')
      .contains(service)
      .click();
    return this;
  }

  public clickServiceFromStageDetails(stage: string, service: string): this {
    this.selectStage(stage);
    cy.get('ktb-expandable-tile dt-info-group a').contains(service).click();
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

  public assertStageDetailsHeader(stage: string): this {
    cy.get('ktb-stage-details h2').should('contain.text', stage);
    return this;
  }

  public assertStageDetailsFilterEnabled(filterType: 'problem' | 'evaluation' | 'approval', enabled: boolean): this {
    cy.byTestId(`ktb-stage-details-${filterType}-button`).should(enabled ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public assertRootDeepLink(project: string): this {
    cy.location('pathname').should('eq', `/project/${project}`);
    return this;
  }

  public assertStageDeepLink(project: string, stage: string): this {
    cy.location('pathname').should('eq', `/project/${project}/environment/stage/${stage}`);
    return this;
  }

  public assertPauseIconShown(stage: string, service: string): this {
    cy.get('ktb-selectable-tile h2')
      .contains(stage)
      .parentsUntil('ktb-selectable-tile')
      .find('ktb-services-list')
      .contains(service)
      .parentsUntil('dt-table')
      .find('path')
      .should('have.attr', 'd', 'M112 64h104v384H112zM296 64h104v384H296z');
    return this;
  }
}

export default EnvironmentPage;
