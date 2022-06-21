import EnvironmentPage from '../support/pageobjects/EnvironmentPage';
import { ServicesSettingsPage } from '../support/pageobjects/ServicesSettingsPage';

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
    cy.get('ktb-stage-details h2').should('contain.text', stage);
  });

  it('stage-detail component should exist when navigating to /environment/stage url', () => {
    environmentPage.visit(project, stage);
    cy.get('ktb-stage-details h2').should('contain.text', stage);
  });

  it('should redirect to stage', () => {
    environmentPage.selectStage('dev');
    cy.location('pathname').should('eq', `/project/${project}/environment/stage/dev`);
    environmentPage.selectStage('staging');
    cy.location('pathname').should('eq', `/project/${project}/environment/stage/staging`);
  });

  it('should add query parameter if clicking on type', () => {
    environmentPage.clickFilterType(stage, 'problem');
    cy.location('href').should('include', `/project/${project}/environment/stage/${stage}?filterType=problem`);
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

  it('should show 2 evaluations in history and should not show current evaluation in history', () => {
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
