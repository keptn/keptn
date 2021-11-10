import { Deployment, StageDeployment } from './deployment';
import {
  ServiceDeploymentMock,
  ServiceDeploymentWithApprovalMock,
} from '../../../shared/fixtures/service-deployment-response.mock';
import { OpenRemediationsResponse } from '../../../server/fixtures/open-remediations-response.mock';
import { ResultTypes } from '../../../shared/models/result-types';
import { SequenceState } from '../../../shared/models/sequence';
import { ServiceRemediationInformation } from './service-remediation-information';
import { Sequence } from './sequence';

describe('Deployment', () => {
  it('should correctly create new class', () => {
    const deployment = Deployment.fromJSON(ServiceDeploymentWithApprovalMock);
    expect(deployment).toBeInstanceOf(Deployment);
    expect(deployment.stages[0]).toBeInstanceOf(StageDeployment);
    expect(deployment.stages[1]).toBeInstanceOf(StageDeployment);
  });

  it('should correctly update', () => {
    const deployment = Deployment.fromJSON(ServiceDeploymentWithApprovalMock);
    const newDeployment = Deployment.fromJSON({
      state: 'finished',
      keptnContext: '2c0e568b-8bd3-4726-a188-e528423813ed',
      service: 'carts',
      stages: [
        {
          name: 'dev',
          lastTimeUpdated: '2021-10-13T10:49:30.005Z',
          openRemediations: [],
          subSequences: [],
          hasEvaluation: true,
          latestEvaluation: {
            traces: [],
            data: {
              evaluation: {
                gitCommit: '',
                indicatorResults: null,
                result: 'pass',
                score: 0,
                sloFileContent: '',
                timeEnd: '2021-10-13T10:47:11Z',
                timeStart: '2021-10-13T10:46:42Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              message: 'no evaluation performed by lighthouse because no SLO file configured for project sockshop',
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'dev',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: 'b286a8ea-1b11-4106-b892-fade06245bcf',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-13T10:49:27.606Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '2c0e568b-8bd3-4726-a188-e528423813ed',
            shkeptnspecversion: '0.2.3',
            triggeredid: '5a548dc5-b29f-4f2d-a2be-4d4534de523f',
          },
        },
        {
          name: 'staging',
          lastTimeUpdated: '2021-10-13T10:54:43.315Z',
          openRemediations: [OpenRemediationsResponse.states[0]],
          subSequences: [
            {
              name: 'delivery',
              type: 'sh.keptn.event.staging.delivery.triggered',
              result: 'pass',
              time: '2021-10-13T10:49:30.202Z',
              state: 'finished',
              id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
              message: 'Finished release',
              hasPendingApproval: false,
            },
          ],
          hasEvaluation: true,
          latestEvaluation: {
            data: {
              evaluation: {
                comparedEvents: ['3344487d-e384-4cd9-a0e0-fcf157a33ad6'],
                gitCommit: '',
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 336.9946150194969,
                        violated: false,
                      },
                      {
                        criteria: '<600',
                        targetValue: 600,
                        violated: false,
                      },
                    ],
                    score: 1,
                    status: 'pass',
                    value: {
                      metric: 'response_time_p95',
                      success: true,
                      value: 304.2952915485157,
                    },
                    warningTargets: [
                      {
                        criteria: '<=800',
                        targetValue: 800,
                        violated: false,
                      },
                    ],
                  },
                ],
                result: 'pass',
                score: 100,
                sloFileContent:
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-13T10:53:29Z',
                timeStart: '2021-10-13T10:51:08Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '93c2eba9-b77c-4976-b079-29a0188d86eg',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-13T10:54:43.112Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '2c0e568b-8bd3-4726-a188-e528423813ed',
            shkeptnspecversion: '0.2.3',
            triggeredid: '1cc9c272-721a-43de-98f6-9eceae484cf5',
          },
        },
        {
          name: 'production',
          lastTimeUpdated: '2021-10-13T10:54:43.315Z',
          openRemediations: [],
          subSequences: [
            {
              name: 'delivery',
              type: 'sh.keptn.event.staging.delivery.triggered',
              result: 'pass',
              time: '2021-10-13T10:49:30.202Z',
              state: 'finished',
              id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
              message: '',
              hasPendingApproval: false,
            },
          ],
          hasEvaluation: false,
        },
      ],
      labels: {
        DtCreds: 'dynatrace',
      },
      image: 'carts:0.12.3',
    });
    const expectedDeployment = Deployment.fromJSON({
      state: 'finished',
      keptnContext: '2c0e568b-8bd3-4726-a188-e528423813ed',
      service: 'carts',
      stages: [
        {
          name: 'dev',
          lastTimeUpdated: '2021-10-13T10:49:30.005Z',
          openRemediations: [],
          subSequences: [
            {
              name: 'delivery',
              type: 'sh.keptn.event.dev.delivery.triggered',
              result: 'pass',
              time: '2021-10-13T10:45:03.780Z',
              state: 'finished',
              id: '08e89fdb-02db-4fc7-a5fd-386d03e1c4a9',
              message: 'Finished release',
              hasPendingApproval: false,
            },
          ],
          hasEvaluation: true,
          latestEvaluation: {
            traces: [],
            data: {
              evaluation: {
                gitCommit: '',
                indicatorResults: null,
                result: 'pass',
                score: 0,
                sloFileContent: '',
                timeEnd: '2021-10-13T10:47:11Z',
                timeStart: '2021-10-13T10:46:42Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              message: 'no evaluation performed by lighthouse because no SLO file configured for project sockshop',
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'dev',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: 'b286a8ea-1b11-4106-b892-fade06245bcf',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-13T10:49:27.606Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '2c0e568b-8bd3-4726-a188-e528423813ed',
            shkeptnspecversion: '0.2.3',
            triggeredid: '5a548dc5-b29f-4f2d-a2be-4d4534de523f',
          },
        },
        {
          name: 'staging',
          lastTimeUpdated: '2021-10-13T10:54:43.315Z',
          openRemediations: [OpenRemediationsResponse.states[0]],
          subSequences: [
            {
              name: 'delivery',
              type: 'sh.keptn.event.staging.delivery.triggered',
              result: 'pass',
              time: '2021-10-13T10:49:30.202Z',
              state: 'finished',
              id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
              message: 'Finished release',
              hasPendingApproval: false,
            },
          ],
          hasEvaluation: true,
          latestEvaluation: {
            data: {
              evaluation: {
                comparedEvents: ['3344487d-e384-4cd9-a0e0-fcf157a33ad6'],
                gitCommit: '',
                indicatorResults: [
                  {
                    displayName: 'Response time P95',
                    keySli: false,
                    passTargets: [
                      {
                        criteria: '<=+10%',
                        targetValue: 336.9946150194969,
                        violated: false,
                      },
                      {
                        criteria: '<600',
                        targetValue: 600,
                        violated: false,
                      },
                    ],
                    score: 1,
                    status: 'pass',
                    value: {
                      metric: 'response_time_p95',
                      success: true,
                      value: 304.2952915485157,
                    },
                    warningTargets: [
                      {
                        criteria: '<=800',
                        targetValue: 800,
                        violated: false,
                      },
                    ],
                  },
                ],
                result: 'pass',
                score: 100,
                sloFileContent:
                  'LS0tCnNwZWNfdmVyc2lvbjogIjEuMCIKY29tcGFyaXNvbjoKICBhZ2dyZWdhdGVfZnVuY3Rpb246ICJhdmciCiAgY29tcGFyZV93aXRoOiAic2luZ2xlX3Jlc3VsdCIKICBpbmNsdWRlX3Jlc3VsdF93aXRoX3Njb3JlOiAicGFzcyIKICBudW1iZXJfb2ZfY29tcGFyaXNvbl9yZXN1bHRzOiAxCmZpbHRlcjoKb2JqZWN0aXZlczoKICAtIHNsaTogInJlc3BvbnNlX3RpbWVfcDk1IgogICAgZGlzcGxheU5hbWU6ICJSZXNwb25zZSB0aW1lIFA5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSAxMCUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNjAwbXMpCiAgICAgIC0gY3JpdGVyaWE6CiAgICAgICAgICAtICI8PSsxMCUiICAjIHJlbGF0aXZlIHZhbHVlcyByZXF1aXJlIGEgcHJlZml4ZWQgc2lnbiAocGx1cyBvciBtaW51cykKICAgICAgICAgIC0gIjw2MDAiICAgICMgYWJzb2x1dGUgdmFsdWVzIG9ubHkgcmVxdWlyZSBhIGxvZ2ljYWwgb3BlcmF0b3IKICAgIHdhcm5pbmc6ICAgICAgICAgICMgaWYgdGhlIHJlc3BvbnNlIHRpbWUgaXMgYmVsb3cgODAwbXMsIHRoZSByZXN1bHQgc2hvdWxkIGJlIGEgd2FybmluZwogICAgICAtIGNyaXRlcmlhOgogICAgICAgICAgLSAiPD04MDAiCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==',
                timeEnd: '2021-10-13T10:53:29Z',
                timeStart: '2021-10-13T10:51:08Z',
              },
              labels: {
                DtCreds: 'dynatrace',
              },
              project: 'sockshop',
              result: 'pass',
              service: 'carts',
              stage: 'staging',
              status: 'succeeded',
              temporaryData: {
                distributor: {
                  subscriptionID: '',
                },
              },
            },
            id: '93c2eba9-b77c-4976-b079-29a0188d86eg',
            source: 'lighthouse-service',
            specversion: '1.0',
            time: '2021-10-13T10:54:43.112Z',
            type: 'sh.keptn.event.evaluation.finished',
            shkeptncontext: '2c0e568b-8bd3-4726-a188-e528423813ed',
            shkeptnspecversion: '0.2.3',
            triggeredid: '1cc9c272-721a-43de-98f6-9eceae484cf5',
          },
        },
        {
          name: 'production',
          lastTimeUpdated: '2021-10-13T10:54:43.315Z',
          openRemediations: [],
          subSequences: [
            {
              name: 'delivery',
              type: 'sh.keptn.event.staging.delivery.triggered',
              result: 'pass',
              time: '2021-10-13T10:49:30.202Z',
              state: 'finished',
              id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
              message: '',
              hasPendingApproval: false,
            },
          ],
          hasEvaluation: false,
        },
      ],
      labels: {
        DtCreds: 'dynatrace',
      },
      image: 'carts:0.12.3',
    });
    deployment.update(newDeployment);

    expect(deployment).toEqual(expectedDeployment);
  });

  it('should assign subSequences', () => {
    const stageDeployment = StageDeployment.fromJSON({
      name: 'dev',
      lastTimeUpdated: '2021-10-13T10:49:30.005Z',
      openRemediations: [],
      subSequences: [],
      hasEvaluation: false,
    });
    const newStageDeployment = StageDeployment.fromJSON({
      name: 'dev',
      lastTimeUpdated: '2021-10-13T10:49:30.005Z',
      openRemediations: [],
      subSequences: [
        {
          name: 'delivery',
          type: 'sh.keptn.event.dev.delivery.triggered',
          result: ResultTypes.PASSED,
          time: '2021-10-13T10:45:03.780Z',
          state: SequenceState.FINISHED,
          id: '08e89fdb-02db-4fc7-a5fd-386d03e1c4a9',
          message: 'Finished release',
          hasPendingApproval: false,
        },
      ],
      hasEvaluation: false,
    });
    stageDeployment.update(newStageDeployment);
    expect(stageDeployment.subSequences).toEqual(newStageDeployment.subSequences);
  });

  it('should add subSequences', () => {
    const stageDeployment = StageDeployment.fromJSON({
      name: 'dev',
      lastTimeUpdated: '2021-10-13T10:49:30.005Z',
      openRemediations: [],
      subSequences: [
        {
          name: 'delivery',
          type: 'sh.keptn.event.dev.delivery.triggered',
          result: ResultTypes.PASSED,
          time: '2021-10-13T10:45:03.780Z',
          state: SequenceState.FINISHED,
          id: '08e89fdb-02db-4fc7-a5fd-386d03e1c4a9',
          message: 'Finished release',
          hasPendingApproval: false,
        },
      ],
      hasEvaluation: false,
    });
    const newStageDeployment = StageDeployment.fromJSON({
      name: 'dev',
      lastTimeUpdated: '2021-10-13T10:49:30.005Z',
      openRemediations: [],
      subSequences: [
        {
          name: 'delivery',
          type: 'sh.keptn.event.dev.rollback.triggered',
          result: ResultTypes.PASSED,
          time: '2021-10-13T11:45:03.780Z',
          state: SequenceState.FINISHED,
          id: '08e89fdb-02db-4fc7-a5fd-386d03e1c4a1',
          message: 'Finished rollback',
          hasPendingApproval: false,
        },
      ],
      hasEvaluation: false,
    });
    const expectedSubSequences = [
      {
        name: 'delivery',
        type: 'sh.keptn.event.dev.rollback.triggered',
        result: ResultTypes.PASSED,
        time: '2021-10-13T11:45:03.780Z',
        state: SequenceState.FINISHED,
        id: '08e89fdb-02db-4fc7-a5fd-386d03e1c4a1',
        message: 'Finished rollback',
        hasPendingApproval: false,
      },
      {
        name: 'delivery',
        type: 'sh.keptn.event.dev.delivery.triggered',
        result: ResultTypes.PASSED,
        time: '2021-10-13T10:45:03.780Z',
        state: SequenceState.FINISHED,
        id: '08e89fdb-02db-4fc7-a5fd-386d03e1c4a9',
        message: 'Finished release',
        hasPendingApproval: false,
      },
    ];
    stageDeployment.update(newStageDeployment);
    expect(stageDeployment.subSequences).toEqual(expectedSubSequences);
  });

  it('should return latest time updated', () => {
    const deployment = Deployment.fromJSON(ServiceDeploymentWithApprovalMock);
    expect(deployment.latestTimeUpdated).toEqual(new Date('2021-10-13T10:54:43.315Z'));
  });

  it('should remove open remediations', () => {
    const deployment = Deployment.fromJSON(ServiceDeploymentMock);
    const serviceRemediationInformation = ServiceRemediationInformation.fromJSON({
      stages: [
        {
          name: 'dev',
          config: 'configString',
          remediations: [
            Sequence.fromJSON({
              name: 'remediation',
              problemTitle: 'Failure rate increase',
              project: 'sockshop',
              service: 'carts',
              shkeptncontext: '35383737-3630-4639-b037-353138323631',
              stages: [
                {
                  actions: [],
                  latestEvent: {
                    id: '7448420f-5b15-4777-9d39-cc8308e2b0c3',
                    time: '2021-11-04T04:51:23.266Z',
                    type: 'sh.keptn.event.get-action.started',
                  },
                  name: 'dev',
                  state: SequenceState.STARTED,
                },
              ],
              state: SequenceState.STARTED,
              time: '2021-11-04T04:51:21.557Z',
            }),
          ],
        },
      ],
    });
    deployment.updateRemediations(serviceRemediationInformation);

    // only update deployed
    expect(deployment.stages[0].remediationConfig).toBeUndefined();
    expect(deployment.stages[0].openRemediations).toEqual([]);

    expect(deployment.stages[1].remediationConfig).toBeUndefined();
    expect(deployment.stages[1].openRemediations).toEqual([]);

    expect(deployment.stages[2].remediationConfig).toBeUndefined();
    expect(deployment.stages[2].openRemediations).toEqual([]);
  });

  it('should update open remediations', () => {
    const deployment = Deployment.fromJSON(ServiceDeploymentMock);
    deployment.stages[2].remediationConfig = undefined;
    deployment.stages[2].openRemediations = [];
    const serviceRemediationInformation = ServiceRemediationInformation.fromJSON({
      stages: [
        {
          name: 'production',
          config: 'configString',
          remediations: [
            Sequence.fromJSON({
              name: 'remediation',
              problemTitle: 'Failure rate increase',
              project: 'sockshop',
              service: 'carts',
              shkeptncontext: '35383737-3630-4639-b037-353138323631',
              stages: [
                {
                  actions: [],
                  latestEvent: {
                    id: '7448420f-5b15-4777-9d39-cc8308e2b0c3',
                    time: '2021-11-04T04:51:23.266Z',
                    type: 'sh.keptn.event.get-action.started',
                  },
                  name: 'production',
                  state: SequenceState.STARTED,
                },
              ],
              state: SequenceState.STARTED,
              time: '2021-11-04T04:51:21.557Z',
            }),
          ],
        },
      ],
    });
    deployment.updateRemediations(serviceRemediationInformation);

    // only update deployed
    expect(deployment.stages[0].remediationConfig).toBeUndefined();
    expect(deployment.stages[0].openRemediations).toEqual([]);

    expect(deployment.stages[1].remediationConfig).toBeUndefined();
    expect(deployment.stages[1].openRemediations).toEqual([]);

    expect(deployment.stages[2].remediationConfig).toEqual(serviceRemediationInformation.stages[0].config);
    expect(deployment.stages[2].openRemediations).toEqual(serviceRemediationInformation.stages[0].remediations);
  });

  it('should remove approval', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.removeApproval();
    expect(stageDeployment.subSequences[0].hasPendingApproval).toBe(false);
    expect(stageDeployment.approvalInformation).toBeUndefined();
  });

  it('should be faulty', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = [
      {
        name: 'delivery',
        type: 'sh.keptn.event.staging.delivery.triggered',
        result: ResultTypes.PASSED,
        time: '2021-10-13T10:49:30.202Z',
        state: SequenceState.FINISHED,
        id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
        message: '',
        hasPendingApproval: false,
      },
      {
        name: 'delivery',
        type: 'sh.keptn.event.staging.delivery.triggered',
        result: ResultTypes.FAILED,
        time: '2021-10-13T10:49:30.202Z',
        state: SequenceState.FINISHED,
        id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
        message: '',
        hasPendingApproval: false,
      },
    ];
    expect(stageDeployment.isFaulty()).toBe(true);
  });

  it('should not be faulty', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = [
      {
        name: 'delivery',
        type: 'sh.keptn.event.staging.delivery.triggered',
        result: ResultTypes.PASSED,
        time: '2021-10-13T10:49:30.202Z',
        state: SequenceState.FINISHED,
        id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
        message: '',
        hasPendingApproval: false,
      },
      {
        name: 'delivery',
        type: 'sh.keptn.event.staging.delivery.triggered',
        result: ResultTypes.PASSED,
        time: '2021-10-13T10:49:30.202Z',
        state: SequenceState.STARTED,
        id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
        message: '',
        hasPendingApproval: false,
      },
    ];
    expect(stageDeployment.isFaulty()).toBe(false);
  });

  it('should not be successful', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = [
      {
        name: 'delivery',
        type: 'sh.keptn.event.staging.delivery.triggered',
        result: ResultTypes.PASSED,
        time: '2021-10-13T10:49:30.202Z',
        state: SequenceState.FINISHED,
        id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
        message: '',
        hasPendingApproval: false,
      },
      {
        name: 'delivery',
        type: 'sh.keptn.event.staging.delivery.triggered',
        result: ResultTypes.PASSED,
        time: '2021-10-13T10:49:30.202Z',
        state: SequenceState.STARTED,
        id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
        message: '',
        hasPendingApproval: false,
      },
    ];
    expect(stageDeployment.isSuccessful()).toBe(false);
  });

  it('should be successful', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = [
      {
        name: 'delivery',
        type: 'sh.keptn.event.staging.delivery.triggered',
        result: ResultTypes.PASSED,
        time: '2021-10-13T10:49:30.202Z',
        state: SequenceState.FINISHED,
        id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
        message: '',
        hasPendingApproval: false,
      },
      {
        name: 'delivery',
        type: 'sh.keptn.event.staging.delivery.triggered',
        result: ResultTypes.PASSED,
        time: '2021-10-13T10:49:30.202Z',
        state: SequenceState.FINISHED,
        id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
        message: '',
        hasPendingApproval: false,
      },
    ];
    expect(stageDeployment.isSuccessful()).toBe(true);
  });

  it('should not be warning', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = [
      {
        name: 'delivery',
        type: 'sh.keptn.event.staging.delivery.triggered',
        result: ResultTypes.WARNING,
        time: '2021-10-13T10:49:30.202Z',
        state: SequenceState.FINISHED,
        id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
        message: '',
        hasPendingApproval: false,
      },
      {
        name: 'delivery',
        type: 'sh.keptn.event.staging.delivery.triggered',
        result: ResultTypes.FAILED,
        time: '2021-10-13T10:49:30.202Z',
        state: SequenceState.STARTED,
        id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
        message: '',
        hasPendingApproval: false,
      },
    ];
    expect(stageDeployment.isWarning()).toBe(false);
  });

  it('should be warning', () => {
    const stageDeployment = getStageDeployment();
    stageDeployment.subSequences = [
      {
        name: 'delivery',
        type: 'sh.keptn.event.staging.delivery.triggered',
        result: ResultTypes.WARNING,
        time: '2021-10-13T10:49:30.202Z',
        state: SequenceState.FINISHED,
        id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
        message: '',
        hasPendingApproval: false,
      },
      {
        name: 'delivery',
        type: 'sh.keptn.event.staging.delivery.triggered',
        result: ResultTypes.PASSED,
        time: '2021-10-13T10:49:30.202Z',
        state: SequenceState.STARTED,
        id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
        message: '',
        hasPendingApproval: false,
      },
    ];
    expect(stageDeployment.isWarning()).toBe(true);
  });

  function getStageDeployment(): StageDeployment {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    return StageDeployment.fromJSON(ServiceDeploymentWithApprovalMock.stages[1]);
  }
});
