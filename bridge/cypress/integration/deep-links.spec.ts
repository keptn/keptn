import { ProjectBoardPage } from '../support/pageobjects/ProjectBoardPage';
import EnvironmentPage from '../support/pageobjects/EnvironmentPage';
import DashboardPage from '../support/pageobjects/DashboardPage';
import ServicesPage from '../support/pageobjects/ServicesPage';
import { SequencesPage } from '../support/pageobjects/SequencesPage';
import { EventTypes } from '../../shared/interfaces/event-types';
import BasePage from '../support/pageobjects/BasePage';

describe('Test deep links', () => {
  const projectBoardPage = new ProjectBoardPage();
  const environmentPage = new EnvironmentPage();
  const dashboardPage = new DashboardPage();
  const servicesPage = new ServicesPage();
  const sequencePage = new SequencesPage();
  const basePage = new BasePage();
  const mockedKeptnContext = '62cca6f3-dc54-4df6-a04c-6ffc894a4b5e';
  const mockedProject = 'sockshop';
  const mockedService = 'carts';
  const mockedServiceDeploymentContext = 'da740469-9920-4e0c-b304-0fd4b18d17c2';

  beforeEach(() => {
    environmentPage.intercept();
    dashboardPage.intercept();
    servicesPage.intercept();
    sequencePage.intercept();
    projectBoardPage.interceptDeepLinks();
  });

  it('should show environment screen', () => {
    environmentPage.visit(mockedProject);
    cy.location('pathname').should('eq', `/project/${mockedProject}`);
    projectBoardPage.assertOnlyEnvironmentViewSelected();
    dashboardPage.visit();

    cy.location('pathname').should('eq', '/dashboard');
  });

  it('should navigate to dashboard through navigate to root', () => {
    environmentPage.visit(mockedProject);
    cy.location('pathname').should('eq', `/project/${mockedProject}`);
    dashboardPage.visit();

    cy.location('pathname').should('eq', '/dashboard');
  });

  it('should navigate to dashboard through click on header icon', () => {
    environmentPage.visit(mockedProject);
    cy.location('pathname').should('eq', `/project/${mockedProject}`);
    basePage.clickMainHeaderKeptn();
    dashboardPage.waitForProjects();

    cy.location('pathname').should('eq', '/dashboard');
  });

  it('deepLink project/:projectName/service', () => {
    servicesPage.visitServicePage(mockedProject);
    cy.location('pathname').should('eq', `/project/${mockedProject}/service`);

    projectBoardPage.assertOnlyServicesViewSelected();
  });

  it('deepLink project/:projectName/service/:serviceName', () => {
    servicesPage.visitService(mockedProject, mockedService);

    cy.location('pathname').should('eq', `/project/${mockedProject}/service/${mockedService}`);

    projectBoardPage.assertOnlyServicesViewSelected();
    servicesPage
      .assertServiceExpanded(mockedService, true)
      .assertDeploymentSelected(mockedService, 'v0.1.2', false)
      .assertDeploymentSelected(mockedService, 'v0.1.1', false);
  });

  it('deepLink project/:projectName/service/:serviceName/context/:shkeptncontext', () => {
    const stage = 'production';
    servicesPage.visitServiceDeployment(mockedProject, mockedService, mockedServiceDeploymentContext);

    cy.location('pathname').should(
      'eq',
      `/project/${mockedProject}/service/${mockedService}/context/${mockedServiceDeploymentContext}/stage/${stage}`
    );

    projectBoardPage.assertOnlyServicesViewSelected();

    servicesPage
      .assertServiceExpanded(mockedService, true)
      .assertDeploymentSelected(mockedService, 'v0.1.2', true)
      .assertDeploymentSelected(mockedService, 'v0.1.1', false)
      .assertStageSelected(stage, true)
      .assertStageSelected('staging', false);
  });

  it('deepLink project/:projectName/service/:serviceName/context/:shkeptncontext/stage/:stage', () => {
    const stage = 'staging';
    servicesPage.visitServiceDeployment(mockedProject, mockedService, mockedServiceDeploymentContext, stage);

    cy.location('pathname').should(
      'eq',
      `/project/${mockedProject}/service/${mockedService}/context/${mockedServiceDeploymentContext}/stage/${stage}`
    );

    projectBoardPage.assertOnlyServicesViewSelected();
    servicesPage
      .assertServiceExpanded(mockedService, true)
      .assertDeploymentSelected(mockedService, 'v0.1.2', true)
      .assertDeploymentSelected(mockedService, 'v0.1.1', false)
      .assertStageSelected(stage, true)
      .assertStageSelected('production', false);
  });

  it('deepLink project/:projectName/sequence', () => {
    sequencePage.visit(mockedProject);

    cy.location('pathname').should('eq', `/project/${mockedProject}/sequence`);

    projectBoardPage.assertOnlySequencesViewSelected();
  });

  it('deepLink project/:projectName/sequence/:shkeptncontext', () => {
    sequencePage.visitContext(mockedProject, mockedKeptnContext);

    cy.location('pathname').should('eq', `/project/${mockedProject}/sequence/${mockedKeptnContext}/stage/production`);

    projectBoardPage.assertOnlySequencesViewSelected();
    sequencePage
      .assertTimelineStageSelected('dev', false)
      .assertTimelineStageSelected('staging', false)
      .assertTimelineStageSelected('production', true);
  });

  it('deepLink project/:projectName/sequence/:shkeptncontext/stage/:stage', () => {
    const stage = 'staging';
    sequencePage.visitContext(mockedProject, mockedKeptnContext, stage);

    cy.location('pathname').should('eq', `/project/${mockedProject}/sequence/${mockedKeptnContext}/stage/${stage}`);

    projectBoardPage.assertOnlySequencesViewSelected();

    sequencePage
      .assertTimelineStageSelected('dev', false)
      .assertTimelineStageSelected(stage, true)
      .assertTimelineStageSelected('production', false);
  });

  it('deepLink project/:projectName/sequence/:shkeptncontext/event/:eventId', () => {
    const eventId = 'ad13f4f6-2ec2-4e40-95db-ef325eed02d9';
    sequencePage.visitEvent(mockedProject, mockedKeptnContext, eventId);

    cy.location('pathname').should('eq', `/project/${mockedProject}/sequence/${mockedKeptnContext}/event/${eventId}`);

    projectBoardPage.assertOnlySequencesViewSelected();
    sequencePage
      .assertTimelineStageSelected('dev', false)
      .assertTimelineStageSelected('staging', true)
      .assertTimelineStageSelected('production', false)
      .assertTaskExpanded(eventId, true);
  });

  it('deepLink project/:projectName/sequence/:shkeptncontext/event/:eventId with sequence that is not initially loaded', () => {
    const eventId = 'ffd870da-bca7-49a1-bafd-726c234bfd3b';
    const keptnContext = '1663de8a-a414-47ba-9566-10a9730f40ff';
    sequencePage.interceptSequencesPageWithSequenceThatIsNotLoaded().visitEvent(mockedProject, keptnContext, eventId);

    cy.wait('@sequenceTraces')
      .location('pathname')
      .should('eq', `/project/${mockedProject}/sequence/${keptnContext}/event/${eventId}`);

    sequencePage
      .assertTimelineStageSelected('dev', true)
      .assertTimelineStageSelected('staging', false)
      .assertTaskExpanded(eventId, true);
  });

  it('deepLink trace/:shkeptncontext', () => {
    sequencePage.visitByContext(mockedKeptnContext);

    cy.location('pathname').should('eq', `/project/${mockedProject}/sequence/${mockedKeptnContext}/stage/production`);

    projectBoardPage.assertOnlySequencesViewSelected();
    sequencePage.assertTimelineStageSelected('production', true);
  });

  it('deepLink trace/:shkeptncontext/:stage', () => {
    const stage = 'staging';
    sequencePage.visitByContext(mockedKeptnContext, stage);

    cy.location('pathname').should('eq', `/project/${mockedProject}/sequence/${mockedKeptnContext}/stage/${stage}`);

    projectBoardPage.assertOnlySequencesViewSelected();
    sequencePage
      .assertTimelineStageSelected(stage, true)
      .assertTimelineStageSelected('dev', false)
      .assertTimelineStageSelected('production', false);
  });

  it('deepLink trace/:shkeptncontext/:eventType', () => {
    // eventType is actually stage or eventType
    const eventId = 'ffd870da-bca7-49a1-bafd-726c234bfd3b';
    sequencePage.visitByEventType(mockedKeptnContext, EventTypes.DEPLOYMENT_TRIGGERED);

    cy.location('pathname').should('eq', `/project/${mockedProject}/sequence/${mockedKeptnContext}/event/${eventId}`);

    projectBoardPage.assertOnlySequencesViewSelected();
    sequencePage
      .assertTimelineStageSelected('dev', true)
      .assertTimelineStageSelected('staging', false)
      .assertTimelineStageSelected('production', false)
      .assertTaskExpanded(eventId, true)
      .assertTaskExpanded('74fae034-1a4f-46eb-80d6-45bf640845f4', false); // just to check if this function does not always assert to true
  });
});
