import { ResultTypes } from 'shared/models/result-types';
import { interceptEvaluationBoard, interceptEvaluationBoardWithoutDeployment } from '../intercept';

export class EvaluationBoardPage {
  public visit(keptnContext: string, stage?: string): this {
    const stagePath = stage ? `/${stage}` : '';
    cy.visit(`/evaluation/${keptnContext}${stagePath}`).wait('@projects');
    return this;
  }

  public intercept(): this {
    interceptEvaluationBoard();
    return this;
  }

  public interceptWithoutDeployment(): this {
    interceptEvaluationBoardWithoutDeployment();
    return this;
  }

  public assertLoadingExists(status: boolean): this {
    cy.byTestId('ktb-loading').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertDefaultErrorVisible(status: boolean): this {
    cy.byTestId('ktb-trace-default-error').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertTraceErrorVisible(status: boolean): this {
    cy.byTestId('ktb-trace-error').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertKeptnContext(keptnContext: string): this {
    cy.byTestId('ktb-keptn-context').should('have.text', keptnContext);
    return this;
  }

  public assertArtifact(artifact: string): this {
    cy.byTestId('ktb-artifact').should('have.text', artifact);
    return this;
  }

  public assertArtifactExists(status: boolean): this {
    cy.byTestId('ktb-artifact').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertStage(stageName: string): this {
    cy.byTestId('ktb-stage').should('have.text', stageName);
    return this;
  }

  public assertViewServiceDetailsExists(status: boolean): this {
    cy.byTestId('ktb-view-service-details').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertViewSequenceDetailsExists(status: boolean): this {
    cy.byTestId('ktb-view-sequence-details').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertScoreInfo(score: number, equality: '<' | '<=' | '>' | '>=', threshold: number): this {
    cy.byTestId('keptn-evaluation-details-scoreInfo').should('have.text', `${score} ${equality} ${threshold}`);
    return this;
  }
  public assertResultInfo(type: ResultTypes): this {
    cy.byTestId('keptn-evaluation-details-resultInfo').should('have.text', `Result: ${type}`);
    return this;
  }
  public assertKeySliInfo(type: 'passed' | 'failed'): this {
    cy.byTestId('keptn-evaluation-details-keySliInfo').should('have.text', `Key SLI: ${type}`);
    return this;
  }
}
