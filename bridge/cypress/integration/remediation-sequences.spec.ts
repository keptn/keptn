import { SequencesPage } from '../support/pageobjects/SequencesPage';

describe('Sequences', () => {
  const sequencePage = new SequencesPage();

  beforeEach(() => {
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
  });

  it('should show remediation in regular state while running', () => {
    sequencePage
      .visit('sockshop')
      .selectSequence('cfaadbb1-3c47-46e5-a230-2e312cf1828a')
      .assertTaskState('remediation', false, false);
  });

  it('should show remediation green when successful', () => {
    sequencePage
      .visit('sockshop')
      .selectSequence('cfaadbb1-3c47-46e5-a230-d0f055f4f518')
      .assertTaskState('remediation', false, true);
  });

  it('should show remediation red when failed', () => {
    sequencePage
      .visit('sockshop')
      .selectSequence('29355a07-7b65-47fa-896e-06f656283c5d')
      .assertTaskState('remediation', true, false);
  });
});
