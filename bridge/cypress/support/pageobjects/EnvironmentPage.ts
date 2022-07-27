/// <reference types="cypress" />

import { interceptEmptyEnvironmentScreen, interceptEnvironmentScreen } from '../intercept';
import { EvaluationBadgeVariant } from '../../../client/app/_components/ktb-evaluation-badge/ktb-evaluation-badge.utils';

class EnvironmentPage {
  public intercept(): this {
    interceptEnvironmentScreen();
    return this;
  }

  public interceptEmpty(): this {
    interceptEmptyEnvironmentScreen();
    return this;
  }

  public interceptEvaluationHistory(
    project: string,
    stage: string,
    service: string,
    limit: 5 | 6,
    delay?: number,
    fixture = 'get.environment.evaluation.history.mock'
  ): this {
    cy.intercept(this.getEvaluationHistoryURL(project, stage, service, limit), {
      fixture,
      delay,
    }).as(`evaluationHistory-${service}-${stage}-${limit}`);
    return this;
  }

  public visit(project: string, stage = '', filterType = ''): this {
    const query = filterType ? `?filterType=${filterType}` : '';
    cy.visit(stage ? `/project/${project}/environment/stage/${stage}${query}` : `/project/${project}`)
      .wait('@metadata')
      .wait('@project');
    return this;
  }

  public waitForEvaluationHistory(service: string, stage: string, limit: number): this {
    cy.wait(`@evaluationHistory-${service}-${stage}-${limit}`);
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
      .byTestId('ktb-evaluation-badge-history')
      .find('ktb-loading-spinner')
      .should('have.length', count);
    return this;
  }

  public assertEvaluationHistoryCount(service: string, count: number): this {
    this.getServiceDetailsContainer(service)
      .byTestId('ktb-evaluation-badge-history')
      .find('.badge')
      .should('have.length', count);
    return this;
  }

  public assertEvaluationInDetails(
    service: string,
    score: number | '-',
    status: 'success' | 'error' | 'warning' | undefined,
    variant: EvaluationBadgeVariant
  ): this {
    this.getServiceDetailsContainer(service)
      .find('.current-evaluation')
      .assertEvaluationBadge(status, score, variant)
      .should('not.have.class', 'border')
      .should('have.text', score);
    return this;
  }

  public assertEvaluationHistory(
    service: string,
    history: { score: number | '-'; status?: 'success' | 'error' | 'warning'; variant: EvaluationBadgeVariant }[]
  ): this {
    for (let i = 0; i < history.length; ++i) {
      const evaluation = history[i];
      this.getServiceDetailsContainer(service).assertEvaluationBadge(
        evaluation.status,
        evaluation.score,
        evaluation.variant,
        i
      );
    }
    return this;
  }

  public assertEvaluationInOverview(
    stage: string,
    service: string,
    score: number | '-',
    status: 'success' | 'error' | 'warning' | undefined,
    variant: EvaluationBadgeVariant
  ): this {
    this.getServiceInStageOverview(stage, service)
      .assertEvaluationBadge(status, score, variant)
      .should('not.have.class', 'border')
      .should('have.text', score);
    return this;
  }

  public getEvaluationHistoryURL(project: string, stage: string, service: string, limit: 5 | 6): string {
    return `/api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?filter=data.project:${project}%20AND%20data.service:${service}%20AND%20data.stage:${stage}%20AND%20source:lighthouse-service&excludeInvalidated=true&limit=${limit}`;
  }

  private getServiceDetailsContainer(service: string): Cypress.Chainable<JQuery> {
    return cy
      .get('ktb-stage-details ktb-expandable-tile h2')
      .contains(service)
      .parentsUntil('ktb-expandable-tile-header');
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
    this.getServiceInStageOverview(stage, service).assertDtIcon('pause');
    return this;
  }

  private getServiceInStageOverview(stage: string, service: string): Cypress.Chainable<JQuery> {
    return cy
      .get('ktb-selectable-tile h2')
      .contains(stage)
      .parentsUntil('ktb-selectable-tile')
      .find('ktb-services-list')
      .contains(service)
      .parentsUntil('dt-row')
      .parent();
  }
}

export default EnvironmentPage;
