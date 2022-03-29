export class SequencesPage {
  public visit(projectName: string): this {
    cy.visit(`/project/${projectName}/sequence`);
    return this;
  }

  public selectSequence(keptnContext: string): this {
    cy.byTestId(`keptn-root-events-list-${keptnContext}`).click();
    return this;
  }

  public assertTaskState(taskName: string, isFailed: boolean, isSuccess: boolean): this {
    cy.byTestId(`keptn-task-item-${taskName}`)
      .find('ktb-expandable-tile')
      .eq(0)
      .should(isFailed ? 'have.class' : 'not.have.class', 'ktb-tile-error');
    cy.byTestId(`keptn-task-item-${taskName}`)
      .find('ktb-expandable-tile')
      .eq(0)
      .should(isSuccess ? 'have.class' : 'not.have.class', 'ktb-tile-success');
    return this;
  }

  public assertTimelineTime(stage: string, time: string): this {
    cy.get('.stage-info')
      .contains(stage)
      .parentsUntilTestId(`keptn-sequence-timeline-stage-${stage}`)
      .should('contain.text', time);
    return this;
  }

  public assertTimelineTimeLoading(stage: string, exists: boolean): this {
    cy.get('.stage-info')
      .contains(stage)
      .parentsUntilTestId(`keptn-sequence-timeline-stage-${stage}`)
      .find('dt-loading-spinner')
      .should(exists ? 'exist' : 'not.exist');
    return this;
  }
}
