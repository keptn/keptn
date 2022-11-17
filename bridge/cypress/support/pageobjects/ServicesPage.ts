/// <reference types="cypress" />

import { SliResult } from '../../../client/app/_interfaces/sli-result';
import {
  interceptHeatmapComponent,
  interceptProjectBoard,
  interceptServicesPage,
  interceptServicesPageWithLoadingSequences,
  interceptServicesPageWithRemediation,
} from '../intercept';
import { EvaluationBadgeVariant } from '../../../client/app/_components/ktb-evaluation-badge/ktb-evaluation-badge.utils';

type SliColumn = 'name' | 'value' | 'weight' | 'score' | 'result' | 'criteria' | 'pass-criteria' | 'warning-criteria';
type UISliResult = SliResult & {
  availableScore: number;
  calculatedChanges: {
    absolute: number;
    relative: number | undefined;
  };
};

class ServicesPage {
  public interceptAll(): this {
    interceptProjectBoard();
    return this.intercept();
  }

  public intercept(): this {
    interceptServicesPage();
    return this;
  }

  public interceptRunning(): this {
    interceptServicesPageWithLoadingSequences();
    return this;
  }

  public interceptRemediations(): this {
    interceptServicesPageWithRemediation();
    return this;
  }

  public interceptForEvaluationBadge(): this {
    cy.intercept('GET', '/api/project/sockshop/deployment/da740469-9920-4e0c-b304-0fd4b18d17c2', {
      statusCode: 200,
      fixture: 'get.sockshop.service.carts.deployment.evaluation.badge.mock.json',
    }).as('ServiceDeployment');
    return this;
  }

  public interceptSliFallback(projectName: string, comparedEventIds: string[], addDelay = false): this {
    cy.intercept(
      'GET',
      `api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?filter=data.project:${projectName}%20AND%20source:lighthouse-service%20AND%20id:${comparedEventIds.join(
        ','
      )}&excludeInvalidated=true&limit=${comparedEventIds.length}`,
      {
        statusCode: 200,
        fixture: 'get.sockshop.service.carts.evaluation.compared.event.mock',
        delay: addDelay ? 10_000 : 0,
      }
    ).as('sliFallback');
    return this;
  }

  public interceptWithInfoSli(): this {
    interceptHeatmapComponent();
    cy.intercept('GET', 'api/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?*', {
      statusCode: 200,
      fixture: 'get.sockshop.service.carts.info-evaluations.mock.json',
    }).as('heatmapEvaluations');
    return this;
  }

  public visitServicePage(projectName: string): this {
    cy.visit(`/project/${projectName}/service`).wait('@bridgeInfo');
    return this;
  }

  public visitService(projectName: string, serviceName: string): this {
    cy.visit(`/project/${projectName}/service/${serviceName}`).wait('@bridgeInfo');
    return this;
  }

  public visitServiceDeployment(projectName: string, serviceName: string, keptnContext: string, stage?: string): this {
    let url = `/project/${projectName}/service/${serviceName}/context/${keptnContext}`;
    if (stage) {
      url += `/stage/${stage}`;
    }
    cy.visit(url).wait('@bridgeInfo');
    return this;
  }

  public waitForSliFallbackFetch(): this {
    cy.wait('@sliFallback');
    return this;
  }

  public waitForEvaluations(): this {
    cy.wait('@serviceDatastore').get('ktb-heatmap').should('exist'); // wait until heatmap is rendered
    return this;
  }

  public assertServiceExpanded(serviceName: string, status: boolean): this {
    cy.byTestId(`ktb-service-tile-${serviceName}-content`).should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertDeploymentSelected(serviceName: string, version: string, status: boolean): this {
    cy.byTestId(`ktb-service-${serviceName}-deployment-${version}`).should(
      status ? 'have.class' : 'not.have.class',
      'active'
    );
    return this;
  }

  public assertStageSelected(stageName: string, status: boolean): this {
    cy.get('ktb-deployment-timeline>div>div ktb-stage-badge dt-tag')
      .contains(stageName)
      .should(status ? 'have.class' : 'not.have.class', 'focused');
    return this;
  }

  public selectService(serviceName: string, version: string): this {
    cy.byTestId(`keptn-service-view-service-${serviceName}`).click();
    cy.get('dt-row').contains(version).click();
    cy.wait('@ServiceDeployment');
    return this;
  }

  public selectStage(stageName: string): this {
    this.getStageInTimeline(stageName).click();
    return this;
  }

  private getStageInTimeline(stageName: string): Cypress.Chainable<JQuery> {
    return cy.byTestId(`keptn-deployment-timeline-stage-${stageName}`);
  }

  public clickSliBreakdownHeader(columnName: string): this {
    cy.byTestId('keptn-sli-breakdown')
      .find('dt-header-row')
      .first()
      .find('dt-header-cell')
      .contains(columnName)
      .click();
    return this;
  }

  verifySliBreakdownSorting(
    columnIndex: number,
    direction: 'ascending' | 'descending',
    firstElement: string,
    secondElement: string
  ): this {
    cy.byTestId('keptn-sli-breakdown')
      .find('dt-header-row')
      .first()
      .find('dt-header-cell')
      .eq(columnIndex)
      .find('.dt-sort-header-container')
      .first()
      .should('have.class', 'dt-sort-header-sorted');
    cy.byTestId('keptn-sli-breakdown')
      .find('dt-header-row')
      .first()
      .find('dt-header-cell')
      .eq(columnIndex)
      .should('have.attr', 'aria-sort', direction);

    cy.byTestId('keptn-sli-breakdown')
      .find('dt-row')
      .eq(0)
      .find('dt-cell')
      .eq(columnIndex)
      .should('have.text', firstElement);
    cy.byTestId('keptn-sli-breakdown')
      .find('dt-row')
      .eq(1)
      .find('dt-cell')
      .eq(columnIndex)
      .should('have.text', secondElement);

    return this;
  }

  verifySliBreakdown(result: UISliResult, isExpanded: boolean): this {
    cy.byTestId('keptn-sli-breakdown').should('exist');
    if (isExpanded) {
      this.assertSliColumnText(
        result.name,
        'name',
        `${result.name}Absolute change:Relative change:Compared with:`
      ).assertSliValueColumnExpanded(
        result.name,
        +result.value,
        result.calculatedChanges.absolute,
        result.calculatedChanges.relative,
        result.comparedValue ?? 0
      );
    } else {
      let assertVal: string;
      if (result.calculatedChanges.relative) {
        assertVal = `${result.value} (${result.calculatedChanges.relative > 0 ? '+' : '-'}${
          result.calculatedChanges.relative
        }%) `;
      } else {
        assertVal = `${result.value}`;
      }
      this.assertSliColumnText(result.name, 'name', result.name).assertSliColumnText(result.name, 'value', assertVal);
    }

    this.assertSliColumnText(result.name, 'weight', result.weight.toString()).assertSliScoreColumn(
      result.name,
      result.score,
      result.availableScore
    );
    cy.byTestId(`keptn-sli-breakdown-row-${result.name}`)
      .find('dt-cell')
      .eq(2)
      .find('.error')
      .should('have.length', !result.calculatedChanges.relative || result.result == 'pass' ? 0 : 1);

    return this;
  }

  public assertSliValueColumnExpanded(
    sliName: string,
    value: number,
    absoluteChange: number,
    relativeChange: number | undefined,
    comparedValue: number
  ): this {
    const columnName = 'value';
    this._getSliCell(sliName, columnName).byTestId('ktb-sli-breakdown-value-value').should('have.text', value);
    this._getSliCell(sliName, columnName)
      .byTestId('ktb-sli-breakdown-value-absolute')
      .should('have.text', `${absoluteChange > 0 ? '+' : ''}${absoluteChange}`);
    this.assertSliRelativeChange(sliName, relativeChange ?? 'n/a');
    this._getSliCell(sliName, columnName)
      .byTestId('ktb-sli-breakdown-value-compared')
      .should('have.text', comparedValue);
    return this;
  }

  public assertSliRelativeChange(sliName: string, relativeChange: number | 'n/a'): this {
    let changeValue: string;
    if (typeof relativeChange === 'string') {
      changeValue = relativeChange;
    } else {
      changeValue = `${relativeChange > 0 ? '+' : ''}${relativeChange}% `;
    }
    this._getSliCell(sliName, 'value').byTestId('ktb-sli-breakdown-value-relative').should('have.text', changeValue);
    return this;
  }

  public assertSliScoreColumn(sliName: string, score?: number, availableScore?: number, isInfo?: boolean): this {
    return this.assertSliColumnText(sliName, 'score', isInfo ? '-' : `${score}/${availableScore}`);
  }

  private showSliScoreOverlay(sliName: string): this {
    this._getSliCell(sliName, 'score').trigger('mouseenter');
    return this;
  }

  private assertSliScoreOverlayFailedExists(status: boolean): this {
    cy.byTestId('ktb-sli-breakdown-score-overlay-failed').should(status ? 'exist' : 'not.exist');
    return this;
  }

  private assertSliScoreOverlayWarningExists(status: boolean): this {
    cy.byTestId('ktb-sli-breakdown-score-overlay-warning').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertSliScoreOverlayDefault(sliName: string): this {
    this.showSliScoreOverlay(sliName)
      .assertSliScoreOverlayFailedExists(false)
      .assertSliScoreOverlayWarningExists(false);
    cy.get('.dt-overlay-container').should('exist');
    return this;
  }

  public assertSliScoreOverlayWarning(sliName: string): this {
    return this.showSliScoreOverlay(sliName)
      .assertSliScoreOverlayFailedExists(false)
      .assertSliScoreOverlayWarningExists(true);
  }

  public assertSliScoreOverlayFailed(sliName: string): this {
    return this.showSliScoreOverlay(sliName)
      .assertSliScoreOverlayFailedExists(true)
      .assertSliScoreOverlayWarningExists(false);
  }

  private _getSliCell(sliName: string, columnName: SliColumn): Cypress.Chainable<JQuery> {
    return cy.byTestId(`keptn-sli-breakdown-row-${sliName}`).byTestId(`keptn-sli-breakdown-${columnName}-cell`);
  }

  assertSliColumnText(sliName: string, columnName: SliColumn, value: string): this {
    this._getSliCell(sliName, columnName).should('have.text', value);
    return this;
  }

  expandSliBreakdown(name: string): this {
    cy.byTestId(`keptn-sli-breakdown-row-${name}`).find('dt-cell').eq(0).find('button').click();
    return this;
  }

  clickOnServicePanelByName(serviceName: string): this {
    cy.wait(500).get('div.dt-info-group-content').get('h2').contains(serviceName).click();
    return this;
  }

  clickOnServiceInnerPanelByName(serviceName: string): this {
    cy.wait(500).get('span.ng-star-inserted').contains(serviceName).click();
    return this;
  }

  clickEvaluationBoardButton(): this {
    cy.byTestId('keptn-event-item-contextButton-evaluation').click();
    return this;
  }

  clickViewServiceDetails(): this {
    cy.get('ktb-heatmap .heatmap-container').should('be.visible');
    cy.contains('View service details').click();
    return this;
  }

  clickViewSequenceDetails(): this {
    cy.contains('View sequence details').click();
    return this;
  }

  clickGoBack(): this {
    cy.contains('Go back').click();
    return this;
  }

  verifyCurrentOpenServiceNameEvaluationPanel(serviceName: string): this {
    cy.get('div.service-title > span').should('have.text', serviceName);
    cy.get('ktb-heatmap .heatmap-container').should('be.visible');
    return this;
  }

  public assertRemediationSequenceCount(count: number): this {
    cy.byTestId('ktb-sequence-list-item-remediation').should('have.length', count);
    return this;
  }

  public assertRootDeepLink(project: string): this {
    cy.location('pathname').should('eq', `/project/${project}/service`);
    return this;
  }

  public assertServiceDeepLink(project: string, service: string): this {
    cy.location('pathname').should('eq', `/project/${project}/service/${service}`);
    return this;
  }

  public assertDeploymentDeepLink(project: string, service: string, keptnContext: string, stage: string): this {
    cy.location('pathname').should(
      'eq',
      `/project/${project}/service/${service}/context/${keptnContext}/stage/${stage}`
    );
    return this;
  }

  public assertIsStageLoading(stage: string, status: boolean): this {
    cy.byTestId(`ktb-deployment-stage-${stage}-loading`).should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertStageEvaluationBadge(
    stage: string,
    status: 'success' | 'error' | 'warning' | undefined,
    score: number | '-',
    variant: EvaluationBadgeVariant
  ): this {
    this.getStageInTimeline(stage).assertEvaluationBadge(status, score, variant);
    return this;
  }

  public assertSliBreakdownLoading(status: boolean): this {
    cy.byTestId('ktb-sli-breakdown-loading').should(status ? 'exist' : 'not.exist');
    return this;
  }
}
export default ServicesPage;
