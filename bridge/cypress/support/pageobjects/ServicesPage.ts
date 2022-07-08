/// <reference types="cypress" />

import { SliResult } from '../../../client/app/_models/sli-result';
import { interceptProjectBoard, interceptServicesPage, interceptServicesPageWithRemediation } from '../intercept';

type SliColumn = 'name' | 'value' | 'weight' | 'score' | 'result' | 'criteria' | 'pass-criteria' | 'warning-criteria';

class ServicesPage {
  public interceptAll(): this {
    interceptProjectBoard();
    return this.intercept();
  }

  public intercept(): this {
    interceptServicesPage();
    return this;
  }

  public interceptRemediations(): this {
    interceptServicesPageWithRemediation();
    return this;
  }

  public visitServicePage(projectName: string): this {
    cy.visit(`/project/${projectName}/service`).wait('@metadata');
    return this;
  }

  public visitService(projectName: string, serviceName: string): this {
    cy.visit(`/project/${projectName}/service/${serviceName}`).wait('@metadata');
    return this;
  }

  public visitServiceDeployment(projectName: string, serviceName: string, keptnContext: string, stage?: string): this {
    let url = `/project/${projectName}/service/${serviceName}/context/${keptnContext}`;
    if (stage) {
      url += `/stage/${stage}`;
    }
    cy.visit(url).wait('@metadata');
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
    cy.byTestId(`keptn-deployment-timeline-stage-${stageName}`).click();
    return this;
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

  verifySliBreakdown(result: SliResult, isExpanded: boolean): this {
    cy.byTestId('keptn-sli-breakdown').should('exist');
    if (isExpanded) {
      this.assertSliColumnText(
        result.name,
        'name',
        `${result.name}Absolute change:Relative change:Compared with:`
      ).assertSliColumnText(
        result.name,
        'value',
        `${result.value}${result.calculatedChanges?.absolute > 0 ? '+' : '-'}${result.calculatedChanges?.absolute}${
          result.calculatedChanges?.relative > 0 ? '+' : '-'
        }${result.calculatedChanges?.relative}% ${result.comparedValue}`
      );
    } else {
      this.assertSliColumnText(result.name, 'name', result.name).assertSliColumnText(
        result.name,
        'value',
        `${result.value} (${result.calculatedChanges?.relative > 0 ? '+' : '-'}${result.calculatedChanges?.relative}%) `
      );
    }

    this.assertSliColumnText(result.name, 'weight', result.weight.toString()).assertSliColumnText(
      result.name,
      'score',
      result.score.toString()
    );

    cy.byTestId(`keptn-sli-breakdown-row-${result.name}`)
      .find('dt-cell')
      .eq(2)
      .find('.error')
      .should('have.length', result.result == 'pass' ? 0 : 1);

    return this;
  }

  assertSliColumnText(sliName: string, columnName: SliColumn, value: string): this {
    cy.byTestId(`keptn-sli-breakdown-row-${sliName}`)
      .find(`[uitestid="keptn-sli-breakdown-${columnName}-cell"]`)
      .should('have.text', value);
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
    cy.get('button[uitestid="keptn-event-item-contextButton-evaluation"]').click();
    return this;
  }

  clickViewServiceDetails(): this {
    cy.get('.highcharts-plot-background').should('be.visible');
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
    cy.get('.highcharts-plot-background').should('be.visible');
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
}
export default ServicesPage;
