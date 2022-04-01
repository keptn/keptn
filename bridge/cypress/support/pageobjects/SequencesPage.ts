export class SequencesPage {
  public visit(projectName: string): this {
    cy.visit(`/project/${projectName}/sequence`);
    return this;
  }

  public interceptRemediationSequences(): this {
    cy.intercept('/api/v1/metadata', { fixture: 'metadata.mock' });
    cy.intercept('/api/bridgeInfo', { fixture: 'bridgeInfo.mock' });
    cy.intercept('/api/project/sockshop?approval=true&remediation=true', {
      fixture: 'get.project.sockshop.remediation.mock',
    });
    cy.intercept('/api/hasUnreadUniformRegistrationLogs', { body: false });

    cy.intercept('/api/controlPlane/v1/project?disableUpstreamSync=true&pageSize=50', { fixture: 'projects.mock' });
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25', {
      fixture: 'get.sequences.remediation.mock',
    }).as('Sequences');
    cy.intercept('/api/controlPlane/v1/sequence/sockshop?pageSize=25&fromTime=*', {
      body: {
        states: [],
      },
    });
    cy.intercept('/api/project/sockshop/sequences/metadata', { fixture: 'sequence.metadata.mock' }).as(
      'SequencesMetadata'
    );

    cy.intercept('/api/mongodb-datastore/event?keptnContext=cfaadbb1-3c47-46e5-a230-2e312cf1828a&project=sockshop', {
      fixture: 'get.events.cfaadbb1-3c47-46e5-a230-2e312cf1828a.mock.json',
    });
    cy.intercept('/api/mongodb-datastore/event?keptnContext=cfaadbb1-3c47-46e5-a230-d0f055f4f518&project=sockshop', {
      fixture: 'get.events.cfaadbb1-3c47-46e5-a230-d0f055f4f518.mock.json',
    });
    cy.intercept('/api/mongodb-datastore/event?keptnContext=29355a07-7b65-47fa-896e-06f656283c5d&project=sockshop', {
      fixture: 'get.events.29355a07-7b65-47fa-896e-06f656283c5d.mock.json',
    });

    return this;
  }

  public selectSequence(keptnContext: string): this {
    cy.byTestId(`keptn-root-events-list-${keptnContext}`).click();
    return this;
  }

  public assertTaskFailed(taskName: string, isFailed: boolean): this {
    cy.byTestId(`keptn-task-item-${taskName}`)
      .find('ktb-expandable-tile')
      .eq(0)
      .should(isFailed ? 'have.class' : 'not.have.class', 'ktb-tile-error');
    return this;
  }

  public assertTaskSuccessful(taskName: string, isSuccess: boolean): this {
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
