import EnvironmentPage from '../support/pageobjects/EnvironmentPage';
import { ServicesSettingsPage } from '../support/pageobjects/ServicesSettingsPage';
import ServicesPage from '../support/pageobjects/ServicesPage';
import { EvaluationBadgeVariant } from '../../client/app/_components/ktb-evaluation-badge/ktb-evaluation-badge.utils';

describe('Environment Screen empty', () => {
  const environmentPage = new EnvironmentPage();
  beforeEach(() => {
    environmentPage.interceptEmpty().visit('dynatrace');
  });

  it('should redirect to create service and redirect back after creation', () => {
    const serviceSettings = new ServicesSettingsPage();

    environmentPage.clickCreateService('dev');
    cy.location('pathname').should('eq', '/project/dynatrace/settings/services/create');
    serviceSettings.createService('my-new-service');
    cy.location('pathname').should('eq', '/project/dynatrace');
  });
});

describe('Environment Screen default requests', () => {
  const environmentPage = new EnvironmentPage();
  const project = 'sockshop';
  const stage = 'dev';

  beforeEach(() => {
    environmentPage.intercept().visit(project);
  });

  it('should not show evaluation history loading indicators', () => {
    environmentPage
      .selectStage(stage)
      .waitForEvaluationHistory('carts', stage, 6)
      .assertEvaluationHistoryLoadingCount('carts', 0);
  });

  it('should not show evaluation history', () => {
    environmentPage.selectStage(stage).assertEvaluationHistoryCount('carts', 0);
  });

  it('should not show evaluation', () => {
    environmentPage
      .selectStage(stage)
      .assertEvaluationInDetails('carts-db', '-', undefined, EvaluationBadgeVariant.NONE);
  });

  it('stage-detail component should exist after clicking on stage', () => {
    environmentPage.selectStage(stage);
    environmentPage.assertStageDetailsHeader(stage);
  });

  it('stage-detail component should exist when navigating to /environment/stage url', () => {
    environmentPage.visit(project, stage);
    environmentPage.assertStageDetailsHeader(stage);
  });

  it('filter should be set when navigating to /environment/stage?filterType=filter', () => {
    environmentPage.visit(project, 'staging', 'evaluation');
    environmentPage.assertStageDetailsFilterEnabled('evaluation', true);
  });

  it('should redirect to stage', () => {
    environmentPage.selectStage('dev');
    cy.location('pathname').should('eq', `/project/${project}/environment/stage/dev`);
    environmentPage.selectStage('staging');
    cy.location('pathname').should('eq', `/project/${project}/environment/stage/staging`);
  });

  it('should add query parameter if clicking on type', () => {
    environmentPage.clickFilterType('staging', 'evaluation');
    cy.location('href').should('include', `/project/${project}/environment/stage/staging?filterType=evaluation`);
  });

  it('should show pause icon if sequence is paused', () => {
    environmentPage.assertPauseIconShown('staging', 'carts-db');
  });
});

describe('Environment Screen Navigation', () => {
  const environmentPage = new EnvironmentPage();
  const servicesPage = new ServicesPage();
  const project = 'sockshop';
  const stage = 'production';
  const service = 'carts';
  const keptnContext = 'da740469-9920-4e0c-b304-0fd4b18d17c2';

  beforeEach(() => {
    environmentPage.intercept();
    servicesPage.intercept();

    environmentPage.visit(project);
  });

  it('navigate to service if clicking on service from stage-overview', () => {
    environmentPage.clickServiceFromStageOverview(stage, service);
    cy.wait('@serviceDatastore');
    servicesPage.assertDeploymentDeepLink(project, service, keptnContext, stage);
  });

  it('navigate to service if clicking on service from stage-details', () => {
    environmentPage.clickServiceFromStageDetails(stage, service);
    cy.wait('@serviceDatastore');
    servicesPage.assertDeploymentDeepLink(project, service, keptnContext, stage);
  });
});

describe('Environment Screen dynamic requests', () => {
  const environmentPage = new EnvironmentPage();
  const project = 'sockshop';
  const stage = 'dev';

  beforeEach(() => {
    environmentPage.intercept();
  });

  it('should show evaluation history loading indicators', () => {
    const service = 'carts';
    environmentPage
      .interceptEvaluationHistory(project, stage, service, 6, 10_000)
      .visit(project)
      .selectStage(stage)
      .assertEvaluationHistoryLoadingCount(service, 5);
  });

  it('should show evaluations in history if sequence does not have an evaluation task', () => {
    const service = 'carts-db';
    environmentPage
      .interceptEvaluationHistory(
        project,
        stage,
        service,
        5,
        undefined,
        'get.environment.evaluation.history.carts-db.mock'
      )
      .visit(project)
      .selectStage(stage)
      .assertEvaluationHistoryLoadingCount(service, 0)
      .assertEvaluationHistoryCount(service, 5)
      .assertEvaluationInDetails(service, '-', undefined, EvaluationBadgeVariant.NONE);
  });

  it('should show 2 evaluations in history and should not show current evaluation in history', () => {
    const service = 'carts';
    environmentPage
      .interceptEvaluationHistory(
        project,
        stage,
        service,
        6,
        undefined,
        'get.environment.evaluation.history.limited.mock'
      )
      .visit(project)
      .selectStage(stage)
      .assertEvaluationHistoryCount(service, 2)
      .assertEvaluationInDetails(service, 100, 'success', EvaluationBadgeVariant.FILL);
  });

  it('should show 5 evaluation in history', () => {
    const service = 'carts';
    environmentPage
      .interceptEvaluationHistory(project, stage, service, 6, undefined, 'get.environment.evaluation.history.mock')
      .visit(project)
      .selectStage(stage)
      .assertEvaluationHistoryCount(service, 5)
      .assertEvaluationInDetails(service, 100, 'success', EvaluationBadgeVariant.FILL);
  });
});
