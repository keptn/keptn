import { OpenRemediationsResponse } from '../../../../shared/fixtures/open-remediations-response.mock';
import { ResultTypes } from '../../../../shared/models/result-types';
import { SequenceState } from '../../../../shared/interfaces/sequence';
import { Sequence } from '../../_models/sequence';

const updatedDeploymentMock = {
  state: 'finished',
  keptnContext: '2c0e568b-8bd3-4726-a188-e528423813ed',
  service: 'carts',
  stages: [
    {
      name: 'dev',
      state: 'finished',
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
      state: 'finished',
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
      state: 'finished',
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
};

const expectedDeploymentMock = {
  state: 'finished',
  keptnContext: '2c0e568b-8bd3-4726-a188-e528423813ed',
  service: 'carts',
  stages: [
    {
      name: 'dev',
      state: 'finished',
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
      state: 'finished',
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
      state: 'finished',
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
};

const stageDeploymentDeliveryFinishedPass = {
  name: 'dev',
  state: SequenceState.FINISHED,
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
};

const stageDeploymentRollBackFinishedPass = {
  name: 'dev',
  state: SequenceState.FINISHED,
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
};

const serviceRemediationInformationDevWithRemediation = {
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
};

const serviceRemediationInformationProductionWithRemediation = {
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
};

const mergedSubsequencesDeliveryRollback = [
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

const defaultSubsequenceDelivery = {
  name: 'delivery',
  type: 'sh.keptn.event.staging.delivery.triggered',
  result: ResultTypes.PASSED,
  time: '2021-10-13T10:49:30.202Z',
  state: SequenceState.FINISHED,
  id: '3cbb5949-3852-4073-a70d-27ec52e04b93',
  message: '',
  hasPendingApproval: false,
};

const subSequencesFailedAndPassed = [
  defaultSubsequenceDelivery,
  {
    ...defaultSubsequenceDelivery,
    result: ResultTypes.FAILED,
    state: SequenceState.FINISHED,
  },
];

const subSequencesPassedLoading = [
  defaultSubsequenceDelivery,
  {
    ...defaultSubsequenceDelivery,
    result: ResultTypes.PASSED,
    state: SequenceState.STARTED,
  },
];

const subSequencesPassed = [defaultSubsequenceDelivery, defaultSubsequenceDelivery];

const subSequencesWarningFailed = [
  {
    ...defaultSubsequenceDelivery,
    result: ResultTypes.WARNING,
    state: SequenceState.FINISHED,
  },
  {
    ...defaultSubsequenceDelivery,
    result: ResultTypes.FAILED,
    state: SequenceState.STARTED,
  },
];

const subSequencesWarning = [
  {
    ...defaultSubsequenceDelivery,
    result: ResultTypes.WARNING,
    state: SequenceState.FINISHED,
  },
  {
    ...defaultSubsequenceDelivery,
    result: ResultTypes.PASSED,
    state: SequenceState.STARTED,
  },
];

const stageDeploymentEmpty = {
  name: 'dev',
  state: SequenceState.STARTED,
  lastTimeUpdated: '2021-10-13T10:49:30.005Z',
  openRemediations: [],
  subSequences: [],
  hasEvaluation: false,
};

export { updatedDeploymentMock as UpdatedDeploymentMock };
export { expectedDeploymentMock as ExpectedDeploymentMock };
export { stageDeploymentDeliveryFinishedPass as StageDeploymentDeliveryFinishedPassMock };
export { stageDeploymentRollBackFinishedPass as StageDeploymentRollbackFinishedPassMock };
export { serviceRemediationInformationDevWithRemediation as ServiceRemediationInformationDevWithRemediationMock };
export { serviceRemediationInformationProductionWithRemediation as ServiceRemediationInformationProductionWithRemediationMock };
export { mergedSubsequencesDeliveryRollback as MergedSubSequencesDeliveryRollbackMock };
export { subSequencesFailedAndPassed as SubSequencesFailedAndPassedMock };
export { subSequencesPassedLoading as SubSequencesPassedLoadingMock };
export { subSequencesPassed as SubSequencesPassedMock };
export { subSequencesWarningFailed as SubSequencesWarningFailedMock };
export { subSequencesWarning as SubSequencesWarningMock };
export { stageDeploymentEmpty as StageDeploymentEmptyMock };
