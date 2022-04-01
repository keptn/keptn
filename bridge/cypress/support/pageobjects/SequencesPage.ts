export class SequencesPage {
  public visit(projectName: string): this {
    cy.visit(`/project/${projectName}/sequence`);
    return this;
  }

  public selectSequence(keptnContext: string): this {
    cy.byTestId(`keptn-root-events-list-${keptnContext}`).click();
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

  public assertServiceName(name: string, tag?: string) {
    const serviceName = tag ? `${name}:${tag}` : name;
    cy.get('.service-name').should('have.text', serviceName);
  }
}
