import { EvaluationBoardPage } from '../support/pageobjects/EvaluationBoardPage';

describe('Evaluation Board', () => {
  const evaluationBoard = new EvaluationBoardPage();
  const keptnContext = '9b1929ac-8fda-476c-8cea-8135cb428484';

  beforeEach(() => {
    evaluationBoard.intercept();
  });

  it('should show loading', () => {
    cy.intercept(
      'api/mongodb-datastore/event?keptnContext=*&type=sh.keptn.event.evaluation.finished&source=lighthouse-service',
      {
        delay: 10_000,
      }
    );
    evaluationBoard.visit(keptnContext).assertLoadingExists(true);
  });

  it('should show trace error', () => {
    cy.intercept(
      'api/mongodb-datastore/event?keptnContext=*&type=sh.keptn.event.evaluation.finished&source=lighthouse-service',
      {
        statusCode: 403,
      }
    );
    evaluationBoard.visit(keptnContext).assertTraceErrorVisible(true);
  });

  it('should show default error', () => {
    cy.intercept('api/controlPlane/v1/project/dynatrace/stage/quality-gate/service/items', {
      statusCode: 403,
    });
    evaluationBoard.visit(keptnContext).assertDefaultErrorVisible(true);
  });

  it('should display the correct data', () => {
    evaluationBoard
      .visit(keptnContext)
      .assertKeptnContext(keptnContext)
      .assertStage('quality-gate')
      .assertArtifact('items:0.10.2')
      .assertViewServiceDetailsExists(true)
      .assertViewSequenceDetailsExists(true);
  });

  it('should not show artifact id', () => {
    evaluationBoard.interceptWithoutDeployment().visit(keptnContext).assertArtifactExists(false);
  });
});
