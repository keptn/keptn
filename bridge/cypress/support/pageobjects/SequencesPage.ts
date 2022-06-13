import { EventTypes } from '../../../shared/interfaces/event-types';
import { interceptMain, interceptSequencesPage, interceptSequencesPageWithSequenceThatIsNotLoaded } from '../intercept';

export class SequencesPage {
  private readonly sequenceWaitingMessage = ' Sequence is waiting for previous sequences to finish. ';

  public intercept(): this {
    interceptSequencesPage();
    return this;
  }

  public interceptSequencesPageWithSequenceThatIsNotLoaded(): this {
    interceptSequencesPageWithSequenceThatIsNotLoaded();
    return this;
  }

  public visit(projectName: string, queryParams?: { [p: string]: string | string[] }): this {
    cy.visit({
      url: `/project/${projectName}/sequence`,
      qs: queryParams,
    })
      .wait('@metadata')
      .wait('@SequencesMetadata');
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
    }).as('project');
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
    cy.intercept('/api/project/sockshop/sequences/filter', { fixture: 'sequence.filter.mock' }).as('SequencesMetadata');

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

  private setFilterForGroup(filterGroup: string, itemName: string, status: boolean): this {
    cy.byTestId('keptn-sequence-view-filter').find('dt-quick-filter').dtQuickFilterCheck(filterGroup, itemName, status);
    return this;
  }

  public assertFilterIsChecked(filterGroup: string, itemName: string, status: boolean): this {
    cy.byTestId('keptn-sequence-view-filter')
      .find('dt-quick-filter')
      .dtQuickFilterIsChecked(filterGroup, itemName, status);
    return this;
  }

  public checkServiceFilter(serviceName: string, status = true): this {
    return this.setFilterForGroup('Service', serviceName, status);
  }

  public checkStageFilter(stageName: string, status = true): this {
    return this.setFilterForGroup('Stage', stageName, status);
  }

  public checkSequenceFilter(sequenceName: string, status = true): this {
    return this.setFilterForGroup('Sequence', sequenceName, status);
  }

  public checkStatusFilter(statusName: string, status = true): this {
    return this.setFilterForGroup('Status', statusName, status);
  }

  public clearFilter(): this {
    cy.byTestId('keptn-sequence-view-filter').find('dt-quick-filter').clearDtFilter();
    return this;
  }

  public clickLoadOlderSequences(): this {
    cy.byTestId('keptn-show-older-sequences-button').click();
    return this;
  }

  public reload(): this {
    cy.reload();
    return this;
  }

  public assertLoadOlderSequencesButtonExists(exists: boolean): this {
    cy.byTestId('keptn-show-older-sequences-button').should(exists ? 'exist' : 'not.exist');
    return this;
  }

  public assertSequenceCount(count: number): this {
    if (count === 0) {
      cy.byTestId('keptn-sequence-view-roots').should('not.exist');
    } else {
      cy.byTestId('keptn-sequence-view-roots').get('ktb-selectable-tile').should('have.length', count);
    }
    return this;
  }

  public assertServiceNameOfSequences(serviceName: string): this {
    return this.assertSequenceTile('keptn-sequence-info-serviceName', serviceName);
  }

  public assertStageNameOfSequences(stageName: string): this {
    return this.assertStageNamesOfSequences([stageName], false);
  }

  public assertStageNamesOfSequences(stageNames: string[], validateLength = true): this {
    cy.byTestId('keptn-sequence-info-stageDetails').each((el) => {
      for (const stageName of stageNames) {
        cy.wrap(el).find('ktb-stage-badge').contains(stageName).should('exist');
      }
      if (validateLength) {
        cy.wrap(el).find('ktb-stage-badge').should('have.length', stageNames.length);
      }
    });
    return this;
  }

  public assertSequenceNameOfSequences(sequenceName: string): this {
    return this.assertSequenceTile('keptn-sequence-info-sequenceName', sequenceName);
  }

  public assertStatusOfSequences(status: string): this {
    return this.assertSequenceTile('keptn-sequence-info-status', status);
  }

  private assertSequenceTile(testId: string, expectedText: string): this {
    cy.byTestId(testId).each((el) => {
      cy.wrap(el).should('have.text', expectedText);
    });
    return this;
  }

  public assertNoSequencesMessageExists(status: boolean): this {
    cy.byTestId('keptn-noSequencesFiltered').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertNoSequencesFilteredMessageExists(status: boolean): this {
    cy.byTestId('keptn-noSequencesFiltered').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertLoadingOldSequencesButtonExists(status: boolean): this {
    cy.byTestId('keptn-loadingOldSequences').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertIsLoadingSequences(status: boolean): this {
    cy.byTestId('keptn-loadingSequences').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertIsWaitingSequence(keptnContext: string, status: boolean): this {
    cy.byTestId(`keptn-root-events-list-${keptnContext}`)
      .find('.ktb-selectable-tile-content')
      .should(status ? 'have.text' : 'not.have.text', this.sequenceWaitingMessage);
    return this;
  }

  public assertIsSelectedSequenceWaiting(status: boolean): this {
    cy.byTestId('keptn-sequence-view-sequenceDetails')
      .find('dt-alert')
      .should(status ? 'have.text' : 'not.have.text', this.sequenceWaitingMessage);
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
      .find('ktb-loading-spinner')
      .should(exists ? 'exist' : 'not.exist');
    return this;
  }

  public assertServiceName(name: string, tag?: string): this {
    const serviceName = tag ? `${name}:${tag}` : name;
    cy.byTestId('keptn-sequence-view-serviceName').should('have.text', serviceName);
    return this;
  }

  public assertQueryParams(query: string): this {
    cy.location('search').should('eq', query);
    return this;
  }

  public clickEvent(id: string): this {
    cy.byTestId(`ktb-task-${id}`).click();
    return this;
  }

  public assertIsApprovalLoading(status: boolean): this {
    cy.byTestId('ktb-approval-loading').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertApprovalDeployedImage(image: string): this {
    cy.byTestId('ktb-approval-deployed-image').should('have.text', image);
    return this;
  }

  public assertAmountOfQueryParameters(amount: number): this {
    // eslint-disable-next-line promise/catch-or-return
    cy.url()
      .then((url) => {
        const splitUrl = url.split('?');
        return splitUrl.length <= 1 ? [] : splitUrl[1].split('&');
      })
      .should('have.length', amount);
    return this;
  }
}
