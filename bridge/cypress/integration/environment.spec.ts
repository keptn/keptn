import EnvironmentPage from '../support/pageobjects/EnvironmentPage';
import { ServicesSettingsPage } from '../support/pageobjects/ServicesSettingsPage';
import ServicesPage from '../support/pageobjects/ServicesPage';

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
    environmentPage.selectStage(stage).assertEvaluationHistoryLoadingCount('carts', 0);
  });

  it('should not show evaluation history', () => {
    environmentPage.selectStage(stage).assertEvaluationHistoryCount('carts', 0);
  });

  it('should not show evaluation', () => {
    environmentPage.selectStage(stage).assertEvaluationInDetails('carts-db', '-');
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
    cy.intercept(environmentPage.getEvaluationHistoryURL(project, stage, service, 6), {
      delay: 10_000,
    });
    environmentPage.visit(project).selectStage(stage).assertEvaluationHistoryLoadingCount(service, 5);
  });

  it('should show evaluation history loading indicators', () => {
    const service = 'carts';
    cy.intercept(environmentPage.getEvaluationHistoryURL(project, stage, service, 6), {
      delay: 10_000,
    });
    environmentPage.visit(project).selectStage(stage).assertEvaluationHistoryLoadingCount(service, 5);
  });

  it('should show evaluations in history if sequence does not have an evaluation task', () => {
    const service = 'carts-db';
    cy.intercept(environmentPage.getEvaluationHistoryURL(project, stage, service, 5), {
      fixture: 'get.environment.evaluation.history.carts-db.mock',
    });
    environmentPage
      .visit(project)
      .selectStage(stage)
      .assertEvaluationHistoryLoadingCount(service, 0)
      .assertEvaluationHistoryCount(service, 5)
      .assertEvaluationInDetails(service, '-');
  });

  it.only('should show 2 evaluations in history and should not show current evaluation in history', () => {
    const service = 'carts';
    cy.intercept(environmentPage.getEvaluationHistoryURL(project, stage, service, 6), {
      fixture: 'get.environment.evaluation.history.limited.mock', // 3 events, including the current one
    });
    environmentPage
      .visit(project)
      .selectStage(stage)
      .assertEvaluationHistoryCount(service, 2)
      .assertEvaluationInDetails(service, 0, 'success');
  });

  it('should show 5 evaluation in history', () => {
    const service = 'carts';
    cy.intercept(environmentPage.getEvaluationHistoryURL(project, stage, service, 6), {
      fixture: 'get.environment.evaluation.history.mock',
    });
    environmentPage
      .visit(project)
      .selectStage(stage)
      .assertEvaluationHistoryCount(service, 5)
      .assertEvaluationInDetails(service, 0, 'success');
  });
});
