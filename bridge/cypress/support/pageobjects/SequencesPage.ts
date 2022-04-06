import { interceptMain, interceptSequencesPage } from '../intercept';
import { EventTypes } from '../../../shared/interfaces/event-types';

export class SequencesPage {
  public intercept(): void {
    interceptSequencesPage();
  }

  public visit(projectName: string): this {
    cy.visit(`/project/${projectName}/sequence`).wait('@metadata');
    return this;
  }

  public visitContext(projectName: string, keptnContext: string, stage?: string): this {
    let url = `/project/${projectName}/sequence/${keptnContext}`;
    if (stage) {
      url += `/stage/${stage}`;
    }
    cy.visit(url).wait('@metadata');
    return this;
  }

  public visitEvent(projectName: string, keptnContext: string, eventId: string): this {
    cy.visit(`/project/${projectName}/sequence/${keptnContext}/event/${eventId}`).wait('@metadata');
    return this;
  }

  public visitByContext(keptnContext: string, stage?: string): this {
    let url = `/trace/${keptnContext}`;
    if (stage) {
      url += `/${stage}`;
    }
    cy.visit(url).wait('@metadata');
    return this;
  }

  public visitByEventType(keptnContext: string, eventType: EventTypes | string): this {
    cy.visit(`/trace/${keptnContext}/${eventType}`).wait('@metadata');
    return this;
  }

  public interceptRemediationSequences(): this {
    interceptMain();
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

  public assertTimelineStageSelected(stageName: string, status: boolean): this {
    cy.byTestId(`keptn-sequence-timeline-stage-${stageName}`)
      .find('.stage-text')
      .should(status ? 'have.class' : 'not.have.class', 'focused');
    return this;
  }

  public assertTaskExpanded(eventId: string, status: boolean): this {
    cy.byTestId(`ktb-task-${eventId}`)
      .find('.ktb-expandable-tile-content')
      .should(status ? 'be.visible' : 'not.be.visible');
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

  public assertServiceName(name: string, tag?: string): this {
    const serviceName = tag ? `${name}:${tag}` : name;
    cy.byTestId('keptn-sequence-view-serviceName').should('have.text', serviceName);
    return this;
  }
}
